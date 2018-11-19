package main

import (
	"strings"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"fmt"
)

type PredictionTree struct {

	Item string
	Parent *PredictionTree
	Children []*PredictionTree

}

func NewPredictionTreeRoot() *PredictionTree {

	return &PredictionTree{Item: "ROOT", Parent: nil, Children: []*PredictionTree{}}

}

func NewPredictionTree(value string) *PredictionTree {

	return &PredictionTree{Item: value, Parent: nil, Children: []*PredictionTree{}}

}

func AddChild(node *PredictionTree, value string) {

	newChild := NewPredictionTree(value)
	newChild.Parent = node
	node.Children = append(node.Children, newChild)

}

func GetChild(node *PredictionTree, value string) *PredictionTree {

	for _, child := range node.Children {
		if strings.EqualFold(child.Item, value) {
			return child
		}
	}
	return nil

}

func GetChildren(node *PredictionTree) []*PredictionTree {

	return node.Children

}

func HasChild(node *PredictionTree, value string) bool {

	if GetChild(node, value) != nil {
		return true
	}
	return false

}

func RemoveChild(node *PredictionTree, value string) {

	newChildren := make([]*PredictionTree, 0)
	for _, child := range node.Children {
		if !strings.EqualFold(child.Item, value) {
			newChildren = append(newChildren, child)
		}
	}
	node.Children = newChildren

}

type Sequence struct {

	ID string
	Values []string

}

type CPT struct {

	Alphabet []string
	Root *PredictionTree
	InvertedIndex map[string][]*Sequence
	LookupTable map[string]*PredictionTree

}

func NewCPT() *CPT {

	return &CPT{
		Alphabet: make([]string, 0),
		Root: NewPredictionTreeRoot(),
		InvertedIndex: make(map[string][]*Sequence, 0),
		LookupTable: make(map[string]*PredictionTree, 0),
	}

}

func LoadDataset(trainFile string) (train []*Sequence){

	f, e := os.Open(trainFile)
	if e != nil {
		log.Fatal("error: trainFile")
	}
	r := csv.NewReader(bufio.NewReader(f))
	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		train = append(train, &Sequence{ID: fmt.Sprintf("%d", count), Values: record})
		count += 1
		if count > 10 {
			break
		}
	}
	return train

}

func Train(cpt *CPT, trainSet []*Sequence) bool {

	cursorNode := cpt.Root
	for _, seq := range trainSet {
		for _, elem := range seq.Values {
			if HasChild(cursorNode, elem) == false {
				AddChild(cursorNode, elem)
			}
			cursorNode = GetChild(cursorNode, elem)
			if _, ok := cpt.InvertedIndex[elem]; !ok {
				cpt.InvertedIndex[elem] = make([]*Sequence, 0)
			}
			cpt.InvertedIndex[elem] = append(cpt.InvertedIndex[elem], seq)
			cpt.Alphabet = append(cpt.Alphabet, elem)
		}
		cpt.LookupTable[seq.ID] = cursorNode
		cursorNode = cpt.Root
	}
	return true

}

func Score(cpt *CPT, counttable map[string]float64, key string, length int, targetSize int, numberOfSimSeq int, numberOfCounttable int) map[string]float64 {

	weightLevel := 1/float64(numberOfSimSeq)
	weightDistance := 1/float64(numberOfCounttable)
	score := 1 + weightLevel + weightDistance * 0.001
	if _, ok := counttable[key]; !ok {
		counttable[key] = score
	} else {
		counttable[key] = score * counttable[key]
	}
	return counttable

}

func intersect(a []*Sequence, b []*Sequence) (inter []*Sequence) {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					if low[f1] != high[f2] {
						done = true
					}
				}
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		if done {
			break
		}
	}
	return
}

func Unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func FindLatest(a []string, x string) int {
	r := -1
	for i, n := range a {
		if x == n {
			r = i
		}
	}
	return r
}

func ConsequentScore(target *Sequence, similar *Sequence, index int) float64 {

	score := 0.0
	for _, c := range similar.Values[index:] {
		if strings.Contains(strings.Join(target.Values, ""), c) {
			score += 1.0
		}
	}
	return score

}

func Predict(cpt *CPT, train []*Sequence, test []*Sequence, k int, n int) map[string][]*Sequence {

	predictions := map[string][]*Sequence{}
	consequent := map[string]float64{}
	for _, targetSequence := range test {
		uniqueItems := Unique(targetSequence.Values)
		intersection := make([]*Sequence, 0)
		for _, element := range uniqueItems[len(uniqueItems)-k:] {
			intersection = intersect(intersection, cpt.InvertedIndex[element])
		}
		for _, similarSequence := range intersection {
			lastTargetSequenceItem := targetSequence.Values[len(targetSequence.Values)-1]
			i := FindLatest(similarSequence.Values, lastTargetSequenceItem)
			if i != -1 {
				score := ConsequentScore(targetSequence, similarSequence, i)
				if _, ok := consequent[targetSequence.Values[i]]; !ok {
					consequent[targetSequence.Values[i]] = 1 + (1 / float64(len(intersection))) + (1 / float64(len(consequent) + 1)) * 0.001
				} else {
					consequent[targetSequence.Values[i]] = (1 + (1 / float64(len(intersection))) + (1 / float64(len(consequent) + 1)) * 0.001) * score
				}
				fmt.Printf("%v", consequent)
			}
		}
		topp := getNLargest(cpt, consequent, n)
		predictions[targetSequence.ID] = topp
	}
	return predictions

}

func getNLargest(cpt *CPT, dictionary map[string]float64, n int) []*Sequence {

	return []*Sequence{}
	// fmt.Printf(dictionary)

}

func PrintPredictionTree(t *PredictionTree) {

}

func PrintCPT(cpt *CPT) {

	fmt.Println()

}

func main() {
	train := LoadDataset("./data/train.csv")
	test  := LoadDataset("./data/test.csv")
	cpt := &CPT{}
	cpt.Root = NewPredictionTreeRoot()
	cpt.InvertedIndex = make(map[string][]*Sequence, 0)
	cpt.LookupTable = make(map[string]*PredictionTree, 0)
	Train(cpt, train)
	_ = Predict(cpt, train, test,5,3)
	// fmt.Printf(predictions)
}