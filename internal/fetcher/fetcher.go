package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/damirm/links-warehouse/internal/model"
)

type Fetcher interface {
	Fetch(url.URL) (model.Link, error)
}

type HttpFetcher struct {
}

func (f *HttpFetcher) Fetch(linkURL url.URL) (*model.Link, error) {
	resp, err := http.Get(linkURL.String())
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	body := string(bs)
	fmt.Println(body)
	return &model.Link{}, nil
}
