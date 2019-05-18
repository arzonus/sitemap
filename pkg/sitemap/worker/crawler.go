package worker

import (
	"context"
	"fmt"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"io"
	"net/http"
)

type Crawler struct {
	client   *http.Client
	ctx      context.Context
	maxDepth int
}

func NewCrawler(client *http.Client, ctx context.Context, maxDepth int) *Crawler {
	return &Crawler{client: client, ctx: ctx, maxDepth: maxDepth}
}

type CrawlerResult struct {
	node *node.Node
	r    io.Reader
}

func (w Crawler) Work(in <-chan *node.Node, out chan<- *CrawlerResult) {
	for {
		select {
		case <-w.ctx.Done():
			return
		case n, ok := <-in:
			if !ok {
				return
			}
			w.work(n, out)
		}
	}
}

var ErrDepthExceeded = fmt.Errorf("depth exceeded")

func (w Crawler) work(node *node.Node, out chan<- *CrawlerResult) {
	if node.Depth() > w.maxDepth {
		node.SetError(ErrDepthExceeded)
		return
	}

	req := &http.Request{
		Method: "GET",
		URL:    node.URL(),
	}

	resp, err := w.client.Do(req)
	if err != nil {
		node.SetError(err)
		return
	}

	select {
	case <-w.ctx.Done():
		node.SetError(fmt.Errorf("context exceeded"))
		return
	default:
		out <- &CrawlerResult{
			node: node,
			r:    resp.Body,
		}
	}
}
