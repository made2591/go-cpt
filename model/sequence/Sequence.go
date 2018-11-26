package sequence

import (
	"strings"
	"fmt"
	"os"
	"log"
	"encoding/csv"
	"bufio"
	"io"
)

type Sequence struct{
	ID int
	Values []string
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
	result := []string{}
	for _, c := range sequence.Values {
		if !stringInSlice(c, sequence.Values) {
			result = append(result, c)
		}
	}
	return result
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ComputeConsequent(seq1 *Sequence, seq2 *Sequence) *Sequence {
	result := &Sequence{Values: make([]string, 0)}
	for i := len(seq2.Values)-1; i >= 0; i-- {
		if strings.Compare(seq2.Values[i], seq1.Values[len(seq1.Values)-1]) == 0 {
			return result
		} else {
			if !stringInSlice(seq2.Values[i], result.Values) {
				result.Values = append(result.Values, seq2.Values[i])
			}
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
	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		result = append(result, &Sequence{ID: count, Values: record})
		count += 1
		// TODO REMOVE LIMIT
		if count > 10 {
			break
		}
	}
	return result

}