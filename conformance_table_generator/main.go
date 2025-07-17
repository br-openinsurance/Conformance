package main

import (
	"flag"
	"log"

	"github.com/br-openinsurance/Conformance/tree/main/conformance_table_generator/utils"
)

var (
	Target  string
	Version string
)

func init() {
	flag.StringVar(&Target, "t", "phase2", "Target Table")
	flag.StringVar(&Version, "v", "current", "API Versions")
	flag.Parse()
}

// go run main.go -t <phaseNo> -v <versionGroup>
// example
// go run main.go -t phase2 -v current

// Add a .env file containing your github access token to the root of this project
// GITHUB_AT=ghp_...
func main() {
	log.Println("Starting table generation process...")
	log.Printf("Target phase: %s", Target)
	log.Printf("Version group: %s", Version)

	var apis []string

	switch Target {
	case "phase2":
		switch Version {
		case "current":
			apis = []string{
				"acceptance-and-branches-abroad_v1.3",
				"customers-business_v1.5",
				"patrimonial_v1.4",
				"customers-personal_v1.5",
				"resources_v2.4",
				"responsibility_v1.3",
				"auto_v1.3",
				"rural_v1.3",
				"transport_1.2",
				"consents_v2.6",
				"resources_2.4",
				"customers-business_v1.5",
				"customers-personal_v1.5",
				"acceptance-and-branches-abroad_v1.3",
				"auto_v1.3",
				"housing_v1.3",
				"patrimonial_v1.4",
				"transport_v1.2",
				"responsibility_v1.3",
				"financial-risk_v1.3.1",
				"rural_v1.3",
				"insurance-pension-plan_v1.4",
				"insurance-capitalization-title_v1.4",
				"insurance-financial-assistance_v1.2",
				"insurance-person_v1.6",
				"insurance-life-pension_v1.4",
			}
		case "legacy":
			apis = []string{
				"acceptance-and-branches-abroad_v1.0",
				"business_v1.0",
				"consents_v2.3",
				"consents_v2.2",
				"financial-risk_v1.2",
				"patrimonial_v1.0",
				"patrimonial_v1.3-old",
				"personal_v1.0",
				"resources_v1.0",
				"resources_v1.2",
				"responsibility_v1.0",
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: legacy, current", Version)
		}
	case "phase3":
		switch Version {
		case "current":
			apis = []string{
				"claim-notification-damages_v1.3.0",
				"claim-notification-person_v1.3.0",
				"endorsement_v1.2.0",
				"quote-patrimonial-home_v1.10.0",
				"quote-acceptance-and-branches-abroad_v1.8.0",
				"quote-auto_v1.9.0",
				"quote-financial-risk_v1.8.0",
				"quote-housing_v1.8.0",
				"quote-responsibility_v1.8.0",
				"quote-rural_v1.8.0",
				"quote-transport_v1.8.0",
				"contract-life-pension_v1.13.0",
				"withdrawal-pension_v1.3.0",
				"withdrawal-capitalization-title_v1.3.0",
				"quote-person-life_v1.11.0",
				"quote-person-travel_v1.11.0",
				"quote-capitalization-title_v1.10.0"
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: current", Version)
		}
	default:
		log.Fatalf("Invalid target entered: %s. Possible values: phase2, phase3", Target)
	}

	log.Printf("APIs to be processed: %v", apis)
	utils.GenerateTable(apis, Target, Version)
	log.Println("Table generation completed successfully!")
}

