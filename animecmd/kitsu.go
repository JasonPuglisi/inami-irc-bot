package animecmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bradfitz/gomemcache/memcache"
)

// container stores data returned by the Kitsu API.
type container struct {
	Results []result `json:"data"`
	Links   links    `json:"links"`
}

// result stores an entry returned by the Kitsu API.
type result struct {
	ID         string     `json:"id"`
	Attributes attributes `json:"attributes"`
}

// attributes stores data in an entry returned by the Kitsu API.
type attributes struct {
	Slug   string `json:"slug"`
	Title  string `json:"canonicalTitle"`
	Number int    `json:"number"`
}

// links stores references to additional data returned by query.
type links struct {
	Next string `json:"next"`
}

// search queries the Kitsu API and returns search results as a slice.
func search(query string) ([]result, error) {
	// Escape query string and send it to Kitsu API.
	resp, err := http.Get(
		fmt.Sprintf("https://kitsu.io/api/edge/anime?page[limit]=5&filter[text]=%s",
			url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read data body into json string.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse json string into container struct.
	data := container{}
	json.Unmarshal(body, &data)

	// Return slice of results with a maximum of five elements.
	if len(data.Results) < 5 {
		return data.Results, nil
	}
	return data.Results[:5], nil
}

// show queries the Kitsu API and returns show data.
func show(id string) (result, error) {
	// Escape query string and send it to Kitsu API for show data.
	resp, err := http.Get(
		fmt.Sprintf("https://kitsu.io/api/edge/anime?page[limit]=1&filter[id]=%s",
			url.QueryEscape(id)))
	if err != nil {
		return result{}, err
	}
	defer resp.Body.Close()

	// Read data body into json string.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result{}, err
	}

	// Parse json string into container struct.
	data := container{}
	json.Unmarshal(body, &data)

	// Return struct of show.
	return data.Results[0], nil
}

// episodes queries the Kitsu API and returns episode data.
func episodes(id string) ([]result, error) {
	// Create slice to store results.
	var results []result

	// Connect to memcached and attempt to return episodes from cache.
	mc := memcache.New("127.0.0.1:11211")
	it, err := mc.Get(fmt.Sprintf("episodes_%s", id))
	if err == nil {
		err = json.Unmarshal(it.Value, &results)
		if err == nil {
			return results, nil
		}
	}

	// Initiate loop to hit all pages.
	more, link := true, fmt.Sprintf(
		"https://kitsu.io/api/edge/anime/%s/episodes?page[limit]=20",
		url.QueryEscape(id))
	for more {
		// Escape query string and send it to Kitsu API.
		resp, err2 := http.Get(link)
		if err2 != nil {
			return nil, err2
		}
		defer resp.Body.Close()

		// Read data body into json string.
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return nil, err2
		}

		// Parse json string into container struct.
		data := container{}
		json.Unmarshal(body, &data)

		// Update results slice.
		results = append(results, data.Results...)

		if data.Links.Next != "" {
			link = data.Links.Next
		} else {
			more = false
		}
	}

	// Send episodes to cache with expiry of 12 hours.
	episodes, err := json.Marshal(results)
	if err == nil {
		mc.Set(&memcache.Item{Key: fmt.Sprintf("episodes_%s", id), Value: episodes})
		mc.Touch(fmt.Sprintf("episodes_%s", id), 43200)
	}

	// Return slice of episodes.
	return results, nil
}
