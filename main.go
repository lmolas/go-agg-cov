package main

import (
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
	var coverFile string

	app := cli.NewApp()
	app.Name = "go-agg-cov"
	app.Usage = "Calculates a single coverage percentage from a go coverage file"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "coverFile",
			Value:       "",
			Usage:       "coverage file to scan",
			Destination: &coverFile,
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Printf("Analyzing file %s\n", coverFile)
		if coverFile == "" {
			log.Fatal("CoverFile is mandatory")
		}

		profiles, err := cover.ParseProfiles(coverFile)
		if err != nil {
			log.Fatal(err)
		}

		coverage := &coverage{}
		coverage.calculateCoverage(profiles)

		log.Printf("Nb Statements: %d Coverage percentage: %f %%", coverage.Statements, coverage.PercentageCovered)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *coverage) calculateCoverage(profiles []*cover.Profile) {

	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			c.Statements += block.NumStmt
			if block.Count > 0 {
				c.StatementsCovered += block.NumStmt
			}
		}
	}

	c.PercentageCovered = (float64(c.StatementsCovered) / float64(c.Statements)) * 100.0
}
