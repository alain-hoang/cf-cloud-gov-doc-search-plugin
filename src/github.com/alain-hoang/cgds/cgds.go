package main

import (
	"fmt"
	"os/exec"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/PuerkitoBio/goquery"
	"runtime"
	"os"
	"flag"
	"encoding/json"
)

type CloudGovDocSearchPlugin struct{}


type SearchResult struct {
	Url      string      `json:"url"`
	Descr    string      `json:"description"`
}

const QUERY_URL = "https://search.usa.gov/search?affiliate=cloud.gov&utf8=%E2%9C%93&query="
const OPTIONS_HELP_STR = "[-format <human|json>] [-url <uri>] <search-term>"

func getResults(searchTerm string, queryUrl string) []SearchResult {
	doc, err := goquery.NewDocument(queryUrl + searchTerm)
	if err != nil {
		fmt.Printf("Error reading %v\n", err)
		return nil
	}

	return parseResults(doc)
}

func parseResults(doc *goquery.Document) []SearchResult {
	var res []SearchResult

	doc.Find(".content-block-item.result").Each(func(i int, s *goquery.Selection) {
		s1 := s.Find("a")
		link, ok := s1.Attr("href")
		if ok {
			sr := SearchResult{Url: link, Descr: s1.Text()}
			res = append(res, sr)
		}
	})

	return res
}

func printResults(sr []SearchResult, format string) {
	switch format {
	case "human":
		for idx, r := range sr {
			fmt.Printf("%d | %s | %s\n", idx, r.Url, r.Descr)
		}
	case "json":
		js, err := json.Marshal(sr)
		if err != nil {
			fmt.Printf("Error to json: %v", err)
			return
		}
		fmt.Printf("%s\n", js)
	default:
		fmt.Printf("Don't know how to output that")
	}
}

func chooseFromResults(sr []SearchResult) error {
	var choice int

	fmt.Printf("\nEnter Choice: ")
	_, err := fmt.Scanf("%d\n", &choice)
	if err != nil {
		fmt.Printf("Error entering choice: %v\n", err)
	} else {
		if choice < len(sr) {
			fmt.Printf("Chose %d, going to %s\n", choice, sr[choice].Url)
			err = openUrl(sr[choice].Url)
		} else {
			fmt.Printf("Invalid choice %d\n", choice)
		}
	}

	return err
}

func openUrl(url string) error {
	var err error
	var cmds = map[string][]string{
		"windows": []string{"cmd", "/c", "start"},
		"darwin":  []string{"open"},
		"linux":   []string{"xdg-open"},
	}

	cmd, ok := cmds[runtime.GOOS]
	if ok {
		cmd = append(cmd, url)
		c := exec.Command(cmd[0], cmd[1:]...)
		err = c.Run()
		if err != nil {
			fmt.Printf("Error trying run execute: %v\n", err)
		}
	} else {
		fmt.Printf("Error running on platform %s\n", runtime.GOOS)
	}

	return err
}

func canFormat(format string) bool {
	switch format {
	case "human":
		return true
	case "json":
		return true
	default:
		return false
	}
}



func usage() {
	fmt.Printf("cf cloud-gov-doc-search %s\n", OPTIONS_HELP_STR) 
}

func (c *CloudGovDocSearchPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// Ensure that we called the command basic-plugin-command
	if args[0] == "cloud-gov-doc-search" {
		flags := flag.NewFlagSet("cgds-args", flag.ExitOnError)
		urlPtr := flags.String("url", QUERY_URL, "search url")
		fmtPtr := flags.String("format", "human", "Output format <human|json>")

		err := flags.Parse(args[1:])
		if err != nil {
			usage()
			fmt.Printf("Error parsing:  %v\n", err)
			os.Exit(1)
		}
		if canFormat(*fmtPtr) {
			sq := flags.Arg(0)
			sr := getResults(sq, *urlPtr)
			printResults(sr, *fmtPtr)
			if *fmtPtr == "human" {
				err = chooseFromResults(sr)
				if err != nil {
					fmt.Printf("Error in choosing: %v\n", err)
					os.Exit(1)
				}
			}
		} else {
			usage()
			os.Exit(1)
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
					Usage: "cgds-plugin\n   cf cgds-plugin " + OPTIONS_HELP_STR,
					Options: map[string]string{
						"url": "If this param is set ",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(CloudGovDocSearchPlugin))
}
