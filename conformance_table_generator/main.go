package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/br-openinsurance/Conformance/tree/main/conformance_table_generator/utils"
)

var (
	Target  string
	Version string
)

func init() {
	flag.StringVar(&Target, "t", "phase2", "Target Table")
	flag.StringVar(&Version, "v", "1", "API Version")
	flag.Parse()
}

// go run main.go -t <phaseNo> -v <versionNo>
// example
// go run main.go -t phase2 -v 2

// Add a .env file containing your github access token to the root of this project
// GITHUB_AT=ghp_...
func main() {
	var apis []string

	if Target == "phase2" || Target == "all" {
		apis = []string{
			"acceptance-and-branches-abroad", "business",
			"consents", "financial-risk",
			"patrimonial", "personal",
			"resources", "responsibility",
		}

		utils.GenerateTable(apis, "phase2", Version)
	}

	if _, err := strconv.Atoi(Version); err != nil {
		log.Fatalf("Invalid version entered: %s. Error: %s", Version, err)
	}
}