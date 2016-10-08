package animecmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// animeV1 holds anime data returned by the Hummingbird v1 API.
type animeV1 struct {
	ID              int       `json:"id"`
	MalID           int       `json:"mal_id"`
	Slug            string    `json:"slug"`
	Status          string    `json:"status"`
	URL             string    `json:"url"`
	Title           string    `json:"title"`
	AlternateTitle  string    `json:"alternate_title"`
	EpisodeCount    int       `json:"episode_count"`
	EpisodeLength   int       `json:"episode_length"`
	CoverImage      string    `json:"cover_image"`
	Synopsis        string    `json:"synopsis"`
	ShowType        string    `json:"show_type"`
	StartedAiring   string    `json:"started_airing"`
	FinishedAiring  string    `json:"finished_airing"`
	CommunityRating float32   `json:"community_rating"`
	AgeRating       string    `json:"age_rating"`
	Genres          []genreV1 `json:"genres"`
}

// genreV1 holds anime genre data returned by the Hummingbird v1 API.
type genreV1 struct {
	Name string `json:"name"`
}

// anime holds a show returned by the Hummingbird v1 API.
type anime struct {
	Anime  animeData  `json:"anime"`
	Linked linkedData `json:"linked"`
}

// animeData holds anime data returned by the Hummingbird API v2.
type animeData struct {
	ID                 int      `json:"id"`
	Titles             title    `json:"titles"`
	Slug               string   `json:"slug"`
	Synopsis           string   `json:"synopsis"`
	StartedAiringDate  string   `json:"started_airing_date"`
	FinishedAiringDate string   `json:"finished_airing_date"`
	YoutubeVideoID     string   `json:"youtube_video_id"`
	AgeRating          string   `json:"age_rating"`
	EpisodeCount       int      `json:"episode_count"`
	EpisodeLength      int      `json:"episode_length"`
	ShowType           string   `json:"show_type"`
	PosterImage        string   `json:"poster_image"`
	CoverImage         string   `json:"cover_image"`
	CommunityRating    float32  `json:"community_rating"`
	Genres             []string `json:"genres"`
	BayesianRating     float32  `json:"bayesian_rating"`
	Links              []link   `json:"links"`
}

// title holds anime title data returned by the Hummingbird API.
type title struct {
	Canonical string `json:"canonical"`
	English   string `json:"english"`
	Romanji   string `json:"romanji"`
	Japanese  string `json:"japanese"`
}

// link holds anime link data returned by the Hummingbird API.
type link struct {
	GalleryImages int `json:"gallery_images"`
	Episodes      int `json:"episodes"`
}

// linkedData holds linked data returned by the Hummingbird API.
type linkedData struct {
	GalleryImages []galleryImage `json:"gallery_images"`
	Episodes      []episode      `json:"episodes"`
}

// galleryImage holds anime gallery image data returned by the Hummingbird API.
type galleryImage struct {
	ID       int    `json:"id"`
	Thumb    string `json:"thumb"`
	Original string `json:"original"`
}

// episode holds anime episode data returned by the Hummingbird API.
type episode struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Synopsis     string `json:"synopsis"`
	AirDate      string `json:"airdate"`
	Number       int    `json:"number"`
	SeasonNumber int    `json:"season_number"`
}

// hummingbirdSearch queries the Hummingbird v1 API and returns search results
// as a slice of animeV1 structs.
func hummingbirdSearch(query string) ([]animeV1, error) {
	// Escape query string and send it to Hummingbird API v1.
	resp, err := http.Get("http://hummingbird.me/api/v1/search/anime?query=" +
		url.QueryEscape(query))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Read data body into json string.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse json string into slice of animeV1 structs.
	shows := []animeV1{}
	json.Unmarshal(body, &shows)

	// Return slice with a maximum of five elements.
	if len(shows) < 5 {
		return shows, nil
	}
	return shows[:5], nil
}

// hummingbirdShow queries the Hummingbird API and returns show data.
func hummingbirdShow(slug string) (*anime, error) {
	// Create new http request client, escape slug, set client ID header, and
	// send it to Hummingbird API.
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://hummingbird.me/api/v2/anime/"+
		url.QueryEscape(slug), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Client-Id", "4a4d0d2045810f2a975f")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Read data body into json string.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse json string into anime struct.
	show := anime{}
	json.Unmarshal(body, &show)

	// Return show.
	return &show, nil
}
