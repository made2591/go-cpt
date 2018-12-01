package main

import (

	"os"
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"path/filepath"

	"github.com/made2591/go-cpt/model/sequence"
	"github.com/made2591/go-cpt/model/predictionTree"
	"github.com/made2591/go-cpt/model/invertedIndexTable"
	"github.com/made2591/go-cpt/model/compactPredictionTree"

	"strings"
)

const maxUploadSize = 20 * 1024 * 1024 // 2 mb
const uploadPath = "./uploads"

func local() {
	trainingSequences := sequence.ReadCSVSequencesFile("./data/dummy.csv")
	testingSequences := sequence.ReadCSVSequencesFile("./data/dumbo.csv")
	trainingSequences = sequence.ReadCSVSequencesFile("./data/train.csv", 1, 11)
	testingSequences = sequence.ReadCSVSequencesFile("./data/test.csv", 1, 11)
	for _, seq := range trainingSequences {
		fmt.Println(sequence.String(seq))
	}
	for _, seq := range testingSequences {
		fmt.Println(sequence.String(seq))
	}
	invertedIndex := invertedIndexTable.NewInvertedIndexTable(trainingSequences)
	predTree := predictionTree.NewPredictionTree("ROOT")
	cpt := compactPredictionTree.NewCompactPredictionTree(invertedIndex, predTree, trainingSequences, testingSequences)
	compactPredictionTree.InitCompactPredictionTree(cpt)
	fmt.Println(predictionTree.String(cpt.PredictionTree))

	predictions := compactPredictionTree.PredictionOverTestingSequence(cpt,5, 3)
	for i := 0; i < len(testingSequences); i++ {
		fmt.Println(testingSequences[i].Values)
		fmt.Println(predictions[i])
	}
}

func initcpt() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trainingSequences := sequence.ReadCSVSequencesFile(strings.Join([]string{"./", uploadPath, "/train.csv"}, ""), 1, 11)
		testingSequences := sequence.ReadCSVSequencesFile(strings.Join([]string{"./", uploadPath, "/test.csv"}, ""), 1, 11)
		invertedIndex := invertedIndexTable.NewInvertedIndexTable(trainingSequences)
		predTree := predictionTree.NewPredictionTree("ROOT")
		cpt := compactPredictionTree.NewCompactPredictionTree(invertedIndex, predTree, trainingSequences, testingSequences)
		compactPredictionTree.InitCompactPredictionTree(cpt)
		fmt.Println(predictionTree.String(cpt.PredictionTree))

		predictions := compactPredictionTree.PredictionOverTestingSequence(cpt, 5, 3)
		for i := 0; i < len(testingSequences); i++ {
			fmt.Println(testingSequences[i].Values)
			fmt.Println(predictions[i])
		}
	})
}

func uploadTrain() http.HandlerFunc {
	return uploadFileHandler("train.csv")
}

func uploadTest() http.HandlerFunc {
	return uploadFileHandler("test.csv")
}

func uploadFileHandler(position string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// validate file size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		fileType := r.PostFormValue("type")
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		switch filetype {
			case "text/plain; charset=utf-8":
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
				return
		}
		newPath := filepath.Join(uploadPath, position)
		fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func main() {

	http.HandleFunc("/upload/train", uploadTrain())
	http.HandleFunc("/upload/test", uploadTest())
	http.HandleFunc("/initcpt", initcpt())

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Print("Server started on localhost:8080, use /upload/[train/test] for uploading train/test files and /files/{fileName} for downloading. Use /initcpt to start training and obtain predictions")
	log.Fatal(http.ListenAndServe(":8080", nil))

}