package compactPredictionTree

import (
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"
	"strings"
	"fmt"
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

func PredictionOverTestingSequence(compactPredictionTree *CompactPredictionTree, targetSequence *sequence.Sequence) string {

	uniqueValues := sequence.UniqueElements(targetSequence)
	similarSequences := make([]*sequence.Sequence, 0)
	for _, uniqueValue := range uniqueValues {
		for _, seq := range compactPredictionTree.invertedIndexTable.Table[uniqueValue] {
			found := false
			for _, alreadySeq := range similarSequences {
				if sequence.EqualSequence(seq, alreadySeq) {
					found = true
					break
				}
			}
			if !found {
				similarSequences = append(similarSequences, seq)
			}
		}
	}

	fmt.Println(similarSequences)

	consequents := make(map[int]*sequence.Sequence, 0)
	for _, similarSequence := range similarSequences {
		consequents[similarSequence.ID] = sequence.ComputeConsequent(targetSequence, similarSequence)
	}

	countables := make(map[string]float64, 0)
	for _, consequent := range consequents {
		for _, elem := range consequent.Values {
			score := 0.0
			if score, ok := countables[elem]; !ok {
				score = float64(1 + (1/len(similarSequences)) +(1/len(countables)+1)) * 0.001
			} else {
				score = score * float64(1 + (1/len(similarSequences)) +(1/len(countables)+1)) * 0.001
			}
			countables[elem] = score
		}
	}

	max := -1.0
	result := ""
	for elem, value := range countables {
		if value > max {
			result = elem
		}
	}
	return result

}

func String(compactPredictionTree CompactPredictionTree) (result string) {
	ii := invertedIndexTable.String(compactPredictionTree.invertedIndexTable)
	lt := lookupTable.String(compactPredictionTree.lookupTable)
	pt := predictionTree.String(compactPredictionTree.predictionTree)
	result = strings.Join([]string{ii, lt, pt}, "\n\n")
	return result
}