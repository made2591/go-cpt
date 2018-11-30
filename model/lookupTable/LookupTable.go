package lookupTable

import (
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/sequence"
	"strings"
	"fmt"
)

type LookupTable struct {

	keys []int
	values []*predictionTree.PredictionTree
	Table map[int]*sequence.Sequence

}

func NewLookupTable(sequences []*sequence.Sequence) (lookupTable *LookupTable) {
	lookupTable = &LookupTable{}
	lookupTable.keys = make([]int, 0)
	lookupTable.Table = make(map[int]*sequence.Sequence)
	for _, seq := range sequences {
		lookupTable.keys = append(lookupTable.keys, seq.ID)
		lookupTable.Table[seq.ID] = nil
	}
	return lookupTable
}

func String(lookupTable *LookupTable) (result string) {
	for _, key := range lookupTable.keys {
		result = strings.Join([]string{fmt.Sprintf("%d", key), " -> ", sequence.String(lookupTable.Table[key]), "\n"}, "")
	}
	return result
}