package main

import (
	"fmt"

	"github.com/made2591/go-cpt/model/compactPredictionTree"
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"
)

func main() {
	trainingSequences := sequence.ReadCSVSequencesFile("./data/dummy.csv")
	testingSequences := sequence.ReadCSVSequencesFile("./data/dumbo.csv")
	trainingSequences = sequence.ReadCSVSequencesFile("./data/train.csv")[1:11]
	testingSequences = sequence.ReadCSVSequencesFile("./data/test.csv")[1:11]
	//testingSequences = sequence.ReadCSVSequencesFile("./data/test.csv")[3:4]
	//for _, seq := range trainingSequences {
	//	fmt.Println(sequence.String(seq))
	//}
	//for _, seq := range testingSequences {
	//	fmt.Println(sequence.String(seq))
	//}
	invertedIndex := invertedIndexTable.NewInvertedIndexTable(trainingSequences)
	lookup := lookupTable.NewLookupTable(trainingSequences)
	predTree := predictionTree.NewPredictionTree("ROOT")
	cpt := compactPredictionTree.NewCompactPredictionTree(invertedIndex, lookup, predTree, trainingSequences, testingSequences)
	compactPredictionTree.InitCompactPredictionTree(cpt, trainingSequences)
	fmt.Println(predictionTree.String(cpt.PredictionTree))

	fmt.Println(compactPredictionTree.PredictionOverTestingSequence(cpt,5, 3))

	//for _, s := range train {
	//	PrintSequence(s)
	//}
	//for _, s := range test {
	//	PrintSequence(s)
	//}
	//cpt := &CPT{}
	//cpt.Root = New("ROOT")
	//cpt.InvertedIndex = make(map[string][]*Sequence, 0)
	//cpt.LookupTable = make(map[string]*Node, 0)
	//Train(cpt, train)
	//PrintCPT(cpt)
	//_ = Predict(cpt, train, test,5,3)
	//// fmt.Printf(predictions)
}
