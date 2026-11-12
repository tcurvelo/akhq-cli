// Package akhq is a client for the AKHQ HTTP API.
package akhq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	BaseURL string
	Cluster string
	Token   string // JWT sent as a cookie; set by the caller or from Login
	HTTP    *http.Client
}

func New(baseURL, cluster string, insecure bool) *Client {
	return &Client{
		BaseURL: baseURL,
		Cluster: cluster,
		HTTP: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
			// AKHQ answers a successful login with a redirect that carries the JWT cookie.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login returns the JWT from the login response, or "" if there is none.
func (c *Client) Login(user, pass string) (string, error) {
	body, err := json.Marshal(loginRequest{Username: user, Password: pass})
	if err != nil {
		return "", err
	}
	resp, err := c.HTTP.Post(c.BaseURL+"/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JWT" {
			return cookie.Value, nil
		}
	}
	return "", nil
}

type ConsumerGroup struct {
	Offsets []PartitionOffset `json:"offsets"`
}

type PartitionOffset struct {
	OffsetLag *int64 `json:"offsetLag"` // nil when the group never committed on this partition
}

// TopicGroups returns the consumer groups of a topic with per-partition offset lag.
func (c *Client) TopicGroups(topic string) ([]ConsumerGroup, error) {
	var groups []ConsumerGroup
	err := c.get(c.topicPath(topic, "groups"), &groups)
	return groups, err
}

type SearchParams struct {
	After               string
	Partition           int32
	Sort                string // OLDEST or NEWEST
	Timestamp           string
	SearchByKey         string
	SearchByValue       string
	SearchByHeaderKey   string
	SearchByHeaderValue string
}

// Search searches a topic's data.
// shortcut: result is untyped because the response shape is unconfirmed (maybe SSE); type it in #8.
func (c *Client) Search(topic string, p SearchParams) (any, error) {
	q := url.Values{}
	q.Set("partition", strconv.Itoa(int(p.Partition)))
	for key, val := range map[string]string{
		"after":               p.After,
		"sort":                p.Sort,
		"timestamp":           p.Timestamp,
		"searchByKey":         p.SearchByKey,
		"searchByValue":       p.SearchByValue,
		"searchByHeaderKey":   p.SearchByHeaderKey,
		"searchByHeaderValue": p.SearchByHeaderValue,
	} {
		if val != "" {
			q.Set(key, val)
		}
	}
	var data any
	err := c.get(c.topicPath(topic, "data/search")+"?"+q.Encode(), &data)
	return data, err
}

func (c *Client) topicPath(topic, rest string) string {
	return fmt.Sprintf("%s/api/%s/topic/%s/%s", c.BaseURL, url.PathEscape(c.Cluster), url.PathEscape(topic), rest)
}

func (c *Client) get(rawURL string, out any) error {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{Name: "JWT", Value: c.Token})
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return json.NewDecoder(resp.Body).Decode(out)
}
