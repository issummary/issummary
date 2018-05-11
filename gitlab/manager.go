package gitlab

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/graph"
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
	workMap     map[int]*WorkNode
	depIIDQueue []*Dependency
}

func NewWorkManager() *WorkManager {
	return &WorkManager{
		g:    simple.NewDirectedGraph(),
		gMap: map[int64]*WorkNode{},
	}
}

func (wg *WorkManager) AddWork(work *Work) {
	node := wg.g.NewNode()
	wg.gMap[node.ID()] = &WorkNode{node: node, work: work}
	wg.workMap[work.Issue.IID] = &WorkNode{node: node, work: work}
	// 依存関係をqueueに入れておく
	for _, issue := range work.Dependencies.Issues {
		wg.depIIDQueue = append(wg.depIIDQueue, &Dependency{fromIID: work.Issue.IID, toIID: issue.IID})
	}
}

func (wg *WorkManager) ConnectByDependencies() error {
	for _, dep := range wg.depIIDQueue {
		fromWorkNode, ok := wg.workMap[dep.fromIID]
		if !ok {
			return fmt.Errorf("dependency (from: %s) cant resolve\n", dep) // FIXME
		}
		toWorkNode, ok := wg.workMap[dep.toIID]
		if !ok {
			return fmt.Errorf("dependency (to: %s) cant resolve\n", dep) // FIXME
		}

		wg.g.Edge(fromWorkNode.node.ID(), toWorkNode.node.ID())
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
	nodes, err := topo.Sort(wg.g)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		works = append(works, wg.gMap[node.ID()].work)
	}

	return

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
