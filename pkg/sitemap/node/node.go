package node

import (
	"context"
	"net/url"
)

type Node struct {
	doneChan chan<- struct{}
	nodeChan chan struct{}
	errChan  chan struct{}

	ctx   context.Context
	depth int
	url   *url.URL

	nodes  []*Node
	result *Result
}

func NewNode(
	ctx context.Context,
	url *url.URL,
	doneChan chan<- struct{},
) *Node {
	return newNode(ctx, url, doneChan, 0)
}

func newNode(
	ctx context.Context,
	url *url.URL,
	doneChan chan<- struct{},
	depth int,
) *Node {
	node := &Node{
		ctx:      ctx,
		url:      url,
		doneChan: doneChan,
		depth:    depth,
	}
	go node.wait()
	return node
}

func (n Node) Error() error {
	return n.result.Error
}

func (n Node) URL() *url.URL {
	return n.url
}

func (n Node) Depth() int {
	return n.depth
}

func (n Node) Nodes() []*Node {
	return n.nodes
}

func (n Node) Prefix() string {
	var str string
	for i := 0; i < n.depth; i++ {
		str += "*"
	}
	return str
}

type Result struct {
	Error error
	URLs  []*url.URL
}

func (n *Node) SetError(err error) {
	n.SetResult(&Result{
		Error: err,
	})
}

func (n *Node) SetURLs(urls []*url.URL) {
	n.SetResult(&Result{
		URLs: urls,
	})
}

func (n *Node) SetResult(result *Result) {
	n.result = result
	if result.Error != nil {
		n.errChan = make(chan struct{})
		n.errChan <- struct{}{}
		return
	}
	n.setNodes(result.URLs)
}

func (n *Node) setNodes(urls []*url.URL) {
	n.nodeChan = make(chan struct{}, len(urls))
	n.nodes = make([]*Node, len(urls))

	for i, url := range urls {
		n.nodes[i] = newNode(n.ctx, url, n.nodeChan, n.depth+1)
	}
}

func (n *Node) close() {
	if n.nodeChan != nil {
		close(n.nodeChan)
	}
	if n.errChan != nil {
		close(n.errChan)
	}
}

func (n *Node) wait() {
	defer n.close()

	var i int
	for {
		select {
		case <-n.ctx.Done():
			return

		case <-n.nodeChan:
			i++
			if i == len(n.nodes) {
				n.doneChan <- struct{}{}
				return
			}
		case <-n.errChan:
			return
		}
	}
}
