package worker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/url"
	"strings"
	"testing"
)

func NewHTML(baseURL string, urls ...string) io.Reader {
	var str = `
<html>`

	str += newHead(baseURL)
	str += newBody(urls)
	str += `
</html>`

	return strings.NewReader(str)
}

func newHead(baseURL string) string {
	var str = `
		<head>`
	if baseURL != "" {
		str += fmt.Sprintf(`
			<base href="%s">`, baseURL)
	}
	str += `
		</head>`
	return str
}

func newBody(urls []string) string {
	var str = `
		<body>`
	for _, url := range urls {
		str += fmt.Sprintf(`
			<p><a href="%s">%s</a></p>`, url, url)
	}
	str += `
		</body>`
	return str
}

func newURL(t *testing.T, urlRaw string) *url.URL {
	u, err := url.Parse(urlRaw)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestParse(t *testing.T) {
	var cases = []struct {
		Name   string
		Reader io.Reader

		ExpError error
		ExpURLs  []*url.URL
	}{
		{
			Name:     "empty",
			Reader:   NewHTML(""),
			ExpError: nil,
			ExpURLs:  nil,
		},
		{
			Name:     "full url",
			Reader:   NewHTML("", "https://vk.com"),
			ExpError: nil,
			ExpURLs: []*url.URL{
				newURL(t, "https://vk.com"),
			},
		},
		{
			Name:     "base url",
			Reader:   NewHTML("https://vk.com/", "img.img"),
			ExpError: nil,
			ExpURLs: []*url.URL{
				newURL(t, "https://vk.com/img.img"),
			},
		},
		{
			Name:     "full url, base url",
			Reader:   NewHTML("https://vk.com/", "img.img", "https://facebook.com"),
			ExpError: nil,
			ExpURLs: []*url.URL{
				newURL(t, "https://vk.com/img.img"),
				newURL(t, "https://facebook.com"),
			},
		},
		{
			Name:     "empty base url, uri",
			Reader:   NewHTML("", "img.img"),
			ExpError: nil,
			ExpURLs:  nil,
		},
		//{
		//	Name:     "javascript void",
		//	Reader:   NewHTML("http://vk.com/", "javascript:void(8)"),
		//	ExpError: nil,
		//	ExpURLs:  nil,
		//},
		//{
		//	Name:     "err url",
		//	Reader:   NewHTML("http://vk.com/", "/.>"),
		//	ExpError: nil,
		//	ExpURLs:  nil,
		//},
	}

	for _, c := range cases {
		urls, err := Parse(c.Reader)
		assert.Equal(t, c.ExpURLs, urls)
		assert.Equal(t, c.ExpError, err)
	}
}
