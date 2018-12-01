// Package compactPredictionTree provides a wrapper around the structure needed to build
// compact prediction tree
package compactPredictionTree

import (

	"os"
	"sort"
	"strings"

	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"

	"github.com/op/go-logging"
	set "github.com/deckarep/golang-set"

)

var log = logging.MustGetLogger("go-cpt")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type Secret string

func (p Secret) Redacted() interface{} {
	return logging.Redact(string(p))
}

func init() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	//log.Debugf("debug %s", Secret("secret"))
	//log.Info("info")
	//log.Notice("notice")
	//log.Warning("warning")
	//log.Error("err")
	//log.Critical("crit")
}

// The struct CompactPredictionTree
type CompactPredictionTree struct {
	InvertedIndexTable *invertedIndexTable.InvertedIndexTable
	PredictionTree     *predictionTree.PredictionTree
	TrainingSet        map[int]*sequence.Sequence
	TestingSet         map[int]*sequence.Sequence
}

// NewCompactPredictionTree create a new CompactPredictionTree
func NewCompactPredictionTree(
	invertedIndexTable *invertedIndexTable.InvertedIndexTable,
	predictionTree *predictionTree.PredictionTree,
	trainingSet map[int]*sequence.Sequence,
	testingSet map[int]*sequence.Sequence) (compactPredictionTree *CompactPredictionTree) {

	// setup compact prediction tree
	compactPredictionTree = &CompactPredictionTree{}

	// setup inverted index table
	compactPredictionTree.InvertedIndexTable = invertedIndexTable

	// setup prediction tree
	compactPredictionTree.PredictionTree = predictionTree

	// setup training and testing set
	compactPredictionTree.TrainingSet = trainingSet
	compactPredictionTree.TestingSet = testingSet

	// return compact prediction tree
	return compactPredictionTree

}

// InitCompactPredictionTree init structure
func InitCompactPredictionTree(compactPredictionTree *CompactPredictionTree) {

	// save the cursore node to create the compact prediction tree
	cursorNode := compactPredictionTree.PredictionTree

	// inverted index table creation
	invertedIndex := compactPredictionTree.InvertedIndexTable

	// for each sequence in the training set (looping to maintain order)
	for seqIndex := 0; seqIndex < len(compactPredictionTree.TrainingSet); seqIndex++ {

		// for each symbol in the sequence
		for _, symbol := range compactPredictionTree.TrainingSet[seqIndex].Values {

			// if the symbol is not found as child of the current prediction tree pointed by the cursor
			if found, child := predictionTree.GetChildWithValue(cursorNode, symbol); found == false {

				// then add the symbol and use is as a new cursor
				cursorNode = predictionTree.AddChild(cursorNode, symbol)

			} else {

				// otherwise move to cursor to the found node - to create the compression
				cursorNode = child

			}

			// add symbol to toe inverted index table
			invertedIndexTable.AddSequenceIfMissing(invertedIndex, symbol, compactPredictionTree.TrainingSet[seqIndex])

		}

		// reset cursor node to the root node of the prediction tree
		cursorNode = compactPredictionTree.PredictionTree

	}

}

// PredictionOverTestingSequence init structure
func PredictionOverTestingSequence(compactPredictionTree *CompactPredictionTree, k int, n int) map[int][]string {

	// create results struct
	results := make(map[int][]string, len(compactPredictionTree.TestingSet))

	// for each sequence in the training set (looping to maintain order)
	for seqIndex := 0; seqIndex < len(compactPredictionTree.TestingSet); seqIndex++ {

		// prepare result of predictions for the targeted sequence
		result := make([]string, 0)

		// isolate tested sequence
		testedSequence := compactPredictionTree.TestingSet[seqIndex]

		// create a new intersection set with the IDs of all the testing sequences
		intersection := set.NewSet()
		for i := 0; i < len(compactPredictionTree.TestingSet); i++ {
			intersection.Add(i)
		}

		log.Debugf("original target: %v", testedSequence.Values)

		// reduce the context
		if k < len(compactPredictionTree.TestingSet[seqIndex].Values) {

			// cut the sequence targeted in test base on the parameter in input
			testedSequence = &sequence.Sequence{Values: testedSequence.Values[len(testedSequence.Values)-k:]}
			log.Debugf("cut target: %v", testedSequence.Values)

		}

		log.Debugf("tested sequence: %v", testedSequence.Values)

		// for every symbol in the tested sequence
		for _, symbol := range testedSequence.Values {

			// if the symbol appears in some sequence
			if len(compactPredictionTree.InvertedIndexTable.Table[symbol]) != 0 {

				log.Debugf("\tsymbol: %s, seqIDs: ", symbol)

				// create a new set of IDs of sequence in which the symbol appears
				seqID := set.NewSet()
				for _, seq := range compactPredictionTree.InvertedIndexTable.Table[symbol] {
					seqID.Add(seq.ID)
					log.Debugf("%d ", seq.ID)
				}
				log.Debugf("\t\tintersection set before intersection: %s", intersection.String())

				// intersect the set
				intersection = intersection.Intersect(seqID)
				log.Debugf("\t\tintersection set after intersection: %s", intersection.String())

			}

		}

		// create iterator over intersection set of IDs in which every symbol of the targeted sequence appear
		it := intersection.Iterator()

		// init similar sequences set
		similarSequences := make([]*sequence.Sequence, 0)

		// init the slice with the right dimension
		similarSequences = make([]*sequence.Sequence, len(it.C))

		log.Infof("number of similar sequences: %d", len(similarSequences))
		for _, seq := range similarSequences {
			log.Debugf("\t%v", sequence.String(seq))
		}

		// fullfill the set of sequences
		for seqID := range it.C {
			log.Debugf("\telement: %d", seqID.(int))
			similarSequences = append(similarSequences, compactPredictionTree.TrainingSet[seqID.(int)])
		}

		// init map of consequent
		consequents := make(map[int][]string, 0)

		// for each similar sequence
		for _, similarSequence := range similarSequences {

			log.Debugf("\tsequence: %d", similarSequence.ID)

			// compute the consequent, i.e. the longest subsequence of symbols
			// after the latest occurence of the latest symbol of tested sequence
			consequent := sequence.ComputeConsequent(testedSequence, similarSequence)

			// if there are some symbols
			if len(consequent) > 0 {

				// saved it for later
				consequents[similarSequence.ID] = consequent
				log.Debugf("\tconsequent: %v", consequent)

			}

		}

		// init score dictionary
		scores := make(map[string]float64, 0)

		// for every consequents
		for _, consequent := range consequents {

			// for each symbol in the consequent
			for _, symbol := range consequent {

				// compute the score
				scores = computeScore(scores, symbol, len(testedSequence.Values), len(testedSequence.Values), len(similarSequences), len(scores)+1)
				log.Debugf("\tsymbol: %s, %f", symbol, scores[symbol])

			}

		}

		// order by score
		pairs := rankByScore(scores)

		// return predictions if any
		if len(pairs) > 0 {
			for i := 0; i < n-1 && i < len(pairs); i++ {
				result = append(result, pairs[i].Key)
				log.Debugf("\tsymbol: %s, %f", pairs[i].Key, pairs[i].Value)
			}
		}

		// add sequence to results
		results[seqIndex] = result

	}

	log.Debugf("results: %v", results)
	return results

}

// Stringify the compactPredictionTree
func String(compactPredictionTree CompactPredictionTree) (result string) {

	ii := invertedIndexTable.String(compactPredictionTree.InvertedIndexTable)
	pt := predictionTree.String(compactPredictionTree.PredictionTree)
	result = strings.Join([]string{ii, pt}, "\n\n")

	return result

}


// computeScore compute the score
func computeScore(scores map[string]float64,
	key string,
	length int,
	target_size int,
	number_of_similar_sequences int,
	number_items_counttable int) map[string]float64 {


	weight_level := float64(1 / number_of_similar_sequences)
	weight_distance := float64(1 / number_items_counttable)
	score := 1.0 + weight_level + weight_distance*0.001

	if oldScore, ok := scores[key]; !ok {
		scores[key] = score
	} else {
		scores[key] = score * oldScore
	}

	return scores
}


type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// rankByScore order score
func rankByScore(scores map[string]float64) PairList {

	pl := make(PairList, len(scores))
	i := 0
	for k, v := range scores {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl

}