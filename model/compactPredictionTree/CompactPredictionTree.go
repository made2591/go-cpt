package compactPredictionTree

import (
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"
	"strings"
)

type CompactPredictionTree struct {

	invertedIndexTable *invertedIndexTable.InvertedIndexTable
	lookupTable *lookupTable.LookupTable
	predictionTree *predictionTree.PredictionTree

}

func NewCompactPredictionTree(
	invertedIndexTable *invertedIndexTable.InvertedIndexTable,
	lookupTable *lookupTable.LookupTable,
	predictionTree *predictionTree.PredictionTree) (compactPredictionTree *CompactPredictionTree) {
        compactPredictionTree = &CompactPredictionTree{}
		compactPredictionTree.invertedIndexTable = invertedIndexTable
		compactPredictionTree.lookupTable = lookupTable
		compactPredictionTree.predictionTree = predictionTree

	return compactPredictionTree
}

func InitCompactPredictionTree(compactPredictionTree *CompactPredictionTree, sequences []*sequence.Sequence) {

	cursorNode := compactPredictionTree.predictionTree
	invertedIndex := compactPredictionTree.invertedIndexTable
	lookup := compactPredictionTree.lookupTable
	for _, seq := range sequences {
		for index, elem := range seq.Values {
			if found, child := predictionTree.GetChildWithValue(cursorNode, elem); found == false {
				cursorNode = predictionTree.AddChild(cursorNode, elem)
			} else {
				cursorNode = child
			}
			invertedIndexTable.AddSequenceIfMissing(invertedIndex, elem, seq)
			if index == len(seq.Values)-1 {
				lookup.Table[seq.ID] = cursorNode
			}
		}
		cursorNode = compactPredictionTree.predictionTree
	}

}

func String(compactPredictionTree CompactPredictionTree) (result string) {
	ii := invertedIndexTable.String(compactPredictionTree.invertedIndexTable)
	lt := lookupTable.String(compactPredictionTree.lookupTable)
	pt := predictionTree.String(compactPredictionTree.predictionTree)
	result = strings.Join([]string{ii, lt, pt}, "\n\n")
	return result
}