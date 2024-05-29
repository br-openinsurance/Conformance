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
	var apis []string

	switch Target {
	case "phase2":
		switch Version {
		case "current":
			apis = []string{
				"acceptance-and-branches-abroad_v1.2", "customers-business_v1.4",
				"patrimonial_v1.3", "customers-personal_v1.4",
				"resources_v2.3", "responsibility_v1.2", "auto_v1.3", "rural_v1.3", "transport_1.2",
				"consents_v2.5","resources_2.4","customers-business_v1.5","customers-personal_v1.5",
				"acceptance-and-branches-abroad_v1.3","auto_v1.3","housing_v1.3","patrimonial_v1.4",
				"transport_v1.2","responsibility_v1.3","financial-risk_v1.3.1","rural_v1.3",
				"insurance-pension-plan_v1.4","insurance-capitalization-title_v1.4",
				"insurance-financial-assistance_v1.2","insurance-person_v1.5",
				"insurance-life-pension_v1.4"
			}
		case "legacy":
			apis = []string{
				"acceptance-and-branches-abroad_v1.0", "business_v1.0",
				"consents_v2.3", "consents_v2.2", "financial-risk_v1.2",
				"patrimonial_v1.0", "patrimonial_v1.3-old", "personal_v1.0",
				"resources_v1.0", "resources_v1.2", "responsibility_v1.0"
				
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: legacy, current", Version)
		}
	case "phase3":
		switch Version {
		case "current":
			apis = []string{ "endorsement_v1.1.3", "claim-notification-damages_v1.2.3","claim-notification-person_v1.2.3","quote-patrimonial-home_v1.8.1" }
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: current", Version)
		}
	default:
		log.Fatalf("Invalid target entered: %s. Possible values: phase2, phase3", Target)
	}

	utils.GenerateTable(apis, Target, Version)
}

