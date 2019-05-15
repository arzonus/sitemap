package sitemap

import (
	"context"
	"fmt"
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
	url string
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

func NewWalker(url string, options ...Option) *Walker {
	w := &Walker{
		url: url,
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

func (w *Walker) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()
	return w.run(ctx)
}

func (w *Walker) RunContext(ctx context.Context) error {
	return w.run(ctx)
}

func (w *Walker) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	url, err := url.Parse(w.url)
	if err != nil {
		return err
	}
	log.Print("try to fetch from ", url)

	crawler := worker.NewCrawler(w.client, ctx, w.maxDepth)
	parser := worker.NewParser(ctx)

	nodes := make(chan *node.Node, w.workerCount)
	results := make(chan *worker.CrawlerResult, w.workerCount*5)
	done := make(chan struct{})

	log.Print("create workers ", w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		go crawler.Work(nodes, results)
		go parser.Work(results, nodes)
	}
	log.Print("create new node")
	log.Println(len(nodes), cap(nodes))
	node := node.NewNode(ctx, url, done)
	log.Print("send first node")
	nodes <- node
	log.Println(len(nodes), cap(nodes))
	log.Print("send to channel")
	<-done
	log.Print("finished")

	fmt.Println(node.String())
	if node.Error() != nil {
		return nil
	}
	return nil
}
