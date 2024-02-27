package utils

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSVFile(filename string) ([][]string, error) {
	csvFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}

	if err := csvFile.Close(); err != nil {
		return nil, fmt.Errorf("error closing file: %w", err)
	}

	return csvLines, nil
}
