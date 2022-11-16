package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/br-openinsurance/Conformance/tree/main/conformance_table_generator/models"
	"github.com/joho/godotenv"
)

func importData(baseUrl string, path string, data models.GithubTree) (models.GithubTree, error) {
	var url string
	// no previous data means we are searching for root
	if data == nil {
		url = baseUrl + path
	} else if _, err := strconv.Atoi(path); err == nil {
		// if path is numeric, it means that it is a version
		url, _ = findShaByVersion(data, path)
		url = baseUrl + url
	} else {
		// normal path
		url = baseUrl + findShaByPath(data, path)
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if err := godotenv.Load(".env"); err == nil {
		token := "Bearer " + os.Getenv("GITHUB_AT")
		req.Header.Add("Authorization", token)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	var respBody models.GithubResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody.Tree, nil
}

func navigateFromRootToFolder(baseUrl string, path string) models.GithubTree {
	var data models.GithubTree
	// Get root of the repository
	data, err := importData(baseUrl, "main", data)
	if err != nil {
		log.Fatalf("Failed to import data from %s path from repo: %s", path, err)
	}
	// If path is empty, you just want the root
	if path == "" {
		return data
	}

	pathSequence := strings.Split(path, "/")

	for _, path := range pathSequence {
		data, err = importData(baseUrl, path, data)
		if err != nil {
			log.Fatal("Failed to import data from github API: ", err)
		}
	}

	return data
}

func findShaByPath(data models.GithubTree, path string) string {
	for _, element := range data {
		if element.Path == path {
			return element.Sha
		}
	}
	return ""
}

// used for both getting the Sha to make a new request and getting the exact version name
func findShaByVersion(data models.GithubTree, version string) (string, string) {
	pattern, err := regexp.Compile(version + `.\d+.\d+`)
	if err != nil {
		log.Fatal("Failed to compile regular expression: ", err)
	}

	for _, element := range data {
		if pattern.MatchString(element.Path) {
			return element.Sha, element.Path
		}
	}
	return "", ""
}

func buildZipUrlFromBaseUrl(baseUrl string, api string, exactVersion string, fileName string) string {
	old := "api.github.com/repos"
	new := "github.com"
	zipUrl := strings.Replace(baseUrl, old, new, 1)

	old = "git/trees/"
	new = "tree/main/submissions/functional/"
	zipUrl = strings.Replace(zipUrl, old, new, 1)

	zipUrl += api + "/" + exactVersion + "/" + fileName

	return zipUrl
}

func getEveryFileForApiAndVersion(baseUrl string, version string, api string, baseData models.GithubTree) []map[string]string {
	apiData, err := importData(baseUrl, api, baseData)
	if err != nil {
		log.Fatalf("Failed to import api %s from repo: %s", api, err)
	}

	_, exactVersion := findShaByVersion(apiData, version)
	data, err := importData(baseUrl, version, apiData)
	if err != nil {
		log.Fatalf("Failed to import version %s from api %s from repo: %s", version, api, err)
	}

	var output []map[string]string
	for _, file := range data {
		// remove .zip from the end of the file
		fileName := file.Path[:len(file.Path) - 4]
		fileInformation := strings.Split(fileName, "_")
		infoLen := len(fileInformation)

		orgId          := fileInformation[0]
		if len(orgId) > 8 {
			orgId = orgId[:8]
		}
		deploymentName := strings.Join(fileInformation[1 : infoLen - 3], " ")
		apiName        := fileInformation[infoLen - 3]
		versionName    := fileInformation[infoLen - 2]
		submissionDate := fileInformation[infoLen - 1]
		zipUrl         := buildZipUrlFromBaseUrl(baseUrl, api, exactVersion, file.Path)

		organisationsMap := makeOrganisationsMap()

		output = append(output, map[string]string{
			"Org Id":         orgId,
			"Organisation":   organisationsMap[orgId],
			"Deployment":     deploymentName,
			"API":            apiName,
			"Version":        versionName,
			"Date":           submissionDate,
			"Zip URL":        zipUrl,
		})
	}

	return output
}

func makeOrganisationsMap() map[string]string {
	organisations := make(map[string]string)

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://data.directory.opinbrasil.com.br/roles", nil)
	if err != nil {
		log.Fatal("Failed to make a request to roles endpoint: ", err)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to obtain a response from roles endpoint: ", err)
	}

	defer resp.Body.Close()

	var roles models.Roles
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		log.Fatal("Failed to decode response from roles endpoint: ", err)
	}

	orgsWithParent := make(map[string]string)
	for _, role := range roles {
		if role.ParentOrganisationReference != nil && role.ParentOrganisationReference != role.RegistrationNumber {
			orgsWithParent[role.RegistrationNumber[:8]] = fmt.Sprintf("%v", role.ParentOrganisationReference)[:8]
		} else {
			organisations[role.RegistrationNumber[:8]] = role.RegisteredName
		}
	}
	for child, parent := range orgsWithParent {
		organisations[child] = organisations[parent]
	}

	return organisations
}
