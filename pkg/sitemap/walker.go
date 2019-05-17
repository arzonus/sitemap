package sitemap

import (
	"context"
	"github.com/arzonus/sitemap/pkg/sitemap/node"
	"github.com/arzonus/sitemap/pkg/sitemap/worker"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type Walker struct {
	option
}

type option struct {
	workerCount int
	maxDepth    int
	outputFile  string
	timeout     time.Duration
	ctx         context.Context
	client      *http.Client
}

type Option func(option *option)

func OutputFileOption(path string) Option {
	return func(option *option) {
		option.outputFile = path
	}
}

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
			outputFile:  "sitemap.out",
			timeout:     15 * time.Second,
			workerCount: runtime.NumCPU(),
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
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
		nodes   = make(chan *node.Node, w.workerCount)
		results = make(chan *worker.CrawlerResult, w.workerCount*5)
		done    = make(chan struct{})
	)

	for i := 0; i < w.workerCount; i++ {
		go crawler.Work(nodes, results)
		go parser.Work(results, nodes)
	}

	n := node.NewNode(ctx, u, done)
	nodes <- n
	<-done

	log.Print("finished")

	if n.Error() != nil {
		return n, err
	}
	return n, nil
}
