package node

import (
	"context"
	"net/url"
	"sync"
)

type Node struct {
	parent *Node

	i        int
	mx       sync.Mutex
	url      *url.URL
	ctx      context.Context
	size     int
	depth    int
	nodes    []*Node
	result   *Result
	isClosed bool
	doneChan chan<- struct{}
	nodeChan chan struct{}
}

func NewNode(
	ctx context.Context,
	url *url.URL,
	doneChan chan<- struct{},
) *Node {
	return newNode(ctx, url, doneChan, 0, nil)
}

func newNode(
	ctx context.Context,
	url *url.URL,
	doneChan chan<- struct{},
	depth int,
	parent *Node,
) *Node {
	node := &Node{
		ctx:      ctx,
		url:      url,
		doneChan: doneChan,
		depth:    depth,
		nodeChan: make(chan struct{}),
		parent:   parent,
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

func (n Node) Parent() *Node {
	return n.parent
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
		n.close()
		return
	}
	n.setNodes(result.URLs)
}

func (n *Node) setNodes(urls []*url.URL) {
	n.mx.Lock()
	defer n.mx.Unlock()
	if n.isClosed {
		return
	}

	n.nodes = make([]*Node, len(urls))
	n.size = len(urls)

	for i, u := range urls {
		n.nodes[i] = newNode(n.ctx, u, n.nodeChan, n.depth+1, n)
	}

}

func (n *Node) close() {
	n.mx.Lock()
	defer n.mx.Unlock()

	if n.isClosed {
		return
	}
	n.isClosed = true

	if n.i < n.size {
		for range n.nodeChan {
			n.i++
			if n.i == n.size {
				break
			}
		}
	}

	close(n.nodeChan)
	n.doneChan <- struct{}{}
}

func (n *Node) wait() {
	defer n.close()

	for {
		select {
		case <-n.ctx.Done():
			return

		case _, ok := <-n.nodeChan:
			if !ok {
				return
			}

			n.i++
			if n.i == n.size {
				return
			}
		}
	}
}

func (n Node) String() string {
	return n.url.String()
}
