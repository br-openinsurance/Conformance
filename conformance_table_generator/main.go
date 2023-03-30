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
	flag.StringVar(&Version, "v", "latest", "API Versions")
	flag.Parse()
}

// go run main.go -t <phaseNo> -v <versionGroup>
// example
// go run main.go -t phase2 -v first

// Add a .env file containing your github access token to the root of this project
// GITHUB_AT=ghp_...
func main() {
	var apis []string

	switch Target {
	case "phase2":
		switch Version {
		case "latest":
			apis = []string{
				"acceptance-and-branches-abroad_v1.2", "business_v1.3",
				"consents_v2.2", "financial-risk_v1.2",
				"patrimonial_v1.3", "personal_v1.3",
				"resources_v1.2", "responsibility_v1.2",
			}
		case "first":
			apis = []string{
				"acceptance-and-branches-abroad_v1.0", "business_v1.0",
				"consents_v1.0", "financial-risk_v1.0",
				"patrimonial_v1.0", "personal_v1.0",
				"resources_v1.0", "responsibility_v1.0",
			}
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: first, latest", Version)
		}
	case "phase3":
		switch Version {
		case "latest":
			apis = []string{ "endorsement_v1.1", "claim-notification_v1.2" }
		default:
			log.Fatalf("Invalid version entered: %s. Possible values: latest", Version)
		}
	default:
		log.Fatalf("Invalid target entered: %s. Possible values: phase2, phase3", Target)
	}

	utils.GenerateTable(apis, Target, Version)
}