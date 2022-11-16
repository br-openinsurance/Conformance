package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func searchFileInTable(table [][]string, file map[string]string) int {
	for i, row := range table {
		if row[0] == file["Organisation"] && row[1] == file["Deployment"] {
			return i
		}
	}
	return -1
}

func exportTable(table [][]string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Failed to create csv file: ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range table {
		if err := writer.Write(row); err != nil {
			log.Fatal("Failed to write to file: ", err)
		}
	}
}
