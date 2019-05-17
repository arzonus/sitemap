package worker

import (
	"context"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/url"
)

type Parser struct {
	ctx context.Context
}

func NewParser(ctx context.Context) *Parser {
	return &Parser{ctx: ctx}
}

func (w Parser) Work(in <-chan *CrawlerResult, out chan<- *node.Node) {
	for {
		select {
		case <-w.ctx.Done():
			log.Print("parser closed")
			return
		case result, ok := <-in:
			if !ok {
				return
			}
			w.work(result, out)
		}
	}
}

func (w Parser) work(r *CrawlerResult, out chan<- *node.Node) {

	urls, err := w.parse(r.r)
	r.node.SetResult(&node.Result{
		URLs:  urls,
		Error: err,
	})

	go func() {
		for i := 0; i < len(r.node.Nodes()); i++ {
			select {
			case <-w.ctx.Done():
				return
			default:
				out <- r.node.Nodes()[i]
			}
		}
	}()
}

// parse provides method for getting urls from html site
func (w Parser) parse(r io.Reader) ([]*url.URL, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var (
		f       func(*html.Node)
		urls    []*url.URL
		baseUrl string
	)

	f = func(n *html.Node) {
		// searching base element
		if n.Type == html.ElementNode && n.Data == "base" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
					if err != nil || !u.IsAbs() {
						return
					}

					baseUrl = attr.Val
				}
			}
		}
		// searching a element
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					var val = attr.Val

					if baseUrl != "" {
						val = baseUrl + val
					}

					u, err := url.Parse(val)
					if err != nil || !u.IsAbs() {
						return
					}
					urls = append(urls, u)
					return
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	return urls, nil
}
