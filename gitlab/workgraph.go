package gitlab // FIXME

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type WorkGraph struct {
	g           *simple.DirectedGraph
	gMap        map[int64]*WorkNode
	workNodeMap map[int]*WorkNode
}

type WorkNode struct {
	node graph.Node
	work *Work
}

func (w *WorkGraph) AddWork(work *Work) {
	node := w.g.NewNode()
	w.g.AddNode(node)
	w.gMap[node.ID()] = &WorkNode{node: node, work: work}
	w.workNodeMap[work.Issue.ID] = &WorkNode{node: node, work: work}
}

func (w *WorkGraph) getWorkNodes() (workNodes []*WorkNode) {
	for _, workNode := range w.gMap {
		workNodes = append(workNodes, workNode)
	}
	return
}

func (w *WorkGraph) GetWorks() (works []*Work) {
	for _, workNode := range w.getWorkNodes() {
		works = append(works, workNode.work)
	}
	return
}

func (w *WorkGraph) getWorkByNodeID(id int64) (*Work, bool) {
	workNode, ok := w.gMap[id]
	return workNode.work, ok
}

func (w *WorkGraph) toWorkNode(work *Work) (*WorkNode, bool) {
	return w.getWorkNodeByID(work.Issue.ID)
}

func (w *WorkGraph) getWorkNodeByID(id int) (*WorkNode, bool) {
	workNode, ok := w.workNodeMap[id]
	return workNode, ok
}

func (w *WorkGraph) GetWorkByID(id int) (*Work, bool) {
	if workNode, ok := w.getWorkNodeByID(id); ok {
		return workNode.work, true
	}
	return nil, false
}

func (w *WorkGraph) SetEdge(aWork, bWork *Work) error {
	aWorkNode, ok := w.toWorkNode(aWork)
	if !ok {
		return fmt.Errorf("work %v not found", aWork)
	}

	bWorkNode, ok := w.toWorkNode(bWork)
	if !ok {
		return fmt.Errorf("work %v not found", bWork)
	}

	w.g.SetEdge(w.g.NewEdge(aWorkNode.node, bWorkNode.node))
	return nil
}
