package cocktail

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the Cocktail API client structure
type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

// C is the exported default client
var C = Client{
	BaseURL: &url.URL{
		Host:   "www.thecocktaildb.com",
		Path:   "/api/json/v1/1/",
		Scheme: "https",
	},
	HTTPClient: &http.Client{
		Timeout: time.Second * 10,
	},
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	var err error
	var resp *http.Response

	if resp, err = c.HTTPClient.Do(req); err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// GetRandomDrink returns a single random FullDrink object
func (c *Client) GetRandomDrink() (*FullDrink, error) {
	var err error
	var req *http.Request
	var d *FullDrink
	var ds *FullDrinkList

	if req, err = c.newRequest("GET", "random.php", nil); err != nil {
		return d, err
	}

	_, err = c.do(req, &ds)
	if len(ds.Drinks) > 0 {
		d = ds.Drinks[0]
	}
	return d, err
}
