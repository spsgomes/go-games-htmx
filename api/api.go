package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	sqlite "go-games-htmx/database"

	"github.com/joho/godotenv"
)

type SearchResults struct {
	Error       string `json:"error"`
	ErrorStatus int    `json:"status_code"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
	Pages       int    `json:"number_of_page_results"`
	Total       int    `json:"number_of_total_results"`
	Results     []struct {
		Id               int    `json:"id,omitempty"`
		Guid             string `json:"guid,omitempty"`
		Name             string `json:"name,omitempty"`
		ShortDescription string `json:"deck,omitempty"`
		Images           struct {
			Small  string `json:"small_url"`
			Medium string `json:"medium_url"`
			Large  string `json:"super_url"`
		} `json:"image,omitempty"`
		ReleaseDate string `json:"original_release_date,omitempty"`
		GameRating  []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			ApiDetailUrl string `json:"api_detail_url"`
		} `json:"original_game_rating,omitempty"`
		Platforms []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			ApiDetailUrl string `json:"api_detail_url"`
		} `json:"platforms,omitempty"`
		DateAdded    string `json:"date_added,omitempty"`
		ResourceType string `json:"resource_type,omitempty"`
		ApiDetailUrl string `json:"api_detail_url"`
	} `json:"results"`
}

func newSearchResults() SearchResults {
	return SearchResults{}
}

type GameResultsError struct {
	Error string `json:"error"`
}

type GameResults struct {
	Error       string `json:"error"`
	ErrorStatus int    `json:"status_code"`
	Results     struct {
		Id               int    `json:"id,omitempty"`
		Guid             string `json:"guid,omitempty"`
		Name             string `json:"name,omitempty"`
		ShortDescription string `json:"deck,omitempty"`
		Description      string `json:"description,omitempty"`
		Images           struct {
			Small  string `json:"small_url"`
			Medium string `json:"medium_url"`
			Large  string `json:"super_url"`
		} `json:"image,omitempty"`
		ReleaseDate string `json:"original_release_date,omitempty"`
		GameRating  []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			ApiDetailUrl string `json:"api_detail_url"`
		} `json:"original_game_rating,omitempty"`
		Platforms []struct {
			Id           int    `json:"id"`
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
			ApiDetailUrl string `json:"api_detail_url"`
		} `json:"platforms,omitempty"`
		Developers []struct {
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"developers,omitempty"`
		Genres []struct {
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"genres,omitempty"`
		Publishers []struct {
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"publishers,omitempty"`
		SimilarGames []struct {
			Id   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"similar_games,omitempty"`
		DateAdded    string `json:"date_added,omitempty"`
		ApiDetailUrl string `json:"api_detail_url"`
	} `json:"results"`
}

func newGameResults() GameResults {
	return GameResults{}
}

const API_URL string = "https://www.giantbomb.com/api"

func CheckApiKey() error {
	_, err := getApiKey()
	if err != nil {
		return err
	}

	return nil
}

func Search(s string, page string) (SearchResults, error) {

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return newSearchResults(), err
	}

	var responseBody []byte

	cached, err := sqlite.GetSearchResults(s, pageInt)

	// Return cached results if they exist
	if err == nil && cached.Id > 0 {
		responseBody = []byte(cached.Results)

		// Fetch from API if no cached results exist
	} else {
		apiKey, err := getApiKey()
		if err != nil {
			return newSearchResults(), err
		}

		req, err := http.NewRequest(http.MethodGet, API_URL+"/search/", nil)
		if err != nil {
			return newSearchResults(), err
		}

		q := req.URL.Query()
		q.Add("api_key", apiKey)
		q.Add("format", "json")
		q.Add("limit", "10")
		q.Add("page", page)
		q.Add("query", s)
		q.Add("field_list", "api_detail_url,date_added,date_updated,deck,guid,id,image,name,original_game_rating,original_release_date,platforms")
		q.Add("resources", "game")
		req.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return newSearchResults(), err
		}

		if response.StatusCode != http.StatusOK {
			return newSearchResults(), fmt.Errorf("status code %v (\"%s\") provided", response.StatusCode, response.Status)
		}

		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			return newSearchResults(), err
		}
	}

	var searchResults SearchResults
	err = json.Unmarshal(responseBody, &searchResults)
	if err != nil {
		return newSearchResults(), err
	}

	// Add to database
	sqlite.AddSearchResults(s, pageInt, string(responseBody))

	return searchResults, nil
}

func Game(guid string) (GameResults, error) {

	var responseBody []byte

	cached, err := sqlite.GetGameResults(guid)

	// Return cached results if they exist
	if err == nil && cached.Id > 0 {
		responseBody = []byte(cached.Results)

		// Fetch from API if no cached results exist
	} else {
		apiKey, err := getApiKey()
		if err != nil {
			return newGameResults(), err
		}

		req, err := http.NewRequest(http.MethodGet, API_URL+"/game/"+guid, nil)
		if err != nil {
			return newGameResults(), err
		}

		q := req.URL.Query()
		q.Add("api_key", apiKey)
		q.Add("format", "json")
		q.Add("limit", "10")
		q.Add("field_list", "api_detail_url,date_added,date_updated,deck,description,guid,id,image,name,original_game_rating,original_release_date,platforms,developers,genres,publishers,similar_games")
		req.URL.RawQuery = q.Encode()

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return newGameResults(), err
		}

		if response.StatusCode != http.StatusOK {
			return newGameResults(), fmt.Errorf("status code %v (\"%s\") provided", response.StatusCode, response.Status)
		}

		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			return newGameResults(), err
		}
	}

	var GameResultsError GameResultsError
	err = json.Unmarshal(responseBody, &GameResultsError)
	if err != nil {
		return newGameResults(), err
	}

	// Error
	if GameResultsError.Error != "OK" {
		return newGameResults(), err
	}

	var GameResults GameResults
	err = json.Unmarshal(responseBody, &GameResults)
	if err != nil {
		return newGameResults(), err
	}

	// Add to database
	sqlite.AddGame(guid, string(responseBody))

	return GameResults, nil
}

func getApiKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", errors.New("error loading .env file")
	}

	apiKey := os.Getenv("GB_API_KEY")

	if apiKey == "" {
		return "", errors.New("error reading the API Key")
	}

	return apiKey, nil
}
