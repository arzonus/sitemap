package node

import (
	"bytes"
	"fmt"
	"io"
)

const (
	edgeTypeLink = "│"
	edgeTypeMid  = "├──"
	edgeTypeLast = "└──"
)

func (n *Node) Tree() string {
	return string(n.TreeBytes())
}

func (n *Node) TreeBytes() []byte {
	buf := new(bytes.Buffer)
	level := 0
	var levelsEnded []int
	if n.parent == nil {
		buf.WriteString(fmt.Sprintf("%v", n.String()))
		buf.WriteByte('\n')
	} else {
		edge := edgeTypeMid
		if len(n.nodes) == 0 {
			edge = edgeTypeLast
			levelsEnded = append(levelsEnded, level)
		}
		printValues(buf, 0, levelsEnded, edge, n.String())
	}
	if len(n.nodes) > 0 {
		printNodes(buf, level, levelsEnded, n.nodes)
	}
	return buf.Bytes()
}

func printNodes(wr io.Writer,
	level int, levelsEnded []int, nodes []*Node) {

	for i, node := range nodes {
		edge := edgeTypeMid
		if i == len(nodes)-1 {
			levelsEnded = append(levelsEnded, level)
			edge = edgeTypeLast
		}
		printValues(wr, level, levelsEnded, edge, node.String())
		if len(node.nodes) > 0 {
			printNodes(wr, level+1, levelsEnded, node.nodes)
		}
	}
}

func printValues(wr io.Writer,
	level int, levelsEnded []int, edge string, val string) {

	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprint(wr, "    ")
			continue
		}
		fmt.Fprintf(wr, "%s   ", edgeTypeLink)
	}
	fmt.Fprintf(wr, "%s %v\n", edge, val)
}

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}
