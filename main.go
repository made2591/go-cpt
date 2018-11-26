package main

import (
	"strings"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"fmt"
	"sort"
	"github.com/made2591/go-cpt/model/sequence"
	"github.com/made2591/go-cpt/model/compactPredictionTree"
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/lookupTable"
	"github.com/made2591/go-cpt/model/predictionTree"
)

const (
	TraverseInOrder TraverseType = iota
	TraversePreOrder
	TraversePostOrder
	TraverseLevelOrder
)

const (
	TraverseLeaves TraverseFlags = 1 << iota
	TraverseNonLeaves
	TraverseMask = 0x3
	TraverseAll  = TraverseLeaves | TraverseNonLeaves
)

type TraverseFunc func(*Node, string) bool
type TraverseType int
type TraverseFlags int

type Node struct {
	Value    string
	Next     *Node
	Previous *Node
	Parent   *Node
	Children *Node
}

type nodeVal struct {
	Value         string
	NodeReference *Node
}

func New(v string) *Node {
	return &Node{Value: v}
}

func Unlink(n *Node) {
	if n == nil {
		return
	}

	if n.Previous != nil {
		n.Previous.Next = n.Next
	} else if n.Parent != nil {
		n.Parent.Children = n.Next
	}

	n.Parent = nil
	if n.Next != nil {
		n.Next.Previous = n.Previous
		n.Next = nil
	}
	n.Previous = nil
}

func Depth(n *Node) int {
	depth := 0

	for n != nil {
		depth++
		n = n.Parent
	}

	return depth
}

func Insert(parent, n *Node) *Node {
	if parent == nil || n == nil || !IsRoot(n) {
		return nil
	}

	return AppendChild(parent, n)
}

func IsRoot(n *Node) bool {
	return n.Parent == nil && n.Previous == nil && n.Next == nil
}

func NodeCount(root *Node, flags TraverseFlags) int {
	if root == nil || flags > TraverseMask {
		return 0
	}

	n := 0
	nodeCountFunc(root, flags, &n)
	return n
}

func nodeCountFunc(n *Node, flags TraverseFlags, count *int) {
	if n.Children != nil {
		if flags&TraverseNonLeaves != 0 {
			(*count)++
		}

		child := n.Children
		for child != nil {
			nodeCountFunc(child, flags, count)
			child = child.Next
		}
	} else if flags&TraverseLeaves != 0 {
		(*count)++
	}
}

func AppendChild(parent, n *Node) *Node {
	if parent == nil || n == nil || !IsRoot(n) {
		return nil
	}

	n.Parent = parent
	if parent.Children != nil {
		sibling := parent.Children
		for sibling.Next != nil {
			sibling = sibling.Next
		}
		n.Previous = sibling
		sibling.Next = n
	} else {
		n.Parent.Children = n
	}

	return n
}

func GetRoot(n *Node) (*Node, int) {
	if n == nil {
		return nil, 0
	}

	depth := 1

	current := n
	for current.Parent != nil {
		depth++
		current = current.Parent
	}

	return current, depth
}

func HasChild(root *Node, data string) bool {
	if root == nil {
		return false
	}
	for root.Next != nil {
		if strings.EqualFold(root.Value, data) {

		}
	}
	return false
}

func FindNode(root *Node, order TraverseType, flags TraverseFlags, data string) *Node {
	if root == nil {
		return nil
	}

	d := &nodeVal{Value: data}
	Traverse(root, order, flags, -1, nodeFindFunc, d.Value)

	return d.NodeReference
}

func nodeFindFunc(n *Node, data string) bool {
	if n.Value != data {
		return false
	}
	return true
}

// Traverse traverses a node as the root node based on the passed in TraverseType, TraverseFlags, and Depth.  Each visited node will
// have TraverseFunc called, passing along Data to each node.
func Traverse(root *Node, order TraverseType, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data string) {
	if root == nil || traverseFunc == nil || order > TraverseLevelOrder || flags > TraverseMask || (depth < -1 || depth == 0) {
		return
	}

	switch order {
	default:
		fallthrough
	case TraversePreOrder:
		if depth < 0 {
			traversePreOrder(root, flags, traverseFunc, data)
		} else {
			depthTraversePreOrder(root, flags, depth, traverseFunc, data)
		}
	case TraverseInOrder:
		if depth < 0 {
			traverseInOrder(root, flags, traverseFunc, data)
		} else {
			depthTraverseInOrder(root, flags, depth, traverseFunc, data)
		}
	case TraversePostOrder:
		if depth < 0 {
			traversePostOrder(root, flags, traverseFunc, data)
		} else {
			depthTraversePostOrder(root, flags, depth, traverseFunc, data)
		}

		// case Traverse_LevelOrder:
		// 	panic("Not Implemented")
		// 	// 	g_node_depth_traverse_level (root, flags, depth, func, data);
	}
}

func traversePreOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

		child := n.Children
		for child != nil {
			current := child
			child = current.Next
			if traversePreOrder(current, flags, traverseFunc, data) {
				return true
			}
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func depthTraversePreOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

		depth--
		if depth == 0 {
			return false
		}

		child := n.Children
		for child != nil {
			current := child
			child = current.Next
			if depthTraversePreOrder(current, flags, depth, traverseFunc, data) {
				return true
			}
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func traverseInOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		child := n.Children
		current := child
		child = current.Next
		if traverseInOrder(current, flags, traverseFunc, data) {
			return true
		}

		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

		for child != nil {
			current = child
			child = current.Next
			if traverseInOrder(current, flags, traverseFunc, data) {
				return true
			}
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func depthTraverseInOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		depth--
		if depth > 0 {
			child := n.Children
			current := child
			child = current.Next

			if depthTraverseInOrder(current, flags, depth, traverseFunc, data) {
				return true
			}

			if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
				return true
			}

			for child != nil {
				current = child
				child = current.Next
				if depthTraverseInOrder(current, flags, depth, traverseFunc, data) {
					return true
				}
			}
		} else if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}
	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func traversePostOrder(n *Node, flags TraverseFlags, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		child := n.Children
		for child != nil {

			current := child
			child = current.Next
			if traversePostOrder(current, flags, traverseFunc, data) {
				return true
			}
		}

		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true

	}

	return false
}

func depthTraversePostOrder(n *Node, flags TraverseFlags, depth int, traverseFunc TraverseFunc, data string) bool {
	if n.Children != nil {
		depth--
		if depth > 0 {

			child := n.Children
			for child != nil {

				current := child
				child = current.Next
				if depthTraversePostOrder(current, flags, depth, traverseFunc, data) {
					return true
				}
			}
		}

		if (flags&TraverseNonLeaves != 0) && traverseFunc(n, data) {
			return true
		}

	} else if (flags&TraverseLeaves != 0) && traverseFunc(n, data) {
		return true
	}

	return false
}

func (n *Node) String() string {
	if n == nil {
		return "()"
	}

	currentLevel := 0
	lastNode := n
	levels := make([]string, 0, 10)
	levels = append(levels, "")
	tFunc := func(node *Node, value string) bool {
		currentLevel = 0
		n := node.Parent
		for n != nil {
			currentLevel++
			if len(levels) <= currentLevel {
				levels = append(levels, "")
			}
			n = n.Parent
		}
		levels[currentLevel] += fmt.Sprintf("(%v)", node.Value) + "\t"
		lastNode = node
		return false
	}

	Traverse(n, TraversePreOrder, TraverseAll, -1, tFunc, "")
	s := ""
	for _, v := range levels {
		s += v + "\n\n"
	}
	return s
}

type Sequence struct {

	ID string
	Values []string

}

type Pair struct {
	Key string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

type CPT struct {

	Alphabet []string
	Root *Node
	InvertedIndex map[string][]*Sequence
	LookupTable map[string]*Node

}

func NewCPT() *CPT {

	return &CPT{
		Alphabet: make([]string, 0),
		Root: New("ROOT"),
		InvertedIndex: make(map[string][]*Sequence, 0),
		LookupTable: make(map[string]*Node, 0),
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
				newNode := &Node{Value: elem}
				AppendChild(cursorNode, newNode)
			}
			cursorNode = cursorNode.Children
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

func getNLargest(wordFrequencies map[string]float64, n int) (res []string) {

	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	for c := 0; c < n && c < len(pl); c++ {
		res = append(res, pl[c].Key)
	}
	return res

}

func Predict(cpt *CPT, train []*Sequence, test []*Sequence, k int, n int) map[string][]string {

	predictions := map[string][]string{}
	consequent := map[string]float64{}
	for _, targetSequence := range test {
		uniqueTargetSequence := &Sequence{ID: targetSequence.ID, Values: Unique(targetSequence.Values)}
		intersection := make([]*Sequence, 0)
		for c, element := range uniqueTargetSequence.Values[len(uniqueTargetSequence.Values)-k:] {
			if _, ok := cpt.InvertedIndex[element]; ok {
				if c == 0 {
					intersection = cpt.InvertedIndex[element]
				}
				intersection = intersect(intersection, cpt.InvertedIndex[element])
			}
		}
		for _, inter := range intersection {
			PrintSequence(inter)
		}
		for _, similarSequence := range intersection {
			lastTargetSequenceItem := targetSequence.Values[len(targetSequence.Values)-1]
			i := FindLatest(similarSequence.Values, lastTargetSequenceItem)
			fmt.Println(i)
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
		topp := getNLargest(consequent, n)
		predictions[targetSequence.ID] = topp
	}
	return predictions

}

func PrintSequence(s *Sequence) {
	fmt.Print(s.ID, " ")
	fmt.Println(s.Values)
}

func PrintNode(t *Node, level ...int) {

	l := 0
	if len(level) > 0 {
		l = level[0]
	}
	for i := 0; i < l; l++ {
		fmt.Print("\t")
	}
	fmt.Print("%s\n", t.Value)
	l++
	c := t.Children
	for c != nil {
		PrintNode(c, l)
		c = c.Children
	}
	return

}

func PrintInvertedIndex(invertedIndex map[string][]*Sequence) {

	for k, v := range invertedIndex {
		fmt.Print(k, " | ")
		for _, s := range v {
			PrintSequence(s)
		}
	}

}

func PrintCPT(cpt *CPT) {

	PrintNode(cpt.Root)
	PrintInvertedIndex(cpt.InvertedIndex)

}

func main() {
	trainingSequences := sequence.ReadCSVSequencesFile("./data/train.csv")
	testingSequences := sequence.ReadCSVSequencesFile("./data/test.csv")
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
	for _, seq := range testingSequences {
		//fmt.Println(sequence.String(seq))
		fmt.Println(compactPredictionTree.PredictionOverTestingSequence(cpt, seq))
	}

	//fmt.Println(compactPredictionTree.String(*cpt))

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