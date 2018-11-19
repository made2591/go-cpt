package invertedIndexTable

import (
	"github.com/made2591/go-cpt/model/sequence"
	"strings"
)

type InvertedIndexTable struct {

	keys []string
	values []*sequence.Sequence
	Table map[string][]*sequence.Sequence

}

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func NewInvertedIndexTable(sequences []*sequence.Sequence) (invertedIndexTable *InvertedIndexTable) {
	invertedIndexTable = &InvertedIndexTable{}
	invertedIndexTable.keys = make([]string, 0)
	invertedIndexTable.Table = make(map[string][]*sequence.Sequence, 0)
	for _, seq := range sequences {
		for _, value := range seq.Values {
			AppendIfMissing(invertedIndexTable.keys, value)
		}
	}
	for _, symbol := range invertedIndexTable.keys {
		invertedIndexTable.Table[symbol] = nil
	}
	return
}

func GetElementWithValue(invertedIndexTable InvertedIndexTable, key string) (bool, []*sequence.Sequence) {
	if value, found := invertedIndexTable.Table[key]; found == true {
		return found, value
	}
	return false, nil
}

func String(invertedIndexTable *InvertedIndexTable) (result string) {
	for _, key := range invertedIndexTable.keys {
		result = key
		sequences := "["
		for _, seq := range invertedIndexTable.Table[key] {
			sequences = strings.Join([]string{sequences, "[", sequence.String(seq), "]"}, " ")
		}
		sequences = strings.Join([]string{sequences, "]\n"}, "")
		result = strings.Join([]string{result, sequences}, " -> ")
	}
	return result
}