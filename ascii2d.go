package ascii2dgo

import (
	"net/url"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

var (
	baseURL = "https://ascii2d.net/search/url/"
)

type Result struct {
	Sources []*Source
}

type Source struct {
	Author    *Author
	Title     string
	Thumbnail string
	URL       string
}

type Author struct {
	Name string
	URL  string
}

func Search(file string) (*Result, error) {
	col := colly.NewCollector(colly.Async(true))

	var (
		result = &Result{make([]*Source, 0)}
	)

	col.OnRequest(func(r *colly.Request) {
		logrus.Infof("Making a request. URL: %s", r.URL)
	})

	col.OnError(func(r *colly.Response, err error) {
		logrus.Errorf("Request errored: %v", err)
	})

	col.OnHTML(".item-box", func(h *colly.HTMLElement) {
		thumbnail := h.ChildAttr(".image-box > img", "src")

		links := h.ChildAttrs(".info-box > .detail-box > h6 > a", "href")
		texts := h.ChildTexts(".info-box > .detail-box > h6 > a")

		if links == nil || texts == nil {
			return
		}

		source := &Source{
			Title:     texts[0],
			Thumbnail: "https://ascii2d.net" + thumbnail,
			URL:       links[0],
			Author:    &Author{Name: texts[1], URL: links[1]},
		}
		result.Sources = append(result.Sources, source)
	})

	col.Visit(baseURL + url.QueryEscape(file))
	col.Wait()

	return result, nil
}
