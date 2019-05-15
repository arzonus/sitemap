package node

import (
	"context"
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

	ii int
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
		n.close("ERROR RESULT")
		return
	}
	n.setNodes(result.URLs)
}

func (n *Node) setNodes(urls []*url.URL) {
	n.nodes = make([]*Node, len(urls))
	n.size = len(urls)
	log.Printf("set urls id %d, size %d, node count %d", n.id, n.size, len(n.nodes))

	for i, url := range urls {
		n.nodes[i] = newNode(n.ctx, url, n.nodeChan, n.depth+1, n)
	}

}

func (n *Node) close(str string) {
	n.mx.Lock()
	defer n.mx.Unlock()
	defer func() {
		if p := recover(); p != nil {
			log.Printf("PANIC: id %d, node count %d, parent %d, parent node count %d, %d", n.id, len(n.nodes), n.parent.id, len(n.parent.nodes), n.parent.ii)
			panic(p)
		}
	}()

	if n.isClosed {
		return
	}

	n.isClosed = true
	log.Printf("try to close chan, id %d, node count %d, parent %d, parent node count %d, i %d, , %d", n.id, len(n.nodes), n.parent.id, len(n.parent.nodes), n.ii, n.size)
	close(n.nodeChan)
	log.Printf("close chan, try to send, id %d, node count %d, parent %d, parent node count %d, i %d", n.id, len(n.nodes), n.parent.id, len(n.parent.nodes), n.ii)
	n.doneChan <- struct{}{}
	log.Printf("success closed id %d, node count %d, parent %d, parent node count %d", n.id, len(n.nodes), n.parent.id, len(n.parent.nodes))

}

func (n Node) sz() int {
	log.Printf("sz is:%d %d", n.id, n.size)
	return n.size
}

func (n *Node) wait() {

	var i int
	for {
		select {
		case <-n.ctx.Done():
			log.Printf("ctx done: %d %d %d", n.id, len(n.nodes), n.size)
			if n.sz() == 0 {
				log.Printf("0: ctx done: %d %d %d %d", n.id, len(n.nodes), n.size, n.sz())
				log.Printf("0: ctx done: %#v", n)
				n.close("ctx: 0")
				return
			}

			for range n.nodeChan {
				log.Printf("loop: ctx done: %d %d", n.id, len(n.nodes))
				i++
				if i == n.size {
					log.Printf("i: ctx done: %d, i: %d, count: %d", n.id, i, len(n.nodes))
					n.close("i: ctx")
					return
				}
			}

			log.Printf("NODE CHAN CLOSED id %d, node count %d, parent %d, parent node count %d, i %d", n.id, len(n.nodes), n.parent.id, len(n.parent.nodes), i)

			return

		case _, ok := <-n.nodeChan:

			if !ok {
				n.close("ok")
				return
			}

			i++
			log.Printf("node chan: %d, i: %d, count: %d", n.id, i, len(n.nodes))
			if i == n.size {
				log.Printf("success: node chan: %d, i: %d, count: %d", n.id, i, len(n.nodes))
				n.close("success")
				return
			}
		}
	}
}
