package main

import (
	"fmt"
	"github.com/made2591/go-cpt/model/sequence"
	"github.com/made2591/go-cpt/model/compactPredictionTree"
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
)

func main() {
    trainingSequences := sequence.ReadCSVSequencesFile("./data/dummy.csv")
	testingSequences := sequence.ReadCSVSequencesFile("./data/dumbo.csv")
	trainingSequences = sequence.ReadCSVSequencesFile("./data/train.csv")[1:11]
	testingSequences = sequence.ReadCSVSequencesFile("./data/test.csv")[1:11]
	//for _, seq := range trainingSequences {
	//	fmt.Println(sequence.String(seq))
	//}
	//for _, seq := range testingSequences {
	//	fmt.Println(sequence.String(seq))
	//}
	invertedIndex := invertedIndexTable.NewInvertedIndexTable(trainingSequences)
	lookup := lookupTable.NewLookupTable(trainingSequences)
	predTree := predictionTree.NewPredictionTree("ROOT")
 	cpt := compactPredictionTree.NewCompactPredictionTree(invertedIndex, lookup, predTree)
	compactPredictionTree.InitCompactPredictionTree(cpt, trainingSequences)
	// fmt.Println(compactPredictionTree.String(*cpt))

	for _, seq := range testingSequences {
		fmt.Println(sequence.String(seq))
		fmt.Println(compactPredictionTree.PredictionOverTestingSequence(cpt, seq, 5, 1))
	}

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