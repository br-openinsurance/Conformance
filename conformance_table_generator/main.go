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
				"consents_v2.3", "financial-risk_v1.2",
				"patrimonial_v1.3", "customers-personal_v1.4",
				"resources_v2.3", "responsibility_v1.2",
			}
		case "legacy":
			apis = []string{
				"acceptance-and-branches-abroad_v1.0", "business_v1.0",
				"consents_v1.0", "consents_v2.2", "financial-risk_v1.0",
				"patrimonial_v1.0", "patrimonial_v1.3-old", "personal_v1.0",
				"resources_v1.0", "resources_v1.2", "responsibility_v1.0",
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: legacy, current", Version)
		}
	case "phase3":
		switch Version {
		case "current":
			apis = []string{ "endorsement_v1.1", "claim-notification_v1.2" }
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: current", Version)
		}
	default:
		log.Fatalf("Invalid target entered: %s. Possible values: phase2, phase3", Target)
	}

	utils.GenerateTable(apis, Target, Version)
}