package compactPredictionTree

import (
	"fmt"
	"sort"
	"strings"

	set "github.com/deckarep/golang-set"
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"
)

type CompactPredictionTree struct {
	InvertedIndexTable *invertedIndexTable.InvertedIndexTable
	LookupTable        *lookupTable.LookupTable
	PredictionTree     *predictionTree.PredictionTree
	TrainingSet        []*sequence.Sequence
	TestingSet         []*sequence.Sequence
}

func NewCompactPredictionTree(
	invertedIndexTable *invertedIndexTable.InvertedIndexTable,
	lookupTable *lookupTable.LookupTable,
	predictionTree *predictionTree.PredictionTree,
	trainingSet []*sequence.Sequence,
	testingSet []*sequence.Sequence) (compactPredictionTree *CompactPredictionTree) {
	compactPredictionTree = &CompactPredictionTree{}
	compactPredictionTree.InvertedIndexTable = invertedIndexTable
	compactPredictionTree.LookupTable = lookupTable
	compactPredictionTree.PredictionTree = predictionTree
	compactPredictionTree.TrainingSet = trainingSet
	compactPredictionTree.TestingSet = testingSet

	return compactPredictionTree
}

func InitCompactPredictionTree(compactPredictionTree *CompactPredictionTree, sequences []*sequence.Sequence) {

	cursorNode := compactPredictionTree.PredictionTree
	invertedIndex := compactPredictionTree.InvertedIndexTable
	lookup := compactPredictionTree.LookupTable
	for _, seq := range sequences {
		for index, elem := range seq.Values {
			if found, child := predictionTree.GetChildWithValue(cursorNode, elem); found == false {
				cursorNode = predictionTree.AddChild(cursorNode, elem)
			} else {
				cursorNode = child
			}
			invertedIndexTable.AddSequenceIfMissing(invertedIndex, elem, seq)
			if index == len(seq.Values)-1 {
				lookup.Table[seq.ID] = seq
			}
		}
		cursorNode = compactPredictionTree.PredictionTree
	}

}

func PredictionOverTestingSequence(compactPredictionTree *CompactPredictionTree, k int, n int) [][]string {

	results := make([][]string, 0)

	for _, targetSequence := range compactPredictionTree.TestingSet {

		intersection := set.NewSet()
		for _, targetSequence := range compactPredictionTree.TestingSet {
			intersection.Add(targetSequence.ID)
		}

		fmt.Println("original target: ", targetSequence.Values)
		if k < len(targetSequence.Values) {
			targetSequence = &sequence.Sequence{Values: targetSequence.Values[len(targetSequence.Values)-k:]}
		}
		fmt.Println("cut target: ", targetSequence.Values)

		for _, element := range targetSequence.Values {
			if len(compactPredictionTree.InvertedIndexTable.Table[element]) != 0 {
				fmt.Println("each target: ", element)
				fmt.Println("before target: ", intersection.String())
				seqID := set.NewSet()
				for _, seq := range compactPredictionTree.InvertedIndexTable.Table[element] {
					seqID.Add(seq.ID)
				}
				intersection = intersection.Intersect(seqID)
				fmt.Println("after target: ", intersection.String())
			}
		}

		similarSequences := make([]*sequence.Sequence, 0)
		it := intersection.Iterator()
		for element := range it.C {
			fmt.Println(element)
			currentNode := compactPredictionTree.LookupTable.Table[element.(int)]
			similarSequences = append(similarSequences, currentNode)
		}

		// similarSequences := make([]*sequence.Sequence, 0)

		fmt.Println("number of similar seqs: ", len(similarSequences))
		for _, seq := range similarSequences {
			fmt.Println("\t", sequence.String(seq))
		}

		consequents := make(map[int][]string, 0)
		for _, similarSequence := range similarSequences {
			consequent := sequence.ComputeConsequent(targetSequence, similarSequence)
			if len(consequent) > 0 {
				consequents[similarSequence.ID] = consequent
				// fmt.Println("\t", similarSequence.ID, consequents[similarSequence.ID])
			}
		}

		fmt.Println("consequents: ")
		countables := make(map[string]float64, 0)
		for _, consequent := range consequents {
			//fmt.Println(consequent)
			for _, elem := range consequent {
				countables = computeScore(countables, elem, len(targetSequence.Values),len(targetSequence.Values),len(similarSequences), len(countables)+1)
				fmt.Println("\t", elem, countables[elem])
			}
		}

		fmt.Println("number of counttable keys: ", len(countables))
		pairs := rankByScore(countables)
		result := make([]string, 0)

		if len(pairs) > 0 {
			for i := 0; i < n && i < len(pairs); i++ {
				result = append(result, pairs[i].Key)
				fmt.Println("\t", pairs[i].Key, pairs[i].Value)
			}
		}
		results = append(results, result)
		fmt.Println("results until now:", results)
	}
	return results

}

func computeScore(countables map[string]float64, key string, length int, target_size int, number_of_similar_sequences int, number_items_counttable int) map[string]float64 {

	weight_level := float64(1 / number_of_similar_sequences)
	weight_distance := float64(1 / number_items_counttable)
	score := 1.0 + weight_level + weight_distance * 0.001

	if oldScore, ok := countables[key]; !ok {
		countables[key] = score
	} else {
		countables[key] = score * oldScore
	}

	return countables
}


func rankByScore(scores map[string]float64) PairList {
	pl := make(PairList, len(scores))
	i := 0
	for k, v := range scores {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	fmt.Println("scores: ", scores)
	fmt.Println("pl: ", pl)
	return pl
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func String(compactPredictionTree CompactPredictionTree) (result string) {
	ii := invertedIndexTable.String(compactPredictionTree.InvertedIndexTable)
	lt := lookupTable.String(compactPredictionTree.LookupTable)
	pt := predictionTree.String(compactPredictionTree.PredictionTree)
	result = strings.Join([]string{ii, lt, pt}, "\n\n")
	return result
}
