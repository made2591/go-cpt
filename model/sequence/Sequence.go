package sequence

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Sequence struct {
	ID     int
	Values []string
}

func AddIfNotExists(seqs []*Sequence, seq *Sequence) (bool, int) {
	for _, s := range seqs {
		if EqualSequence(s, seq) {
			return false, -1
		}
	}
	seqs = append(seqs, seq)
	return true, len(seqs)
}

func EqualSequence(seq1 *Sequence, seq2 *Sequence) bool {
	if seq1 == nil && seq2 != nil {
		return false
	}
	if seq1 != nil && seq2 == nil {
		return false
	}
	if seq1.ID == seq2.ID {
		return true
	}
	return false
}

func NewSequence(id int, values []string) (sequence *Sequence) {
	sequence.ID = id
	sequence.Values = values
	return sequence
}

func FillSequence(sequence *Sequence, values []string) *Sequence {
	sequence.Values = values
	return sequence
}

func UniqueElements(sequence *Sequence) []string {
	result := make([]string, 0)
	for _, c := range sequence.Values {
		if found, _ := StringInSlice(c, result); !found {
			result = append(result, c)
		}
	}
	return result
}

func StringInSlice(a string, list []string) (bool, int) {
	for i := len(list)-1; i >= 0; i-- {
		if strings.EqualFold(a, list[i]) {
			return true, i
		}
	}
	return false, -1
}

func ComputeConsequent(seq1 *Sequence, seq2 *Sequence) []string {
	result := make([]string, 0)
	if found, index := StringInSlice(seq1.Values[len(seq1.Values)-1], seq2.Values); found {
		for i := index; i < len(seq2.Values); i++ {
			if found, _ := StringInSlice(seq2.Values[i], seq1.Values); !found {
				result = append(result, seq2.Values[i])
			}
		}
	}
	return result
}

func LastNSymbols(sequence *Sequence, n int) []string {
	if sequence == nil {
		return make([]string, 0)
	}
	result := make([]string, 0)
	for i, c := range sequence.Values {
		if i >= (len(sequence.Values) - n) {
			result = append(result, c)
		}
	}
	return result
}

func String(sequence *Sequence) string {
	return strings.Join([]string{"ID: ", fmt.Sprintf("%d", sequence.ID), " Values [", fmt.Sprintf("%v", sequence.Values), "]"}, "")
}

func ReadCSVSequencesFile(filepath string) (result []*Sequence) {

	f, e := os.Open(filepath)
	if e != nil {
		log.Fatal("error: trainFile")
	}
	r := csv.NewReader(bufio.NewReader(f))
	count := -1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		result = append(result, &Sequence{ID: count, Values: record})
		count += 1
	}
	return result

}
