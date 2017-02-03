package main

import (
	"fmt"
	"os/exec"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/PuerkitoBio/goquery"
)

type CloudGovDocSearchPlugin struct{}


type SearchResult struct {
	Url      string
	Descr    string
}

const queryUrl = "https://search.usa.gov/search?affiliate=cloud.gov&utf8=%E2%9C%93&query="

func getResults(searchTerm string) []SearchResult {
	doc, err := goquery.NewDocument(queryUrl + searchTerm)
	if err != nil {
		fmt.Printf("Error reading %v\n", err)
		return nil
	}

	return parseResults(doc)
}

func parseResults(doc *goquery.Document) []SearchResult {
	res := make([]SearchResult, 0)

	doc.Find(".content-block-item.result").Each(func(i int, s *goquery.Selection) {
		s1 := s.Find("a")
		link, _ := s1.Attr("href")
		sr := SearchResult{Url: link, Descr: s1.Text()}
		res = append(res, sr)
	})

	return res
}

func printResults(sr []SearchResult) {
	for idx, r := range sr {
		fmt.Printf("%d | %s | %s\n", idx, r.Url, r.Descr)
	}
}

func (c *CloudGovDocSearchPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	if args[0] == "cloud-gov-doc-search" {
		if len(args) > 1 {
			var choice int
			sr := getResults(args[1])
			printResults(sr)

			// Get choice and try navigating to choice
			fmt.Printf("\nEnter Choice: ")
			fmt.Scanf("%d", &choice)
			fmt.Printf("Chose %d, going to %s\n", choice, sr[choice].Url)

			cmd := exec.Command("open", sr[choice].Url)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error trying run execute: %v\n", err)
			}
		}
	}
}

func (c *CloudGovDocSearchPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "CloudGovDocSearchPlugin",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "cloud-gov-doc-search",
				Alias:    "cgds",
				HelpText: "Search cloud.gov documentation for search-term",
				UsageDetails: plugin.Usage{
					Usage: "cgds-plugin\n   cf cgds-plugin search-term",
				},
			},
		},
	}
}


func main() {
	plugin.Start(new(CloudGovDocSearchPlugin))
}
