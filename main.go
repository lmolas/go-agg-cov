package main

import (
	"bufio"
	"log"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/tools/cover"
)

type coverage struct {
	Statements        int
	StatementsCovered int
	PercentageCovered float64
}

func main() {
	var coverFile, businessLogicFile string

	app := cli.NewApp()
	app.Name = "go-agg-cov"
	app.Usage = "Calculates a single coverage percentage from a go coverage file and a an optional list of go files representing the business logic"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "coverFile",
			Value:       "",
			Usage:       "coverage file to scan",
			Destination: &coverFile,
		},
		cli.StringFlag{
			Name:        "businessLogicFile",
			Value:       "",
			Usage:       "file containing list of business logic files",
			Destination: &businessLogicFile,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Analyzing file %s\n", coverFile)
		log.Printf("Business Logic file %s\n", businessLogicFile)
		if coverFile == "" {
			log.Fatal("CoverFile is mandatory")
		}

		business, errBusiness := parseBusinessLogicFile(businessLogicFile)
		if errBusiness != nil {
			log.Fatal(errBusiness)
		}

		profiles, err := cover.ParseProfiles(coverFile)
		if err != nil {
			log.Fatal(err)
		}

		coverage := &coverage{}
		coverage.calculateCoverage(profiles, business)

		log.Printf("Nb Statements: %d Coverage percentage: %f %%", coverage.Statements, coverage.PercentageCovered)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *coverage) calculateCoverage(profiles []*cover.Profile, businessFiles map[string]struct{}) {

	if len(businessFiles) == 0 {
		// Calculate coverage on all coverage file
		for _, profile := range profiles {
			for _, block := range profile.Blocks {
				c.Statements += block.NumStmt
				if block.Count > 0 {
					c.StatementsCovered += block.NumStmt
				}
			}
		}
	} else {
		// Calculate coverage on business files from coverage file
		for _, profile := range profiles {
			_, isBusinessFile := businessFiles[profile.FileName]
			if isBusinessFile {
				for _, block := range profile.Blocks {
					c.Statements += block.NumStmt
					if block.Count > 0 {
						c.StatementsCovered += block.NumStmt
					}
				}
			}
		}
	}

	c.PercentageCovered = (float64(c.StatementsCovered) / float64(c.Statements)) * 100.0
}

func parseBusinessLogicFile(businessLogicFile string) (map[string]struct{}, error) {
	files := make(map[string]struct{}, 0)

	if businessLogicFile == "" {
		return files, nil
	}
	pf, err := os.Open(businessLogicFile)
	if err != nil {
		return nil, err
	}
	defer pf.Close()

	buf := bufio.NewReader(pf)
	s := bufio.NewScanner(buf)
	for s.Scan() {
		files[s.Text()] = struct{}{}
	}

	return files, nil
}
