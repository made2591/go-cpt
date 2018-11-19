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
	Children *PredictionTree
}

func NewPredictionTree(value string) (predictionTree *PredictionTree) {
	predictionTree = &PredictionTree{}
	predictionTree.Item = value
	return predictionTree
}

func AddChild(node *PredictionTree, value string) *PredictionTree {
	pointer := node
	for pointer.Children != nil {
		pointer = pointer.Children
	}
	pointer.Children = &PredictionTree{Item: value}
	return pointer.Children
}

func GetChildWithValue(node *PredictionTree, value string) (found bool, predictionTree *PredictionTree){
	pointer := node
	for pointer.Children != nil {
		if strings.EqualFold(pointer.Children.Item, value) {
			return true, pointer.Children
		}
		pointer = pointer.Children
	}
	return false, nil
}
	
func GetAllChildren(node *PredictionTree) *PredictionTree {
	return node.Children
}

func RemoveChildWithValue(node *PredictionTree, value string) bool {
	pointer := node
	for pointer.Children != nil {
		if strings.EqualFold(pointer.Children.Item, value) {
			pointer.Children = pointer.Children.Children
			return true
		}
		pointer = pointer.Children
	}
	return false
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

func printItems(t *PredictionTree, spaces []bool) string {
	var result string
	pointer := t
	for pointer.Children != nil {
		last := pointer.Children.Children == nil
		result += printText(pointer.Item, spaces, last)
		if pointer.Children.Children != nil {
			spacesChild := append(spaces, last)
			result += printItems(pointer.Children.Children, spacesChild)
		}
		pointer = pointer.Children
	}
	return result
}