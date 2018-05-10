package gitlab

import (
	"sort"

	"gonum.org/v1/gonum/graph/simple"
)

type WorkNode struct {
	node simple.Node
	work *Work
}

type WorkManager struct {
	g    *simple.DirectedGraph
	gMap map[int64]*WorkNode
}

func NewWorkManager() *WorkManager {
	return &WorkManager{
		g:    simple.NewDirectedGraph(),
		gMap: map[int64]*WorkNode{},
	}
}

func (wg *WorkManager) AddWork(work *Work) {
	workNode := wg.g.NewNode()
	wg.gMap[workNode.ID()] = &WorkNode{}
	wg.g.AddNode(workNode)
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
	workNodeFlags := map[int64]struct{}{}
	// TODO: 締め切りが設定されているworkを短い順に取り出す
	for _, workNodes := range wg.ListSortedWorksByDueDate() {
		nodes := wg.g.To(workNodes.node.ID())
		for _, node := range nodes {
			if _, ok := workNodeFlags[node.ID()]; !ok {
				continue
			}

			// TODO: スケジュールに追加する処理
			// TODO: 取り出されたworkごとに依存ソートして配置
			// TODO: 締め切りまでの利用可能日数を計算
			// TODO: 日数が足りなければ並列度を増やす

			workNodeFlags[node.ID()] = struct{}{}
		}
	}

	return
}
