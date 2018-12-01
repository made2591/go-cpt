package predictionTree

import (
	"strings"
)

const (
	newLine      = "\n"
	emptySpace   = "    "
	middleItem   = "├── "
	continueItem = "│   "
	lastItem     = "└── "
)

// The struct PredictionTree
type PredictionTree struct {

	Item     string
	Parent   *PredictionTree
	Children []*PredictionTree

}

// NewPredictionTree creates a new PredictionTree
func NewPredictionTree(value string) (predictionTree *PredictionTree) {

	// setup PredictionTree
	predictionTree = &PredictionTree{}

	// setup prediction tree
	predictionTree.Item = value
	predictionTree.Children = make([]*PredictionTree, 0)

	return predictionTree

}

// AddChild append a Child to PredictionTree
func AddChild(node *PredictionTree, value string) *PredictionTree {

	// if prediction tree is nil
	if node == nil {
		return nil
	}

	// create new child
	newChild := NewPredictionTree(value)
	newChild.Parent = node
	node.Children = append(node.Children, newChild)

	return newChild

}

// IsLeaf appends PredictionTree
func IsLeaf(node *PredictionTree) bool {
	return len(node.Children) == 0
}

// GetChildWithValue provides the pointer of the direct child with the given value, if any
func GetChildWithValue(node *PredictionTree, value string) (found bool, predictionTree *PredictionTree) {

	// if the node is nil, or has no children
	if node == nil || node.Children == nil || len(node.Children) == 0 {
		return false, nil
	}

	// return child if found
	for _, child := range node.Children {
		if strings.EqualFold(child.Item, value) {
			return true, child
		}
	}

	return false, nil

}

// String produce a string of PredictionTree
func String(node *PredictionTree) string {

	return node.Item + newLine + printItems(node.Children, []bool{})

}

func printText(text string, spaces []bool, last bool) string {

	var result string

	for _, space := range spaces {
		if space {
			result += emptySpace
		} else {
			result += continueItem
		}
	}

	indicator := middleItem
	if last {
		indicator = lastItem
	}

	return result + indicator + text + newLine

}

func printItems(t []*PredictionTree, spaces []bool) string {

	var result string

	for i, f := range t {
		last := i == len(t)-1
		result += printText(f.Item, spaces, last)
		if len(f.Children) > 0 {
			spacesChild := append(spaces, last)
			result += printItems(f.Children, spacesChild)
		}
	}

	return result

}
