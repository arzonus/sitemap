package sitemap

import (
	"net/url"
)

//
//import (
//	"context"
//	"golang.org/x/net/html"
//	"io"
//	"net/http"
//	"net/url"
//)
//
//type CrawlerPool struct {
//	workerCount int
//}
//
//func NewCrawlerPool(workerCount int) *CrawlerPool {
//	return &CrawlerPool{workerCount: workerCount}
//}
//
//func (p CrawlerPool) Run(ctx context.Context) error {
//
//}
//
//type Crawler struct {
//	client *http.Client
//}
//
//func (c Crawler) Do(ctx context.Context, url *url.URL) ([]url.URL, error) {
//
//	if Node.depth > c.maxDepth {
//		Node.Response.Error = ErrDepthExceeded
//		return
//	}
//
//	req := http.Request{
//		Method: "GET",
//		URL:    Node.URL,
//	}.WithContext(ctx)
//
//	resp, err := c.client.Do(req)
//	if err != nil {
//		Node.Response.Error = err
//		return
//	}
//
//	urls, err := c.Parse(resp.Body)
//	if err != nil {
//		Node.Response.Error = err
//		return
//	}
//
//	Node.Nodes = make([]Site, len(urls))
//	for i, url := range urls {
//		Node.Nodes[i].URL = url
//		Node.Nodes[i].depth = Node.depth + 1
//	}
//}
//
//func (c Crawler) Parse(r io.Reader) ([]*url.URL, error) {
//	_, err := html.Parse(r)
//	if err != nil {
//		return nil, err
//	}
//
//	return nil, nil
//}

type Site struct {
	depth int
	URL   *url.URL

	Response HTTPResponse
	Sites    []Site
}

type HTTPResponse struct {
	Error error
}



//
//func (s Site) Done() <-chan struct{} {
//
//	var done = make(chan struct{})
//
//	go func(){
//		defer close(done)
//
//
//
//	}()
//
//	return done
//}
//
//func (s *Site) NewSite(url *url.URL) *Site {
//	site := &Site{
//		depth: s.depth + 1,
//	}
//}
//
//
//type HTTPClient struct {
//	client *http.Client
//}
//
//func (c HTTPClient) Do(ctx context.Context, in <-chan *Site, out chan<- *Site) {
//
//}
