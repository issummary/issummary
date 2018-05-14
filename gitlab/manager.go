package gitlab

import (
	"fmt"
	"io/ioutil"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type WorkNode struct {
	node graph.Node
	work *Work
}

type Dependency struct {
	fromIID int
	toIID   int
}

type WorkManager struct {
	g           *simple.DirectedGraph
	gMap        map[int64]*WorkNode
	workNodeMap map[int]*WorkNode
}

func NewWorkManager() *WorkManager {
	return &WorkManager{
		g:           simple.NewDirectedGraph(),
		gMap:        map[int64]*WorkNode{},
		workNodeMap: map[int]*WorkNode{},
	}
}

func (wg *WorkManager) AddWork(work *Work) {
	node := wg.g.NewNode()
	wg.g.AddNode(node)
	wg.gMap[node.ID()] = &WorkNode{node: node, work: work}
	wg.workNodeMap[work.Issue.ID] = &WorkNode{node: node, work: work}
}

func (wg *WorkManager) ConnectByDependencies() error {
	for _, fromWorkNode := range wg.workNodeMap {
		wg.setEdgesByDependIssues(fromWorkNode, fromWorkNode.work.Dependencies.Issues)
		for _, dependLabel := range fromWorkNode.work.Dependencies.Labels {
			wg.setEdgesByDependIssues(fromWorkNode, dependLabel.RelatedIssues)
		}
	}
	return nil
}

func (wg *WorkManager) setEdgesByDependIssues(workNode *WorkNode, issues []*Issue) error {
	for _, issue := range issues {
		toID := issue.ID
		toWorkNode, ok := wg.workNodeMap[toID]
		if !ok {
			return fmt.Errorf("dependency (to: %v) cant resolve\n", toWorkNode.work) // FIXME
		}
		wg.g.SetEdge(wg.g.NewEdge(workNode.node, toWorkNode.node))
	}
	return nil
}

func (wg *WorkManager) AddWorks(works []*Work) {
	for _, work := range works {
		wg.AddWork(work)
	}
}

func (wg *WorkManager) ListSortedWorksByDueDate() (workNodes []*WorkNode) {
	for _, work := range wg.gMap {
		workNodes = append(workNodes, work)
	}

	sort.Slice(workNodes, func(i, j int) bool {
		return workNodes[i].work.Issue.DueDate.After(*(workNodes[j].work.Issue.DueDate))
	})
	return workNodes
}

func (wg *WorkManager) GetSortedWorks() (works []*Work, err error) {

	marshalGraph, err := dot.Marshal(wg.g, "name", "prefix", "  ", false)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile("test.dot", marshalGraph, 0777); err != nil {
		return nil, err
	}

	nodes, err := topo.SortStabilized(wg.g, func(nodes []graph.Node) {
		sort.Slice(nodes, func(i, j int) bool {
			aWork := wg.gMap[nodes[i].ID()]
			bWork := wg.gMap[nodes[j].ID()]

			if aWork.work.Label == nil {
				return true
			}

			if bWork.work.Label == nil {
				return false
			}

			return aWork.work.Label.Name < bWork.work.Label.Name
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := wg.gMap[nodes[i].ID()]
			bWork := wg.gMap[nodes[j].ID()]

			if aWork.work.Label == nil {
				return true
			}

			if bWork.work.Label == nil {
				return false
			}

			if aWork.work.Label.Parent == nil {
				return true
			}

			if bWork.work.Label.Parent == nil {
				return false
			}

			return aWork.work.Label.Parent.Name < bWork.work.Label.Parent.Name
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := wg.gMap[nodes[i].ID()]
			bWork := wg.gMap[nodes[j].ID()]
			return aWork.work.Issue.ProjectName+string(aWork.work.Issue.IID) > bWork.work.Issue.ProjectName+string(bWork.work.Issue.IID)
		})
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		works = append(works, wg.gMap[node.ID()].work)
	}

	return reverseWorks(works), nil

	//workNodeFlags := map[int64]struct{}{}
	//// TODO: 締め切りが設定されているworkを短い順に取り出す
	//for _, workNodes := range wg.ListSortedWorksByDueDate() {
	//	nodes := wg.g.To(workNodes.node.ID())
	//	for _, node := range nodes {
	//		if _, ok := workNodeFlags[node.ID()]; !ok {
	//			continue
	//		}
	//
	//		// TODO: スケジュールに追加する処理
	//		// TODO: 取り出されたworkごとに依存ソートして配置
	//		// TODO: 締め切りまでの利用可能日数を計算
	//		// TODO: 日数が足りなければ並列度を増やす
	//
	//		workNodeFlags[node.ID()] = struct{}{}
	//	}
	//}
	//
	//return
}

func reverseWorks(works []*Work) []*Work {
	if len(works) == 0 {
		return works
	}
	return append(reverseWorks(works[1:]), works[0])
}
