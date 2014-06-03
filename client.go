package chroma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var defaultHttp *http.Client = &http.Client{}

type Client struct {
	// The IP address or hostname to the Hue Bridge.
	BridgeHost string

	// A username is required to interact with the Hue API. Make it unique for
	// your application or use case.
	Username string

	http *http.Client
}

// Authorizes the username to access the Hue Bridge. The link button must be
// pressed.
func (c *Client) Authorize(deviceType string) error {
	var resp interface{}
	_, err := c.post(c.buildApiUrlWithoutUsername(""),
		map[string]interface{}{"devicetype": deviceType, "username": c.Username},
		&resp)

	return err
}

func (c *Client) get(url *url.URL, resp interface{}) (httpResp *http.Response, err error) {
	return c.do(&http.Request{
		Method: "GET",
		URL:    url,
	}, nil, resp)
}

func (c *Client) post(url *url.URL, body interface{}, resp interface{}) (httpResp *http.Response, err error) {
	return c.do(&http.Request{
		Method: "POST",
		URL:    url,
	}, body, resp)
}

func (c *Client) put(url *url.URL, body interface{}, resp interface{}) (httpResp *http.Response, err error) {
	return c.do(&http.Request{
		Method: "PUT",
		URL:    url,
	}, body, resp)
}

func (c *Client) do(req *http.Request, body interface{}, resp interface{}) (httpResp *http.Response, err error) {
	httpClient := c.http
	if httpClient == nil {
		httpClient = defaultHttp
	}

	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		} else {
			if req.Header == nil {
				req.Header = make(http.Header, 0)
			}

			req.Header["Content-Type"] = []string{"application/json"}
			req.Body = nopCloser{bytes.NewBuffer(encoded)}
			req.ContentLength = int64(len(encoded))
		}
	}

	httpResp, err = httpClient.Do(req)
	if err != nil {
		return httpResp, err
	}
	defer httpResp.Body.Close()

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return httpResp, err
	}

	return
}

func (c *Client) buildApiUrl(path string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   c.BridgeHost,
		Path:   fmt.Sprintf("/api/%s/%s", c.Username, path),
	}
}

func (c *Client) buildApiUrlWithoutUsername(path string) *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   c.BridgeHost,
		Path:   fmt.Sprintf("/api/%s", path),
	}
}
