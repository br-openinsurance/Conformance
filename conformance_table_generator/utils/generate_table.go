package utils

import "fmt"

func GenerateTable(apis []string, phase string, version string) {
	baseUrl := "https://api.github.com/repos/br-openinsurance/Conformance/git/trees/"
	baseData := navigateFromRootToFolder(baseUrl, "submissions/functional")

	tableHeaders := []string {"Conglomerado", "Deployment"}
	tableHeaders = append(tableHeaders, apis...)
	table := [][]string {tableHeaders}

	dumpHeaders := []string {"Id da Organização", "Deployment", "API", "Version", "Data"}
	dump := [][]string {dumpHeaders}

	for i, api := range apis {
		files := getEveryFileForApiAndVersion(baseUrl, version, api, baseData)
		for _, file := range files {
			// dump
			dump = append(dump, []string {
				file["Org Id"],
				file["Deployment"],
				file["API"],
				file["Version"],
				file["Date"],
			})

			// table
			if ind := searchFileInTable(table, file); ind == -1 {
				newRow := make([]string, len(tableHeaders))
				newRow[0] = file["Organisation"]
				newRow[1] = file["Deployment"]
				newRow[i + 2] = fmt.Sprintf("[%s](%s)", file["Date"], file["Zip URL"])
				
				table = append(table, newRow)
			} else {
				table[ind][i + 2] = fmt.Sprintf("[%s](%s)", file["Date"], file["Zip URL"])
			}
		}
	}

	dumpFileName := fmt.Sprintf("../results/%s/v%s/%s-v%s-conformance-dump.csv", phase, version, phase, version)
	exportTable(dump, dumpFileName)

	tableFileName := fmt.Sprintf("../results/%s/v%s/%s-v%s-conformance-table.csv", phase, version, phase, version)
	exportTable(table, tableFileName)
}
