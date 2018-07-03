package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func resolveRelativeURL(u *url.URL, path string, values string) *url.URL {

	ru := &url.URL{Path: path, RawQuery: values}

	return u.ResolveReference(ru)
}

// Client for football api
type Client struct {
	APIKey  string
	BaseURL *url.URL
	Client  *http.Client
}

// Get makes a get request
func (c Client) Get(resources string, filters string) (j *json.Decoder, err error) {
	url := resolveRelativeURL(c.BaseURL, resources, filters)
	// fmt.Println("url:", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("X-Auth-Token", c.APIKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Println("Err in reading json to buf", err)
		return
	}
	j = json.NewDecoder(buf)
	return
}

// Filters represents filters
type Filters struct {
	ID       int      `json:"id,omitempty"`
	MatchDay int      `json:"matchday,omitempty"`
	Status   []string `json:"status,omitempty"`
	Venue    string   `json:"venue,omitempty"`
	DateFrom string   `json:"dateFrom,omitempty"`
	DateTo   string   `json:"dateTo,omitempty"`
	Stage    string   `json:"stage,omitempty"`
}

type ResponseMatches struct {
	Count   int     `json:"count,omitempty"`
	Filters Filters `json:"filters,omitempty"`
	Matches []Match `json:"matches,omitempty"`
}

type Match struct {
	ID          int         `json:"id,omitempty"`
	Status      string      `json:"status,omitempty"`
	Competition Competition `json:"competition,omitempty"`
	Season      Season      `json:"season,omitempty"`
	UTCDate     string      `json:"utcDate,omitempty"`
	MatchDay    int         `json:"matchday,omitempty"`
	Stage       string      `json:"stage,omitempty"`
	Group       string      `json:"group,omitempty"`
	LastUpdated string      `json:"lastUpdated,omitempty"`
	HomeTeam    Team        `json:"homeTeam,omitempty"`
	AwayTeam    Team        `json:"awayTeam,omitempty"`
	Score       MatchScore  `json:"score,omitempty"`
}

type Competition struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Season struct {
	ID              int    `json:"id,omitempty"`
	StartDate       string `json:"startDate,omitempty"`
	EndDate         string `json:"endDate,omitempty"`
	CurrentMatchDay int    `json:"currentMatchday,omitempty"`
}

type Team struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type MatchScore struct {
	Winner    string `json:"winner,omitempty"`
	Duration  string `json:"duration,omitempty"`
	FullTime  Score  `json:"fullTime,omitempty"`
	HalfTime  Score  `json:"halfTime,omitempty"`
	ExtraTime Score  `json:"extraTime,omitempty"`
	Penalties Score  `json:"penalties,omitempty"`
}

type Score struct {
	HomeTeam int `json:"homeTeam,omitempty"`
	AwayTeam int `json:"awayTeam,omitempty"`
}
