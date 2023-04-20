package fetcher

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Fetcher interface {
	Fetch(*url.URL) (string, error)
}

type HttpFetcher struct {
}

func (f *HttpFetcher) Fetch(linkURL *url.URL) (string, error) {
	resp, err := http.Get(linkURL.String())
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	return string(bs), nil
}
