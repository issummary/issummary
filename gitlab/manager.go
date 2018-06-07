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

type GraphManager struct {
	g           *simple.DirectedGraph
	gMap        map[int64]*WorkNode
	workNodeMap map[int]*WorkNode
}

type WorkNode struct {
	node graph.Node
	work *Work
}

func (gm *GraphManager) addWork(work *Work) {
	node := gm.g.NewNode()
	gm.g.AddNode(node)
	gm.gMap[node.ID()] = &WorkNode{node: node, work: work}
	gm.workNodeMap[work.Issue.ID] = &WorkNode{node: node, work: work}
}

func (gm *GraphManager) getWorkNodes() (workNodes []*WorkNode) {
	for _, workNode := range gm.gMap {
		workNodes = append(workNodes, workNode)
	}
	return
}

func (gm *GraphManager) getWorks() (works []*Work) {
	for _, workNode := range gm.getWorkNodes() {
		works = append(works, workNode.work)
	}
	return
}

func (gm *GraphManager) getWorkByNodeID(id int64) (*Work, bool) {
	workNode, ok := gm.gMap[id]
	return workNode.work, ok
}

func (gm *GraphManager) getWorkNodeByWork(work *Work) (*WorkNode, bool) {
	return gm.getWorkNodeByWorkID(work.Issue.ID)
}

func (gm *GraphManager) getWorkNodeByWorkID(id int) (*WorkNode, bool) {
	workNode, bool := gm.workNodeMap[id]
	return workNode, bool
}

func (gm *GraphManager) ConnectByDependencies() error {
	for _, fromWorkNode := range gm.workNodeMap {
		if err := gm.setEdgesByDependIssues(fromWorkNode, fromWorkNode.work.Dependencies.Issues); err != nil {
			return fmt.Errorf("work edge setting is failed: %v", err) // FIXME
		}
		for _, dependLabel := range fromWorkNode.work.Dependencies.Labels {
			if err := gm.setEdgesByDependIssues(fromWorkNode, dependLabel.RelatedIssues); err != nil {
				return fmt.Errorf("depend label edge setting is failed: %v", err) // FIXME
			}
		}

		if fromWorkNode.work.Label == nil {
			continue
		}

		for _, label := range fromWorkNode.work.Label.Dependencies {
			for _, issue := range label.RelatedIssues {
				fmt.Println(issue.Title)
			}
			if err := gm.setEdgesByDependIssues(fromWorkNode, label.RelatedIssues); err != nil {
				return fmt.Errorf("label dependency edge setting is failed: %v", err) // FIXME
			}
		}
	}
	return nil
}

func (gm *GraphManager) setEdgesByDependIssues(workNode *WorkNode, issues []*Issue) error {
	for _, issue := range issues {
		toID := issue.ID
		toWorkNode, ok := gm.workNodeMap[toID]
		if !ok {
			return fmt.Errorf("dependency (to: %v) cant resolve\n", toWorkNode.work) // FIXME
		}

		if workNode.node.ID() == toWorkNode.node.ID() {
			return fmt.Errorf("self edge: %s/%s", workNode.work.Issue.ProjectName, workNode.work.Issue.Title)
		}

		gm.g.SetEdge(gm.g.NewEdge(workNode.node, toWorkNode.node))
	}
	return nil
}

type Dependency struct {
	fromIID int
	toIID   int
}

type WorkManager struct {
	gm *GraphManager
}

func NewWorkManager() *WorkManager {
	gm := &GraphManager{
		g:           simple.NewDirectedGraph(),
		gMap:        map[int64]*WorkNode{},
		workNodeMap: map[int]*WorkNode{},
	}

	return &WorkManager{gm}
}

func (wg *WorkManager) AddWork(work *Work) {
	wg.gm.addWork(work)
}

func (wg *WorkManager) AddWorks(works []*Work) {
	for _, work := range works {
		wg.AddWork(work)
	}
}

func (wg *WorkManager) ListSortedWorksByDueDate() (workNodes []*WorkNode) {
	workNodes = wg.gm.getWorkNodes()
	sort.Slice(workNodes, func(i, j int) bool {
		return workNodes[i].work.Issue.DueDate.After(*(workNodes[j].work.Issue.DueDate))
	})
	return workNodes
}

func (gm *GraphManager) GetSortedWorks() (works []*Work, err error) {

	marshalGraph, err := dot.Marshal(gm.g, "name", "prefix", "  ", false)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile("test.dot", marshalGraph, 0777); err != nil {
		return nil, err
	}

	nodes, err := topo.SortStabilized(gm.g, func(nodes []graph.Node) {
		sort.Slice(nodes, func(i, j int) bool {
			aWork := gm.gMap[nodes[i].ID()]
			bWork := gm.gMap[nodes[j].ID()]

			if aWork.work.Label == nil {
				return true
			}

			if bWork.work.Label == nil {
				return false
			}

			return aWork.work.Label.Name < bWork.work.Label.Name
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := gm.gMap[nodes[i].ID()]
			bWork := gm.gMap[nodes[j].ID()]

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
			aWork := gm.gMap[nodes[i].ID()]
			bWork := gm.gMap[nodes[j].ID()]
			return aWork.work.Issue.ProjectName+string(aWork.work.Issue.IID) > bWork.work.Issue.ProjectName+string(bWork.work.Issue.IID)
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := gm.gMap[nodes[i].ID()]
			bWork := gm.gMap[nodes[j].ID()]

			if bWork.work.Issue.DueDate == nil {
				return false
			}

			if aWork.work.Issue.DueDate == nil {
				return true
			}

			return aWork.work.Issue.DueDate.After(*bWork.work.Issue.DueDate)
		})
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		works = append(works, gm.gMap[node.ID()].work)
	}

	return reverseWorks(works), nil

	//workNodeFlags := map[int64]struct{}{}
	//// TODO: 締め切りが設定されているworkを短い順に取り出す
	//for _, workNodes := range gm.ListSortedWorksByDueDate() {
	//	nodes := gm.g.To(workNodes.node.ID())
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

func (wg *WorkManager) ConnectByDependencies() error {
	return wg.gm.ConnectByDependencies()
}

func (wg *WorkManager) GetSortedWorks() (works []*Work, err error) {
	return wg.gm.GetSortedWorks()
}

func (wg *WorkManager) GetDependWorks(work *Work) (works []*Work) {
	workNode, ok := wg.gm.getWorkNodeByWork(work)
	if !ok {
		return
	}

	nodes := wg.gm.g.To(workNode.node.ID())
	for _, node := range nodes {
		if work, ok := wg.gm.getWorkByNodeID(node.ID()); ok {
			works = append(works, work)
		}
	}
	return
}

func reverseWorks(works []*Work) []*Work {
	if len(works) == 0 {
		return works
	}
	return append(reverseWorks(works[1:]), works[0])
}
