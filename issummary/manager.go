package issummary

import (
	"fmt"
	"io/ioutil"
	"sort"

	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

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

func (wg *WorkManager) ResolveDependencies() error {
	for _, work := range wg.w.GetWorks() {
		wg.setEdgesByWork(work)
	}

	wg.setWorkDependencies()

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

func (wg *WorkManager) GetSortedWorks() (works []*Work, err error) {
	sortWorkFunctions := []SortWorkFunc{
		func(aWork, bWork *Work) bool {
			if aWork.Label == nil {
				return true
			}

			if bWork.Label == nil {
				return false
			}

			return aWork.Label.Name < bWork.Label.Name
		},
		func(aWork, bWork *Work) bool {
			if aWork.Label == nil {
				return true
			}

			if bWork.Label == nil {
				return false
			}

			if aWork.Label.Parent == nil {
				return true
			}

			if bWork.Label.Parent == nil {
				return false
			}

			return aWork.Label.Parent.Name < bWork.Label.Parent.Name
		},
		func(aWork, bWork *Work) bool {
			return aWork.Issue.ProjectName+string(aWork.Issue.IID) > bWork.Issue.ProjectName+string(bWork.Issue.IID)
		},
		func(aWork, bWork *Work) bool {
			if bWork.Issue.DueDate == nil {
				return false
			}

			if aWork.Issue.DueDate == nil {
				return true
			}

			return aWork.Issue.DueDate.After(*bWork.Issue.DueDate)
		},
	}

	return wg.w.GetSortedWorks(sortWorkFunctions)

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

func (wg *WorkManager) MarshalGraph() error {
	marshalGraph, err := dot.Marshal(wg.w.g, "name", "prefix", "  ", false)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile("test.dot", marshalGraph, 0777); err != nil {
		return err
	}
	return nil
}

func (wg *WorkManager) GetDependWorks(work *Work) (works []*Work) {
	return wg.w.GetRelatedWorks(work)
}

func (wg *WorkManager) setWorkDependencies() {
	for _, work := range wg.w.GetWorks() {
		work.DependWorks = wg.GetDependWorks(work)
		work.TotalStoryPoint = work.GetTotalStoryPoint()
	}
}

func reverseWorks(works []*Work) []*Work {
	if len(works) == 0 {
		return works
	}
	return append(reverseWorks(works[1:]), works[0])
}
