package issummary

import (
	"fmt"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type LabelGraph struct {
	g            *simple.DirectedGraph
	gMap         map[int64]*labelNode
	labelNodeMap map[int64]*labelNode
	labelNameMap map[string]*labelNode
}

func newLabelGraph() *LabelGraph {
	return &LabelGraph{
		g:            simple.NewDirectedGraph(),
		gMap:         map[int64]*labelNode{},
		labelNodeMap: map[int64]*labelNode{},
		labelNameMap: map[string]*labelNode{}, // FIXME labelNameMap does not work well if duplicated labels exist beyond group
	}
}

type labelNode struct {
	node  graph.Node
	label *Label
}

func (lg *LabelGraph) AddLabel(label *Label) bool {
	if label == nil {
		return false
	}

	labelId := label.GetID()

	if _, ok := lg.gMap[labelId]; ok {
		return false
	}

	lNode := lg.g.NewNode()
	lg.g.AddNode(lNode)
	lg.gMap[lNode.ID()] = &labelNode{node: lNode, label: label}
	lg.labelNodeMap[labelId] = &labelNode{node: lNode, label: label}
	lg.labelNameMap[label.GetName()] = &labelNode{node: lNode, label: label}

	return true
}

func (lg *LabelGraph) getLabelNodes() (labelNodes []*labelNode) {
	for _, labelNode := range lg.gMap {
		labelNodes = append(labelNodes, labelNode)
	}
	return
}

func (lg *LabelGraph) SetEdges() error {
	for _, labelNode := range lg.getLabelNodes() {
		if err := lg.setEdgeByParent(labelNode.label); err != nil {
			return err
		}
	}

	if err := lg.setParents(); err != nil {
		return err
	}

	return nil
}

func (lg *LabelGraph) setEdgeByParent(label *Label) error {
	parentLabelName := label.Description.ParentName
	if parentLabelName == "" {
		return nil
	}

	parentLabelNode, ok := lg.labelNameMap[parentLabelName]
	if !ok {
		return fmt.Errorf("parent label %q not found", parentLabelName)
	}

	parentLabel := parentLabelNode.label
	if err := lg.SetEdge(label, parentLabel); err != nil {
		return err
	}

	return nil
}

func (lg *LabelGraph) getLabelByNodeID(id int64) (*Label, bool) {
	labelNode, ok := lg.gMap[id]
	if !ok {
		return nil, false
	}
	return labelNode.label, ok
}

func (lg *LabelGraph) getLabelNodeByName(labelName string) (*labelNode, bool) {
	labelNode, ok := lg.labelNameMap[labelName]
	return labelNode, ok
}

func (lg *LabelGraph) toLabelsFromNodes(nodes []graph.Node) (labels []*Label, err error) {
	for _, node := range nodes {
		label, ok := lg.getLabelByNodeID(node.ID())
		if !ok {
			return nil, fmt.Errorf("label node node %v not found", node.ID())
		}
		labels = append(labels, label)
	}
	return
}

func (lg *LabelGraph) SetEdge(aLabel, bLabel *Label) error {
	aLabelNode, ok := lg.getLabelNodeByName(aLabel.GetName())
	if !ok {
		return fmt.Errorf("label %q not found", aLabel.GetName())
	}

	bWorkNode, ok := lg.getLabelNodeByName(bLabel.GetName())
	if !ok {
		return fmt.Errorf("label %q not found", bLabel.GetName())
	}

	edge := lg.g.NewEdge(aLabelNode.node, bWorkNode.node)
	lg.g.SetEdge(edge)
	return nil
}

func (lg *LabelGraph) list() (labels []*Label) {
	for _, labelNode := range lg.getLabelNodes() {
		labels = append(labels, labelNode.label)
	}
	return
}

func (lg *LabelGraph) listParents(label *Label) ([]*Label, error) {
	labelNode, ok := lg.getLabelNodeByName(label.GetName())
	if !ok {
		return nil, fmt.Errorf("label %q not found in graph when list parents", label.GetName())
	}

	nodes := lg.g.From(labelNode.node.ID())
	labels, err := lg.toLabelsFromNodes(nodes)
	return labels, errors.Wrap(err, fmt.Sprintf("failed to get labels which depended from %v from label graph\n", label.GetName()))
}

func (lg *LabelGraph) setParents() error {
	for _, label := range lg.list() {
		parentLabels, err := lg.listParents(label)
		if err != nil {
			return err
		}
		label.ParentLabels = parentLabels
	}
	return nil
}
