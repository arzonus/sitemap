package sitemap

import (
	"context"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"github.com/arzonus/sitemap/pkg/sitemap/worker"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

type Walker struct {
	option
}

type option struct {
	workerCount int
	maxDepth    int
	timeout     time.Duration
	ctx         context.Context
	client      *http.Client
	bufCount    int
}

type Option func(option *option)

func WorkerCountOption(count int) Option {
	return func(option *option) {
		option.workerCount = count
	}
}

func MaxDepthOption(depth int) Option {
	return func(option *option) {
		option.maxDepth = depth
	}
}

func TimeoutOption(dur time.Duration) Option {
	return func(option *option) {
		option.timeout = dur
	}
}

func HTTPClientOption(client *http.Client) Option {
	return func(option *option) {
		option.client = client
	}
}

func NewWalker(options ...Option) *Walker {
	w := &Walker{
		option: option{
			maxDepth:    5,
			timeout:     15 * time.Second,
			workerCount: runtime.NumCPU(),
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			bufCount: 1000,
		},
	}

	for _, apply := range options {
		apply(&w.option)
	}

	return w
}

func (w *Walker) Walk(urlRaw string) (*node.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()
	return w.walk(ctx, urlRaw)
}

func (w *Walker) WalkContext(ctx context.Context, urlRaw string) (*node.Node, error) {
	return w.walk(ctx, urlRaw)
}

func (w *Walker) walk(ctx context.Context, urlRaw string) (*node.Node, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u, err := url.Parse(urlRaw)
	if err != nil {
		return nil, err
	}

	var (
		crawler = worker.NewCrawler(w.client, ctx, w.maxDepth)
		parser  = worker.NewParser(ctx)
		nodes   = make(chan *node.Node, w.workerCount*w.bufCount)
		results = make(chan *worker.CrawlerResult, w.workerCount*w.bufCount)
		done    = make(chan struct{})
	)
	defer close(done)

	var (
		nodeChan   = make([]<-chan *node.Node, w.workerCount)
		resultChan = make([]<-chan *worker.CrawlerResult, w.workerCount)
	)

	for i := 0; i < w.workerCount; i++ {
		resultChan[i] = crawler.Work(nodes)
		nodeChan[i] = parser.Work(results)
	}

	mergeNodes(nodeChan, nodes)
	mergeResults(resultChan, results)

	n := node.NewNode(ctx, u, done)
	nodes <- n
	<-done

	if n.Error() != nil {
		return n, err
	}
	return n, nil
}

func mergeNodes(cs []<-chan *node.Node, out chan<- *node.Node) {
	var wg sync.WaitGroup
	output := func(c <-chan *node.Node) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

func mergeResults(cs []<-chan *worker.CrawlerResult, out chan<- *worker.CrawlerResult) {
	var wg sync.WaitGroup
	output := func(c <-chan *worker.CrawlerResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}
