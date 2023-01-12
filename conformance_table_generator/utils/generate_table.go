package utils

import (
	"fmt"
	"strings"

	"github.com/br-openinsurance/Conformance/tree/main/conformance_table_generator/models"
)

func GenerateTable(apis []string, phase string, version string) {
	// import files
	repositoryUrl := "https://api.github.com/repos/br-openinsurance/Conformance/git/trees/main?recursive=1"
	submissionFiles := importSubmittedFiles(repositoryUrl)

	// filter files by chosen apis and version
	filteredFiles := filterFilesByApisAndVersion(submissionFiles, apis, version)

	// create table and dump
	tableHeaders := []string {"Conglomerado", "Deployment"}
	tableHeaders = append(tableHeaders, apis...)
	table := [][]string {tableHeaders}

	dumpHeaders := []string {"Id da Organização", "Deployment", "API", "Version", "Data"}
	dump := [][]string {dumpHeaders}

	organisationsMap := makeOrganisationsMap()

	for _, file := range filteredFiles {
		fileSplit := strings.Split(file, "/")
		fileName := fileSplit[len(fileSplit) - 1]
		fileNameSplit := strings.Split(fileName, "_")
		orgId := fileNameSplit[0]
		orgName := organisationsMap[orgId]
		deploymentName := fileNameSplit[1]
		api := fileNameSplit[2]
		version := fileNameSplit[3]
		date := fileNameSplit[4]
		date = date[:len(date) - 4]

		zipUrl := strings.Replace(repositoryUrl, "api.github.com/repos", "github.com", 1)
		zipUrl = strings.Replace(zipUrl, "git/trees/main?recursive=1", "blob/main/" + file, 1)
		zipUrl = strings.Replace(zipUrl, " ", "%20", -1)

		dump = append(dump, []string {
			orgId,
			deploymentName,
			api,
			version,
			date,
		})

		apiIndex := findApiIndex(apis, api)
		if ind := searchFileInTable(table, orgName, deploymentName); ind == -1 {
			newRow := make([]string, len(tableHeaders))
			newRow[0] = orgName
			newRow[1] = deploymentName
			newRow[apiIndex + 2] = fmt.Sprintf("[%s](%s)", date, zipUrl)
			
			table = append(table, newRow)
		} else {
			table[ind][apiIndex + 2] = fmt.Sprintf("[%s](%s)", date, zipUrl)
		}
	}

	dumpFileName := fmt.Sprintf("../results/%s/v%s/%s-v%s-conformance-dump.csv", phase, version, phase, version)
	exportTable(dump, dumpFileName)

	tableFileName := fmt.Sprintf("../results/%s/v%s/%s-v%s-conformance-table.csv", phase, version, phase, version)
	exportTable(table, tableFileName)
}

func filterFilesByApisAndVersion(submissionFiles models.GithubTree, apis []string, version string) []string {
	var filteredFiles []string

	for _, file := range submissionFiles {
		filePath := file.Path
		fileSplit := strings.Split(filePath, "/")
		fileApi := fileSplit[2]
		fileVersion := fileSplit[3]

		if findApiIndex(apis, fileApi) != -1 && strings.Split(fileVersion, ".")[0] == version {
			filteredFiles = append(filteredFiles, filePath)
		}
	}

	return filteredFiles
}

func findApiIndex(apis []string, api string) int {
	for i, element := range apis {
		if element == api {
			return i
		}
	}
	return -1
}
