package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/br-openinsurance/Conformance/tree/main/export_data/utils"
)

// go run main.go -t <phaseNo> -v <versionGroup>
// example
// go run main.go -t phase2 -v current

var (
	Target            string
	Version           string
)

func init() {
	flag.StringVar(&Target, "t", "phase2", "Target Table")
	flag.StringVar(&Version, "v", "current", "API Versions")
	flag.Parse()
}

func main() {

	resultsPathCsv := fmt.Sprintf("../results/%s/%s/%s-%s-data.csv", Target, Version, Target, Version)
	resultsPathMd  := fmt.Sprintf("../results/%s/%s/%s-%s-data.md" , Target, Version, Target, Version)
	// It's recommended to use semicolon as separator as some organisations might have comma in their names
	separator   := ';'
	var apiFamilyTypes []string
	var apiVersions    []string
	var apiHeaderNames []string
	headers := []string {
		"Conglomerado",
		"Marca",
	}

	switch Target {
	case "phase2":
		apiFamilyTypes = []string {
			"consents",
			"customers-personal",
			"customers-business",
			"resources",
			"insurance-acceptance-and-branches-abroad",
			"financial-risk",
			"patrimonial",
			"responsibility",
		}

		switch Version {
		case "current":
			apiVersions = []string { "2.6", "1.5", "1.5", "2.4", "1.3", "1.3.1", "1.4", "1.3" }
		case "legacy":
			apiVersions = []string { "1.0", "1.0", "1.0", "1.0", "1.0", "1.0", "1.0", "1.0" }
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: legacy, current", Version)
		}

		apiHeaderNames = []string {
			"Consentimento API versão " + apiVersions[0],
			"Dados Cadastrais (PF) API versão " + apiVersions[1],
			"Dados Cadastrais (PJ) API versão " + apiVersions[2],
			"Resources API versão " + apiVersions[3],
			"Aceitação e Sucursal no exterior API versão " + apiVersions[4],
			"Riscos Financeiros API versão " + apiVersions[5],
			"Patrimonial API versão " + apiVersions[6],
			"Responsabilidade API versão " + apiVersions[7],
		}
	case "phase3":
		apiFamilyTypes = []string{
			"claim-notification",
			"dynamic-fields",
			"endorsement",
			"notifications",
			"quote-patrimonial-home",
			"quote-acceptance-and-branches-abroad",
			"quote-auto",
			"quote-financial-risk",
			"quote-housing",
			"quote-responsibility",
			"quote-rural",
			"quote-transport",
			"webhook",
			"contract-life-pension",
			"withdrawal-pension",
			"quote-person-life",
			"quote-person-travel",
			"quote-capitalization-title",
		}

		switch Version {
		case "current":
			apiVersions = []string{
				"1.3",  // claim-notification
				"1.4",  // dynamic-fields
				"1.2",  // endorsement
				"1.0",  // notifications
				"1.10", // quote-patrimonial-home
				"1.8",  // quote-acceptance-and-branches-abroad
				"1.9",  // quote-auto
				"1.8",  // quote-financial-risk
				"1.8",  // quote-housing
				"1.8",  // quote-responsibility
				"1.8",  // quote-rural
				"1.8",  // quote-transport
				"1.1",  // webhook
				"1.13", // contract-life-pension
				"1.3",  // withdrawal-pension
				"1.11", // quote-person-life
				"1.11", // quote-person-travel
				"1.10", // quote-capitalization-title
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: current", Version)
		}

		apiHeaderNames = []string{
			"Claim Notification API versão " + apiVersions[0],
			"Dynamic Fields API versão " + apiVersions[1],
			"Endorsement API versão " + apiVersions[2],
			"Notifications API versão " + apiVersions[3],
			"Quote Patrimonial Home API versão " + apiVersions[4],
			"Quote Acceptance and Branches Abroad API versão " + apiVersions[5],
			"Quote Auto API versão " + apiVersions[6],
			"Quote Financial Risk API versão " + apiVersions[7],
			"Quote Housing API versão " + apiVersions[8],
			"Quote Responsibility API versão " + apiVersions[9],
			"Quote Rural API versão " + apiVersions[10],
			"Quote Transport API versão " + apiVersions[11],
			"Webhook API versão " + apiVersions[12],
			"Contract Life Pension API versão " + apiVersions[13],
			"Withdrawal Pension API versão " + apiVersions[14],
			"Quote Person Life API versão " + apiVersions[15],
			"Quote Person Travel API versão " + apiVersions[16],
			"Quote Capitalization Title API versão " + apiVersions[17],
		}
	default:
		log.Fatalf("Invalid target entered: %s. Possible values: phase2, phase3", Target)
	}

	exportData(apiFamilyTypes, apiHeaderNames, apiVersions, resultsPathCsv, separator)

	// Filter entries that are duplicated
	utils.FilterDuplicateEntries(resultsPathCsv, separator)

	// Specifically for phase 2, we should filter out entries that do not have certification for consents API
	utils.FilterEntriesWithoutApiData(resultsPathCsv, separator)
	if Target == "phase2" {
		utils.FilterEntriesWithoutConsents(resultsPathCsv, separator)
	}

	headers = append(headers, apiHeaderNames...)
	utils.GenerateFromCsv(resultsPathCsv, resultsPathMd, headers, separator)
}

func exportData(apiFamilyTypes []string, apiHeaderNames []string, apiVersions []string, fileName string, separator rune) {
	// Creating the header for the table
	tableHeader := []string{"Conglomerado", "Marca"}
	tableHeader = append(tableHeader, apiHeaderNames...)

	// Requesting data from the API participants endpoint
	participants, err := utils.ImportData("https://data.directory.opinbrasil.com.br/participants")
	if err != nil {
		log.Fatal("Failed to request data from the participants API:", err)
	}

	// Creating the map from registration number to registered name
	organisations := utils.MakeOrganisationsMap()

	// Creating the csv file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Failed to open file: ", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = separator
	defer writer.Flush()

	// Writing header to the file
	if err := writer.Write(tableHeader); err != nil {
		log.Fatal("Failed to write to file:", err)
	}

	// Writing data to the file
	for _, participant := range participants {
		for _, server := range participant.AuthorisationServers {
			rowElements := make(map[string]string)
			// We look for the parent organisation in order to get the conglomerate
			if participant.ParentOrganisationReference != "" {
				rowElements["Conglomerado"] = organisations[participant.ParentOrganisationReference]
			} else {
				rowElements["Conglomerado"] = participant.RegisteredName
			}
			rowElements["Marca"] = server.CustomerFriendlyName

			// Iterate through all servers
			for _, resource := range server.APIResources {
				// The family type must be in apiFamilyTypes and there must be an APICertificationURI
				apiIndex := utils.FindIndex(apiFamilyTypes, resource.APIFamilyType)
				if apiIndex != -1 && resource.APICertificationURI != nil && utils.IsRightVersion(resource.APIVersion, apiVersions[apiIndex]) {
					// Search for the date in the zip containing the certification
					certDate := utils.DateFromZipName(fmt.Sprintf("%v", resource.APICertificationURI))
					if certDate == "" {
						certDate = "No date"
					}
					// If the date is available in the endpoint, it should overwrite the one from the zip
					if resource.CertificationStartDate != nil {
						certDate = utils.ConvertDate(fmt.Sprintf("%v", resource.CertificationStartDate))
					}
					rowElements[resource.APIFamilyType] = fmt.Sprintf(
						"[%s](%s)",
						certDate,
						resource.APICertificationURI,
					)
				}
			}

			row := make([]string, len(apiFamilyTypes) + 2)
			row[0] = rowElements["Conglomerado"]
			row[1] = rowElements["Marca"]
			for i, familyType := range apiFamilyTypes {
				row[i + 2] = rowElements[familyType]
			}

			if err := writer.Write(row); err != nil {
				log.Fatal("Failed to write to file:", err)
			}
		}
	}
}