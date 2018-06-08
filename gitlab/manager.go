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

func (wg *WorkManager) ConnectByDependencies() error {
	for _, work := range wg.w.GetWorks() {
		wg.setEdgesByWork(work)
	}
	return nil
}

func (wg *WorkManager) setEdgesByWork(fromWork *Work) error {
	if err := wg.setEdgesByDependIssues(fromWork, fromWork.Dependencies.Issues); err != nil {
		return fmt.Errorf("work edge setting is failed: %v", err) // FIXME

	}

	for _, dependLabel := range fromWork.Dependencies.Labels {
		if err := wg.setEdgesByDependIssues(fromWork, dependLabel.RelatedIssues); err != nil {
			return fmt.Errorf("depend label edge setting is failed: %v", err) // FIXME
		}
	}

	if fromWork.Label == nil {
		return nil
	}

	for _, label := range fromWork.Label.Dependencies {
		if err := wg.setEdgesByDependIssues(fromWork, label.RelatedIssues); err != nil {
			return fmt.Errorf("label dependency edge setting is failed: %v", err) // FIXME
		}
	}

	return nil
}

func (wg *WorkManager) setEdgesByDependIssues(fromWork *Work, issues []*Issue) error {
	for _, issue := range issues {
		toWork, ok := wg.w.GetWorkByID(issue.ID)
		if !ok {
			return fmt.Errorf("dependency (to: %v) cant resolve\n", toWork) // FIXME
		}

		if fromWork.Issue.ID == toWork.Issue.ID {
			return fmt.Errorf("self edge: %s/%s", fromWork.Issue.ProjectName, fromWork.Issue.Title)
		}

		wg.w.SetEdge(fromWork, toWork)
	}
	return nil
}

type Dependency struct {
	fromIID int
	toIID   int
}

type WorkManager struct {
	w *WorkGraph
}

func NewWorkManager() *WorkManager {
	gm := &WorkGraph{
		g:           simple.NewDirectedGraph(),
		gMap:        map[int64]*WorkNode{},
		workNodeMap: map[int]*WorkNode{},
	}

	return &WorkManager{gm}
}

func (wg *WorkManager) AddWork(work *Work) {
	wg.w.AddWork(work)
}

func (wg *WorkManager) AddWorks(works []*Work) {
	for _, work := range works {
		wg.AddWork(work)
	}
}

func (wg *WorkManager) ListSortedWorksByDueDate() (workNodes []*WorkNode) {
	workNodes = wg.w.getWorkNodes()
	sort.Slice(workNodes, func(i, j int) bool {
		return workNodes[i].work.Issue.DueDate.After(*(workNodes[j].work.Issue.DueDate))
	})
	return workNodes
}

func (w *WorkGraph) GetSortedWorks() (works []*Work, err error) {

	marshalGraph, err := dot.Marshal(w.g, "name", "prefix", "  ", false)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile("test.dot", marshalGraph, 0777); err != nil {
		return nil, err
	}

	nodes, err := topo.SortStabilized(w.g, func(nodes []graph.Node) {
		sort.Slice(nodes, func(i, j int) bool {
			aWork := w.gMap[nodes[i].ID()]
			bWork := w.gMap[nodes[j].ID()]

			if aWork.work.Label == nil {
				return true
			}

			if bWork.work.Label == nil {
				return false
			}

			return aWork.work.Label.Name < bWork.work.Label.Name
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := w.gMap[nodes[i].ID()]
			bWork := w.gMap[nodes[j].ID()]

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
			aWork := w.gMap[nodes[i].ID()]
			bWork := w.gMap[nodes[j].ID()]
			return aWork.work.Issue.ProjectName+string(aWork.work.Issue.IID) > bWork.work.Issue.ProjectName+string(bWork.work.Issue.IID)
		})

		sort.SliceStable(nodes, func(i, j int) bool {
			aWork := w.gMap[nodes[i].ID()]
			bWork := w.gMap[nodes[j].ID()]

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
		works = append(works, w.gMap[node.ID()].work)
	}

	return reverseWorks(works), nil

	//workNodeFlags := map[int64]struct{}{}
	//// TODO: 締め切りが設定されているworkを短い順に取り出す
	//for _, workNodes := range w.ListSortedWorksByDueDate() {
	//	nodes := w.g.To(workNodes.node.ID())
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

func (wg *WorkManager) GetSortedWorks() (works []*Work, err error) {
	return wg.w.GetSortedWorks()
}

func (wg *WorkManager) GetDependWorks(work *Work) (works []*Work) {
	workNode, ok := wg.w.toWorkNode(work)
	if !ok {
		return
	}

	nodes := wg.w.g.To(workNode.node.ID())
	for _, node := range nodes {
		if work, ok := wg.w.getWorkByNodeID(node.ID()); ok {
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
