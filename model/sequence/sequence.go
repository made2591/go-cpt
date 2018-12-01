package sequence

import (
	"io"
	"os"
	"fmt"
	"log"
	"bufio"
	"strings"
	"encoding/csv"
)

// The struct Sequence
type Sequence struct {

	ID     int
	Values []string

}

// NewSequence creates a new Sequence
func NewSequence(id int, values []string) (sequence *Sequence) {

	return &Sequence{ID: id, Values: values}

}

// SameSequence creates a new Sequence
func SameSequence(seq1 *Sequence, seq2 *Sequence) bool {

	// if one or both of the sequence are nil return false
	if seq1 == nil || seq2 == nil {
		return false
	}

	// check the ID
	if seq1.ID == seq2.ID {
		return true
	}

	return false

}

// LatestStringInSlice creates a new Sequence
func LatestStringInSlice(symbol string, list []string) (bool, int) {

	// if starting from the end the symbol is found
	for i := len(list)-1; i >= 0; i-- {

		if strings.EqualFold(symbol, list[i]) {
			return true, i
		}

	}

	return false, -1

}

// ComputeConsequent creates a new Sequence
func ComputeConsequent(seq1 *Sequence, seq2 *Sequence) []string {

	// create result
	result := make([]string, 0)

	// if seq2 is similar to seq1 (contains is last character)
	if found, index := LatestStringInSlice(seq1.Values[len(seq1.Values)-1], seq2.Values); found {

		// for every values from there to end
		for i := index; i < len(seq2.Values); i++ {

			// keep value not found in original seq1
			if found, _ := LatestStringInSlice(seq2.Values[i], seq1.Values); !found {
				result = append(result, seq2.Values[i])
			}

		}

	}

	return result

}

// String provides a string of Sequence
func String(sequence *Sequence) string {

	return strings.Join([]string{"ID: ",
			fmt.Sprintf("%d", sequence.ID), " Values [",
			fmt.Sprintf("%v", sequence.Values), "]"}, "")

}

// ReadCSVSequencesFile provides a map of Sequences given a lower and upper limit
func ReadCSVSequencesFile(filepath string, limits ...int) (map[int]*Sequence) {

	result := map[int]*Sequence{}

	f, e := os.Open(filepath)
	if e != nil {
		log.Fatal("error: trainFile")
	}
	r := csv.NewReader(bufio.NewReader(f))
	row := 0
	id := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if len(limits) > 1 {
			if row >= limits[1] {
				break
			}
		}
		if len(limits) > 0 {
			if row >= limits[0] {
				result[id] = NewSequence(id, record)
				id += 1
			}
		} else {
			result[id] = NewSequence(id, record)
			id += 1
		}
		row += 1
	}

	return result

}
