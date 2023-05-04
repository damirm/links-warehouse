package fetcher

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Fetcher interface {
	Fetch(context.Context, *url.URL) (string, error)
}

type HttpFetcher struct {
}

func (f *HttpFetcher) Fetch(ctx context.Context, linkURL *url.URL) (string, error) {
	// TODO: Handle redirects.
	resp, err := http.Get(linkURL.String())
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	return string(bs), err
}
