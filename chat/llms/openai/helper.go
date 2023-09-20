package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type callback[T any] func(response T, done bool, err error)

var (
	StreamPrefixDATA  = []byte("data: ")
	StreamPrefixERROR = []byte("error: ")
	StreamDataDONE    = []byte("[DONE]")
)

type MultipartFormDataRequestBody interface {
	ToMultipartFormData() (*bytes.Buffer, string, error)
}

func execute[T any](client *OpenAI, req *http.Request, response *T, cb callback[T]) error {
	if client.HTTPClient == nil {
		client.HTTPClient = http.DefaultClient
	}
	httpres, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if httpres.StatusCode >= 400 {
		defer httpres.Body.Close()
		return client.apiError(httpres)
	}
	if cb != nil {
		go listen(httpres, cb)
		return nil
	}
	defer httpres.Body.Close()
	if err := json.NewDecoder(httpres.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response to %T: %v", response, err)
	}
	return nil
}

// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format
func listen[T any](res *http.Response, cb callback[T]) {
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		var entry T
		b := scanner.Bytes()
		switch {
		case len(b) == 0:
			continue
		case bytes.HasPrefix(b, StreamPrefixDATA):
			if bytes.HasSuffix(b, StreamDataDONE) {
				cb(entry, true, nil)
				return
			}
			if err := json.Unmarshal(b[len(StreamPrefixDATA):], &entry); err != nil {
				cb(entry, true, err)
				return
			}
			cb(entry, false, nil)
			// TODO: Any error case?
			// case bytes.HasPrefix(b, StreamPrefixERROR):
			// 	cb(entry, true, fmt.Errorf(string(b)))
			// 	return
			// TODO: Any other case?
			// default:
			// 	cb(entry, true, fmt.Errorf(string(b)))
			// 	return
		}
	}
}

func call[T any](ctx context.Context, client *OpenAI, method string, p string, body interface{}, resp T, cb callback[T]) (T, error) {
	req, err := client.build(ctx, method, p, body)
	if err != nil {
		return resp, err
	}
	err = execute(client, req, &resp, cb)
	return resp, err
}

func (l *OpenAI) endpoint(p string) (string, error) {
	if l.APIHost == "" {
		l.APIHost = DefaultOpenAIAPIURL
	}
	u, err := url.Parse(l.APIHost)
	if err != nil {
		return "", err
	}
	u.Path = strings.Join([]string{
		strings.TrimRight(u.Path, "/"),
		strings.TrimLeft(p, "/"),
	}, "/")
	return u.String(), nil
}

func (l *OpenAI) build(ctx context.Context, method, p string, body interface{}) (req *http.Request, err error) {
	endpoint, err := l.endpoint(p)
	if err != nil {
		return nil, err
	}
	r, contenttype, err := l.bodyToReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to build request buf from given body: %v", err)
	}
	req, err = http.NewRequest(method, endpoint, r)
	if err != nil {
		return nil, fmt.Errorf("failed to init request: %v", err)
	}
	req.Header.Add("Content-Type", contenttype)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", l.APIKey))
	if l.Organization != "" {
		req.Header.Add("OpenAI-Organization", l.Organization)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	return req, nil
}

func (l *OpenAI) bodyToReader(body interface{}) (io.Reader, string, error) {
	var r io.Reader
	switch v := body.(type) {
	// case io.Reader:
	// 	r = v
	case nil:
		r = nil
	case MultipartFormDataRequestBody: // TODO: Refactor
		buf, ct, err := v.ToMultipartFormData()
		if err != nil {
			return nil, "", err
		}
		return buf, ct, nil
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, "", err
		}
		r = bytes.NewBuffer(b)
	}
	return r, "application/json", nil
}

func (l *OpenAI) apiError(res *http.Response) error {

	return fmt.Errorf("status: %s, status_code: %d", res.Status, res.StatusCode)
	// errbody := struct {
	// 	Error APIError `json:"error"`
	// }{APIError{Status: res.Status, StatusCode: res.StatusCode}}

	// if err := json.NewDecoder(res.Body).Decode(&errbody); err != nil {
	// 	return fmt.Errorf("failed to decode error body: %v", err)
	// }

	// return errbody.Error
}
