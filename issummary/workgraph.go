package issummary // FIXME

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type WorkGraph struct {
	g           *simple.DirectedGraph
	gMap        map[int64]*WorkNode
	workNodeMap map[int]*WorkNode
	edgeMap     map[graph.Edge]string
}

type WorkNode struct {
	node graph.Node
	work *Work
}

type SortWorkFunc func(aWork, bWork *Work) bool

type ListWorksOptions struct {
	GroupName   string
	ProjectName string
	Number      int
	LabelNames  []string
}

func (w *WorkGraph) AddWork(work *Work) {
	node := w.g.NewNode()
	w.g.AddNode(node)
	w.gMap[node.ID()] = &WorkNode{node: node, work: work}
	w.workNodeMap[int(work.Issue.GetID())] = &WorkNode{node: node, work: work} // FIXME
}

func (w *WorkGraph) getWorkNodes() (workNodes []*WorkNode) {
	for _, workNode := range w.gMap {
		workNodes = append(workNodes, workNode)
	}
	return
}

func (w *WorkGraph) ListWorks(opt *ListWorksOptions) (works []*Work) {
	for _, workNode := range w.getWorkNodes() {
		work := workNode.work
		if opt != nil {
			if opt.GroupName != "" && opt.GroupName != work.Issue.GroupName {
				continue
			}
			if opt.ProjectName != "" && opt.ProjectName != work.Issue.ProjectName {
				continue
			}
			if opt.Number != 0 && opt.Number != work.Issue.GetNumber() {
				continue
			}

			if !work.Issue.HasAllLabels(opt.LabelNames) {
				continue
			}
		}

		works = append(works, workNode.work)
	}
	return
}

func (w *WorkGraph) GetSortedWorks(sortWorkFunctions []SortWorkFunc) (works []*Work, err error) {
	nodes, err := topo.SortStabilized(w.g, func(nodes []graph.Node) {
		for _, sortWorkFunction := range sortWorkFunctions {
			sortFunction := func(i, j int) bool {
				aWork, _ := w.getWorkByNodeID(nodes[i].ID())
				bWork, _ := w.getWorkByNodeID(nodes[j].ID())
				return sortWorkFunction(aWork, bWork)
			}
			sort.SliceStable(nodes, sortFunction)
		}
	})
	if err != nil {
		return nil, err
	}

	works = w.convertNodesToWorks(nodes)
	return reverseWorks(works), nil
}

func (w *WorkGraph) getWorkNodeByNodeID(id int64) (*WorkNode, bool) {
	workNode, ok := w.gMap[id]
	return workNode, ok
}

func (w *WorkGraph) getWorkByNodeID(id int64) (*Work, bool) {
	workNode, ok := w.gMap[id]
	return workNode.work, ok
}

func (w *WorkGraph) toWorkNode(work *Work) (*WorkNode, bool) {
	return w.getWorkNodeByID(int(work.Issue.GetID()))
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

func (w *WorkGraph) convertNodesToWorks(nodes []graph.Node) (works []*Work) {
	for _, node := range nodes {
		if work, ok := w.getWorkByNodeID(node.ID()); ok {
			works = append(works, work)
		}
	}
	return
}

func (w *WorkGraph) SetEdge(aWork, bWork *Work, edgeName string) error {
	aWorkNode, ok := w.toWorkNode(aWork)
	if !ok {
		return fmt.Errorf("work %v not found", aWork)
	}

	bWorkNode, ok := w.toWorkNode(bWork)
	if !ok {
		return fmt.Errorf("work %v not found", bWork)
	}

	edge := w.g.NewEdge(aWorkNode.node, bWorkNode.node)
	w.g.SetEdge(edge)
	w.edgeMap[edge] = edgeName
	return nil
}

func (w *WorkGraph) getRelatedWorksByNodeID(id int64) (works []*Work) {
	nodes := w.g.From(id)

	beforeWorkNode, _ := w.getWorkNodeByNodeID(id) // FIXME
	for _, node := range nodes {
		if work, ok := w.getWorkByNodeID(node.ID()); ok {
			if edge := w.g.Edge(beforeWorkNode.node.ID(), node.ID()); edge != nil {
				work.WorkRelationType = NewWorkRelationTypeFromString(w.edgeMap[edge])
			}
			works = append(works, work)
		}
	}
	return works
}

func (w *WorkGraph) GetRelatedWorks(work *Work) (works []*Work) {
	workNode, ok := w.toWorkNode(work)
	if !ok {
		return
	}

	return w.getRelatedWorksByNodeID(workNode.node.ID())
}
