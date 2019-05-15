package worker

import (
	"context"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"golang.org/x/net/html"
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
	n, err := html.Parse(r.r)
	if err != nil {
		r.node.SetError(err)
		return
	}

	var urls []*url.URL
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
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

	r.node.SetURLs(urls)

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
