package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCSVFile(filename string) ([][]string, error) {
	log.Println("Open file:", filename)

	csvFile, err := os.Open(filename)
	defer csvFile.Close()

	if err != nil {
		return nil, err
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	return csvLines, nil
}
