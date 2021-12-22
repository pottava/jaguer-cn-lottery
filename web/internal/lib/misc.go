package lib

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPGet send a HTTP request
func HTTPGet(ctx context.Context, client *http.Client, endpoint string, headers *http.Header) (int, []byte, error) {
	return httpOnce(ctx, client, http.MethodGet, endpoint, nil, headers)
}

// HTTPPost send a HTTP request
func HTTPPost(ctx context.Context, client *http.Client, endpoint string, body io.Reader, headers *http.Header) (int, []byte, error) {
	return httpOnce(ctx, client, http.MethodPost, endpoint, body, headers)
}

func httpOnce(ctx context.Context, client *http.Client, method, endpoint string, body io.Reader, headers *http.Header) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return 0, nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}
