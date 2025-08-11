package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	alphabet "github.com/rh2g17/md-brasilian-alphabet-sort"

	"github.com/br-openinsurance/Conformance/tree/main/export_data/models"
)

func ImportData(url string) (models.Participants, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	var data models.Participants
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}

	return data, nil
}

func FindIndex(arr []string, str string) int {
	for i, s := range arr {
		if s == str {
			return i
		}
	}
	return -1
}

func MakeOrganisationsMap() map[string]string {
	organisations := make(map[string]string)

	// Requesting data from the API roles endpoint
	roles, err := ImportData("https://data.directory.opinbrasil.com.br/roles")
	if err != nil {
		log.Fatal("Failed to request data from the roles API:", err)
	}

	for _, role := range roles {
		organisations[role.RegistrationNumber] = role.RegisteredName
	}

	return organisations
}

func IsRightVersion(apiVersion string, targetVersion string) bool {
	return strings.Join(strings.Split(apiVersion, ".")[:2], ".") == targetVersion
}

// Copied from the main.go file in the root directory
func GenerateFromCsv(inputFile string, outputFile string, headers []string, separator rune) {
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer f.Close()

	// Read lines from file
	reader := csv.NewReader(f)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	lines = lines[1:]
	var sortLines []string

	for _, line := range lines {
		joinedLine := strings.Join(line, ",")
		sortLines = append(sortLines, joinedLine)
	}
	sortLines = alphabet.MergeSort(sortLines)
	lines = [][]string{}

	for _, item := range sortLines {
		split := strings.Split(item, ",")
		lines = append(lines, split)
	}

	// Set the table to output as a string
	tableOutput := &strings.Builder{}
	table := tablewriter.NewWriter(tableOutput)

	var indexHeaders []string

	for index := range headers {
		indexHeaders = append(indexHeaders, strconv.Itoa(index))
	}

	// Configure table
	table.SetHeader(indexHeaders)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("|")
	table.AppendBulk(lines)
	table.Render()

	// Open output file
	output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	toWrite := tableOutput.String()

	// Replace headersxs
	for index, value := range headers {
		toWrite = strings.Replace(toWrite, " "+strconv.Itoa(index), " "+value, 1)
	}

	// Write result of table to file
	output.Write([]byte(toWrite))
	output.Close()
}

func FilterEntriesWithoutConsents(inputFile string, separator rune) {
	fileRead, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	// Read lines from file
	reader := csv.NewReader(fileRead)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	fileRead.Close()

	// Reopen file
	fileWrite, err := os.Create(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer fileWrite.Close()

	// Create writer
	writer := csv.NewWriter(fileWrite)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(lines[0]); err != nil {
		log.Fatal("Failed to write to file: ", err)
	}
	lines = lines[1:]

	// Writing the rest of the lines filtering entries without data in consents api column
	for _, line := range lines {
		if line[2] != "" {
			if err := writer.Write(line); err != nil {
				log.Fatal("Failed to write to file: ", err)
			}
		}
	}
}

func FilterEntriesWithoutApiData(inputFile string, separator rune) {
	fileRead, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open file for reading: %v", err)
	}
	defer fileRead.Close()

	reader := csv.NewReader(fileRead)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read csv file: %v", err)
	}

	if len(lines) == 0 {
		// Nothing to do
		return
	}

	// Recreate file (truncate)
	fileWrite, err := os.Create(inputFile)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %v", err)
	}
	defer fileWrite.Close()

	writer := csv.NewWriter(fileWrite)
	writer.Comma = separator
	defer func() {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Fatalf("Failed flushing csv writer: %v", err)
		}
	}()

	// Write header back
	if err := writer.Write(lines[0]); err != nil {
		log.Fatalf("Failed to write header: %v", err)
	}

	// Process data rows
	var (
		removed        int
		kept           int
		removedEntries []string
	)

	for _, row := range lines[1:] {
		// Defensive: ensure row has at least the two mandatory columns
		if len(row) < 2 {
			// malformed, drop
			removed++
			removedEntries = append(removedEntries, fmt.Sprintf("MALFORMED ROW: %v", row))
			continue
		}

		hasApiData := false
		for i := 2; i < len(row); i++ {
			v := strings.TrimSpace(row[i])
			if v != "" {
				// If you want to treat "No date" as *empty*, then uncomment:
				// if strings.EqualFold(v, "No date") {
				//     continue
				// }
				hasApiData = true
				break
			}
		}

		if hasApiData {
			if err := writer.Write(row); err != nil {
				log.Fatalf("Failed to write row: %v", err)
			}
			kept++
		} else {
			removed++
			// Log which organization is being removed
			conglomerado := row[0]
			marca := row[1]
			removedEntries = append(removedEntries, fmt.Sprintf("%s | %s", conglomerado, marca))
		}
	}

	log.Printf("FilterEntriesWithoutApiData: kept=%d removed=%d (total=%d)", kept, removed, kept+removed)

	if len(removedEntries) > 0 {
		log.Printf("REMOVED ENTRIES (no API data):")
		for _, entry := range removedEntries {
			log.Printf("  - %s", entry)
		}
	}
}

func FilterDuplicateEntries(inputFile string, separator rune) {
	fileRead, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}

	// Read lines from file
	reader := csv.NewReader(fileRead)
	reader.Comma = separator
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv file: ", err)
	}
	fileRead.Close()

	// Reopen file
	fileWrite, err := os.Create(inputFile)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer fileWrite.Close()

	// Create writer
	writer := csv.NewWriter(fileWrite)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(lines[0]); err != nil {
		log.Fatal("Failed to write to file: ", err)
	}
	lines = lines[1:]

	// Keep only one instance of exactly equal lines
	isWritten := make(map[string]bool)
	for _, line := range lines {
		if _, written := isWritten[strings.Join(line, " ")]; !written {
			isWritten[strings.Join(line, " ")] = true
			if err := writer.Write(line); err != nil {
				log.Fatal("Failed to write to file: ", err)
			}
		}
	}
}

func DateFromZipName(zip string) string {
	re, err := regexp.Compile(`\d+-\D{3}-\d{4}`)
	if err != nil {
		log.Fatal("Could not create regular expression: ", err)
	}

	return re.FindString(zip)
}

func ConvertDate(date string) string {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		log.Fatalf("Failed to convert date (%s): %s", date, err)
	}

	return t.Format("02-Jan-2006")
}
