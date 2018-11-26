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

type PredictionTree struct{
	Item string
	Parent *PredictionTree
	Children[] *PredictionTree
}

func NewPredictionTree(value string) (predictionTree *PredictionTree) {
	predictionTree = &PredictionTree{}
	predictionTree.Item = value
	predictionTree.Children = make([]*PredictionTree, 0)
	return predictionTree
}

func AddChild(node *PredictionTree, value string) *PredictionTree {
	if node == nil {
		return nil
	}
	newChild := &PredictionTree{Item: value}
	node.Children = append(node.Children, newChild)
	return newChild
}

func IsLeaf(node *PredictionTree) bool {
	return len(node.Children) == 0
}

func GetChildWithValue(node *PredictionTree, value string) (found bool, predictionTree *PredictionTree) {
	if node == nil || node.Children == nil || len(node.Children) == 0 {
		return false, nil
	}
	for _, c := range node.Children {
		if strings.EqualFold(c.Item, value) {
			return true, c
		}
	}
	return false, nil
}
	
func GetAllChildren(node *PredictionTree) []*PredictionTree {
	return node.Children
}

func RemoveChildWithValue(node *PredictionTree, value string) bool {
	found := false
	newNodes := make([]*PredictionTree, 0)
	for _, c := range node.Children {
		if !strings.EqualFold(c.Item, value) {
			newNodes = append(newNodes, c)
			found = true
		}
	}
	if found {
		node.Children = newNodes
	}
	return found
}

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