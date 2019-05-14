package worker

import (
	"context"
	"fmt"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"io"
	"log"
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

var ErrDepthExceeded = fmt.Errorf("err depth exceeded")

func (w Crawler) work(node *node.Node, out chan<- *CrawlerResult) {
	if node.Depth() > w.maxDepth {
		log.Println(node.Prefix(), "got depth exceeded: ", node.URL())
		node.SetError(ErrDepthExceeded)
		return
	}

	req := &http.Request{
		Method: "GET",
		URL:    node.URL(),
	}

	log.Println(node.Prefix(), "try to request url: ", node.URL(), " ", node.Depth())

	resp, err := w.client.Do(req)
	if err != nil {
		log.Println(node.Prefix(), "got err: ", err)
		node.SetError(err)
		return
	}
	log.Println(node.Prefix(), "got body: ", node.URL(), " ", node.Depth())

	out <- &CrawlerResult{
		node: node,
		r:    resp.Body,
	}
}
