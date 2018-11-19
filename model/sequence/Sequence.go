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

func NewSequence(id int, values []string) (sequence *Sequence) {
	sequence.ID = id
	sequence.Values = values
	return sequence
}

func FillSequence(sequence *Sequence, values []string) *Sequence {
	sequence.Values = values
	return sequence
}

func String(sequence *Sequence) string {
	return strings.Join([]string{"ID: ", fmt.Sprintf("%d", sequence.ID), "[", fmt.Sprintf("%v", sequence.Values), "]"}, " ")
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