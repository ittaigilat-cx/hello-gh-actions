package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/checkmarxdev/vorpal/internal"
)

const usage = `
Usage: vorpal -s <sourcePath> [-r <resultsFile>]

Options:
	-s, --source <sourcePath>  	Source path to analyze.
	-r, --result <resultsFile>	Result file to write the output. Supported result file extensions: .json or .csv
	-h, --help                	Show help.
`

var usageDebug string = ""
var version = "0.0.0"

func main() {

	//debug.SetGCPercent(50)
	//fmt.Println("Old GC Percent: ", oldGCPercent)

	//l, i, _ := internal.UniqueLinesInFirstFile("c:\\tmp\\6559.txt", "c:\\tmp\\6458.txt")
	//fmt.Println(l, i)
	//internal.TestParserById("C:\\tmp\\a.cs", "82f82539-391f-4115-b978-59ea183fe06f")
	//.TestParser(`C:\tmp\a.js`, "Hardcoded File Path", ".js")

	if len(usageDebug) == 0 {
		fmt.Println()
		green := "\033[32m"
		reset := "\033[0m"

		fmt.Printf("%sVorpal - %s - Code Security Assistant\n", green, version)
		fmt.Printf("%s", reset)
		fmt.Println("\"Sharp and Deadly\" - The Ultimate Weapon for Developers")
		fmt.Println("\nSupported Languages: Java (.java), C# (.cs), Go (.go), Python (.py), Node.js (.js)")
		fmt.Println("\n©2024 Checkmarx Ltd. All Rights Reserved.")
	}

	var help bool
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&help, "h", false, "")

	sourcePath := ""
	flag.StringVar(&sourcePath, "s", "", "")
	flag.StringVar(&sourcePath, "source", "", "")

	resultsFile := ""
	flag.StringVar(&resultsFile, "r", "", "")
	flag.StringVar(&resultsFile, "result", "", "")

	flag.Usage = func() {
		fmt.Print(usage)
		if len(usageDebug) > 0 {
			fmt.Println(usageDebug)
		} else {
			fmt.Println()
		}
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(1)
	}

	rulesService := internal.NewRulesService()

	if internal.ListRules {
		rulesService.PrintQueriesList()
		os.Exit(0)
	}

	if len(sourcePath) == 0 {
		fmt.Printf("\nMissing source path command argument\n")
		flag.Usage()
		os.Exit(1)
	}

	//check if the results file extension is .json or .csv
	rFile := strings.ToLower(resultsFile)
	if filepath.Ext(rFile) != ".json" && filepath.Ext(rFile) != ".csv" {
		if len(resultsFile) > 0 {
			fmt.Printf("\nInvalid results file '%s'. The file should have an extension of .json or .csv\n", resultsFile)
		}
		fmt.Printf("\nResults will be written to the default location: '%s'\n", internal.DefaultResultsFile)
		resultsFile = internal.DefaultResultsFile
	}

	_, err := os.Stat(resultsFile)
	if !os.IsNotExist(err) {
		os.Remove(resultsFile)
	}

	start := time.Now()

	var queries []internal.Query
	if internal.RuleIdToEvaluate > 0 {
		query := rulesService.GetQueriesByRuleId(internal.RuleIdToEvaluate)
		queries = append(queries, query)
	} else {
		queries = rulesService.GetQueries()
	}

	var results []internal.Result
	err1 := internal.Analyze(sourcePath, &results, queries)
	if err1 == nil {
		if internal.LineCounter > 0 {
			fmt.Println("Lines of Code Analyzed:", internal.LineCounter)
		}

		results, er := internal.WriteToLog(resultsFile, results)
		elapsed := time.Since(start)
		elapsedPer1M := (elapsed.Seconds() * 1000000.0) / float64(internal.LineCounter)
		if internal.Verbose {
			fmt.Printf("Processing Time: %s (%.3f seconds per 1,000,000 lines) \n", elapsed, elapsedPer1M)
		} else {
			fmt.Printf("Processing Time: %.1fs\n", elapsed.Seconds())
		}

		fmt.Println("Security Issues Found:", len(results))
		if er != nil {
			red := "\033[31m"
			reset := "\033[0m" // ANSI escape code to reset color
			fmt.Println(red+"Failed to open the result file", resultsFile, reset)
		} else if len(results) > 0 && len(resultsFile) > 0 {

			cleanedPath := filepath.Clean(resultsFile)
			absPath, err := filepath.Abs(cleanedPath)
			if err == nil {
				fmt.Println("\nResults can be found in:", absPath)
			}
		} else if len(results) == 0 {
			fmt.Println("No security issues found.")
		}
	} else {
		fmt.Println("")
	}
	fmt.Println()

}
