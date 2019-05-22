package worker

import (
	"context"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"golang.org/x/net/html"
	"io"
	"net/url"
)

type Parser struct {
	ctx context.Context
}

func NewParser(ctx context.Context) *Parser {
	return &Parser{ctx: ctx}
}

func (w *Parser) Work(in <-chan *CrawlerResult) <-chan *node.Node {
	var out = make(chan *node.Node)

	go func() {
		defer close(out)
		for {
			select {
			case <-w.ctx.Done():
				return
			case result, ok := <-in:
				if !ok {
					return
				}
				w.work(result, out)
			}
		}
	}()

	return out
}

func (w Parser) work(r *CrawlerResult, out chan<- *node.Node) {
	var p = new(parser)
	urls, err := p.Parse(r.r)
	r.node.SetResult(&node.Result{
		URLs:  urls,
		Error: err,
	})

	for i := 0; i < len(r.node.Nodes()); i++ {
		select {
		case <-w.ctx.Done():
			return
		default:
			out <- r.node.Nodes()[i]
		}
	}
}

type parser struct {
	baseURL string
	urls    []*url.URL
}

func (p *parser) Parse(r io.Reader) ([]*url.URL, error) {
	n, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	p.parse(n)
	return p.urls, nil
}

func (p *parser) parse(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a":
			p.parseATag(n)
		case "base":
			p.parseBaseTag(n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.parse(c)
	}
}

func (p *parser) parseBaseTag(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			u, err := url.Parse(attr.Val)
			if err != nil || !(u.IsAbs() && u.Hostname() != "") {
				return
			}

			p.baseURL = attr.Val
		}
	}
}

func (p *parser) parseATag(n *html.Node) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			u, err := url.Parse(attr.Val)
			if err != nil {
				return
			}

			if u.IsAbs() && u.Hostname() != "" {
				p.urls = append(p.urls, u)
				return
			}

			if p.baseURL == "" {
				return
			}

			// if u is not abs, but baseURL is not empty
			u, err = url.Parse(p.baseURL + attr.Val)
			if err != nil {
				return
			}

			if !(u.IsAbs() && u.Hostname() != "") {

				return
			}

			p.urls = append(p.urls, u)
		}
	}
}
