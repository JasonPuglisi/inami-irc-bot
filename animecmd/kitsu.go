package animecmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type container struct {
	Results []result `json:"data"`
}

type result struct {
	ID         string     `json:"id"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Slug   string `json:"slug"`
	Title  string `json:"canonicalTitle"`
	Number int    `json:"number"`
}

// search queries the Kitsu API and returns search results as a slice.
func search(query string) ([]result, error) {
	// Escape query string and send it to Kitsu API.
	resp, err := http.Get(
		fmt.Sprintf("https://kitsu.io/api/edge/anime?page[limit]=5&filter[text]=%s",
			url.QueryEscape(query)))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

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
	defer resp.Body.Close()
	if err != nil {
		return result{}, err
	}

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
	// Escape query string and send it to Kitsu API.
	resp, err := http.Get(
		fmt.Sprintf("https://kitsu.io/api/edge/anime/%s/episodes",
			url.QueryEscape(id)))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Read data body into json string.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse json string into container struct.
	data := container{}
	json.Unmarshal(body, &data)

	// Return slice of episodes.
	return data.Results, nil
}
