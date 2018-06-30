package issummary

import (
	"fmt"
	"io/ioutil"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type WorkManager struct {
	w *WorkGraph
}

func NewWorkManager() *WorkManager {
	gm := &WorkGraph{
		g:           simple.NewDirectedGraph(),
		gMap:        map[int64]*WorkNode{},
		workNodeMap: map[int]*WorkNode{},
		edgeMap:     map[graph.Edge]string{},
	}

	return &WorkManager{gm}
}

func (wg *WorkManager) ResolveDependencies() error {
	for _, work := range wg.w.ListWorks(nil) {
		wg.setEdgesByWork(work)
	}

	wg.setWorkDependencies()

	return nil
}

func (wg *WorkManager) setEdgesByWork(fromWork *Work) error {

	issueDependencies := fromWork.Issue.ListDependencies()
	for _, issue := range issueDependencies.Issues {
		opt := &ListWorksOptions{
			ProjectName: issue.ProjectName,
			GroupName:   issue.GroupName,
			Number:      issue.Number,
		}
		works := wg.w.ListWorks(opt)

		if len(works) == 0 {
			return fmt.Errorf("depend issue not found: %v", opt)
		}

		wg.w.SetEdge(fromWork, works[0], IssueOfIssueDescriptionRelation.String())
	}

	for _, labelName := range issueDependencies.LabelNames {
		opt := &ListWorksOptions{
			LabelNames: []string{labelName},
		}
		works := wg.w.ListWorks(opt)

		for _, work := range works {
			wg.w.SetEdge(fromWork, work, LabelOfIssueDescriptionRelation.String())
		}
	}

	if fromWork.Label == nil {
		return nil
	}

	for _, dependLabelName := range fromWork.Label.Description.DependLabelNames {
		opt := &ListWorksOptions{
			LabelNames: []string{dependLabelName},
		}
		works := wg.w.ListWorks(opt)
		for _, work := range works {
			wg.w.SetEdge(fromWork, work, LabelOfLabelDescriptionRelation.String())
		}
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

func (wg *WorkManager) GetListSortedWorksByDueDate() (workNodes []*WorkNode) {
	workNodes = wg.w.getWorkNodes()
	sort.Slice(workNodes, func(i, j int) bool {
		return workNodes[i].work.Issue.GetDueDate().After(*(workNodes[j].work.Issue.GetDueDate()))
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

			return aWork.Label.GetName() < bWork.Label.GetName()
		},
		func(aWork, bWork *Work) bool {
			if aWork.Label == nil {
				return true
			}

			if bWork.Label == nil {
				return false
			}

			if aWork.Label.ParentName == "" {
				return true
			}

			if bWork.Label.ParentName == "" {
				return false
			}

			return aWork.Label.ParentName < bWork.Label.ParentName
		},
		func(aWork, bWork *Work) bool {
			return aWork.Issue.ProjectName+string(aWork.Issue.GetNumber()) > bWork.Issue.ProjectName+string(bWork.Issue.GetNumber())
		},
		func(aWork, bWork *Work) bool {
			if bWork.Issue.GetDueDate() == nil {
				return false
			}

			if aWork.Issue.GetDueDate() == nil {
				return true
			}

			return aWork.Issue.GetDueDate().After(*bWork.Issue.GetDueDate())
		},
	}

	return wg.w.GetSortedWorks(sortWorkFunctions)

	//workNodeFlags := map[int64]struct{}{}
	//// TODO: 締め切りが設定されているworkを短い順に取り出す
	//for _, workNodes := range wg.GetListSortedWorksByDueDate()() {
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
	for _, work := range wg.w.ListWorks(nil) {
		work.DependWorks = wg.GetDependWorks(work)
		work.TotalStoryPoint = work.GetTotalStoryPoint()
	}
}

func (wg *WorkManager) listDependencies() (dependencies *IssueDependencies) {
	dependencies = &IssueDependencies{
		Issues:     []*DependIssue{},
		LabelNames: []string{},
	}
	for _, work := range wg.w.ListWorks(nil) {
		dep := work.Issue.Description.Dependencies
		dependencies.Issues = append(dependencies.Issues, dep.Issues...)
		dependencies.LabelNames = append(dependencies.LabelNames, dep.LabelNames...)
	}
	return
}

func (wg *WorkManager) ListMissingIssueDependencies() (dependIssues []*DependIssue) {
	dependencies := wg.listDependencies()
	for _, issue := range dependencies.Issues {
		works := wg.w.ListWorks(&ListWorksOptions{
			GroupName:   issue.GroupName,
			ProjectName: issue.ProjectName,
			Number:      issue.Number,
		})
		if len(works) == 0 {
			dependIssues = append(dependIssues, issue)
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
