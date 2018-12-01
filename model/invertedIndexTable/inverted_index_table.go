package invertedIndexTable

import (
	"strings"

	"github.com/made2591/go-cpt/model/sequence"
)

// The struct InvertedIndexTable
type InvertedIndexTable struct {
	keys   []string
	values []*sequence.Sequence
	Table  map[string][]*sequence.Sequence
}

// NewInvertedIndexTable creates a new InvertedIndexTable
func NewInvertedIndexTable(sequences map[int]*sequence.Sequence) (invertedIndexTable *InvertedIndexTable) {

	// setup inverted index table
	invertedIndexTable = &InvertedIndexTable{}

	// setup keys
	invertedIndexTable.keys = make([]string, 0)
	invertedIndexTable.Table = make(map[string][]*sequence.Sequence, 0)

	// add sequences
	for seqID := 0; seqID < len(sequences); seqID++ {
		for _, value := range sequences[seqID].Values {
			AppendIfMissing(invertedIndexTable.keys, value)
		}
	}
	for _, symbol := range invertedIndexTable.keys {
		invertedIndexTable.Table[symbol] = nil
	}
	return
}

// AppendIfMissing appends element to slice
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// AddSequenceIfMissing appends element to slice
func AddSequenceIfMissing(invertedIndexTable *InvertedIndexTable, key string, seq *sequence.Sequence) (bool, *sequence.Sequence) {

	// add sequence to table
	for _, cseq := range invertedIndexTable.Table[key] {
		if sequence.EqualSequence(seq, cseq) {
			return false, nil
		}
	}
	invertedIndexTable.Table[key] = append(invertedIndexTable.Table[key], seq)
	return true, seq

}

// GetSequencesForSymbol provides sequences for a given symbol
func GetSequencesForSymbol(invertedIndexTable *InvertedIndexTable, key string) (bool, []*sequence.Sequence) {
	if value, found := invertedIndexTable.Table[key]; found == true {
		return found, value
	}
	return false, nil
}

// Stringify the NewInvertedIndexTable
func String(invertedIndexTable *InvertedIndexTable) (result string) {
	result = ""
	for key, sequences := range invertedIndexTable.Table {
		result = strings.Join([]string{result, key, " -> "}, "")
		for _, seq := range sequences {
			result = strings.Join([]string{result, sequence.String(seq)}, " ")
		}
		result = strings.Join([]string{result, "\n"}, "")
	}
	return result
}
