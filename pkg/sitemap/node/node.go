package node

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

var iterator int64

type Node struct {
	id int64

	parent *Node

	doneChan chan<- struct{}
	nodeChan chan struct{}

	ctx   context.Context
	depth int
	url   *url.URL

	nodes  []*Node
	size   int
	result *Result

	mx       sync.Mutex
	sxmx     sync.Mutex
	isClosed bool

	i int
}

func NewNode(
	ctx context.Context,
	url *url.URL,
	doneChan chan<- struct{},
) *Node {
	return newNode(ctx, url, doneChan, 0, &Node{})
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
		id:       atomic.AddInt64(&iterator, 1),
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
		n.close()
		return
	}
	n.setNodes(result.URLs)
}

func (n *Node) setNodes(urls []*url.URL) {
	n.mx.Lock()
	defer n.mx.Unlock()
	if n.isClosed {
		log.Println("IS CLOSED!")
		return
	}

	n.nodes = make([]*Node, len(urls))
	n.size = len(urls)
	log.Printf("set urls id %d, size %d, node count %d", n.id, n.size, len(n.nodes))

	for i, url := range urls {
		n.nodes[i] = newNode(n.ctx, url, n.nodeChan, n.depth+1, n)
	}

}

func (n *Node) close() {
	n.mx.Lock()
	defer n.mx.Unlock()
	defer func() {
		if p := recover(); p != nil {
			log.Println("PANIC:", n.id, n.parent.id)
			panic(p)
		}
	}()

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
	if n.result != nil {
		if n.result.Error != nil {
			return fmt.Sprintf("%s %s err: %s", n.Prefix(), n.url, n.result.Error)
		}
	} else {
		return fmt.Sprintf("%s %s wasn't handled", n.Prefix(), n.url)
	}

	var str = fmt.Sprintf("%s %s", n.Prefix(), n.url)
	for i := range n.nodes {
		str += "\n" + n.nodes[i].String()
	}
	return str
}
