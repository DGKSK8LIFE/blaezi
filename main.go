// Copyright 2021 Aadhav Vignesh

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	// basic flags
	var testsFile = flag.String("tests", "test.json", "Complete path to the tests file.")
	var secure = flag.Bool("secure", false, "Secure connection.")
	var timeout = flag.Int("timeout", 5, "Timeout for client.")

	// parse flags and get args
	flag.Parse()
	args := flag.Args()

	// get the base URL at least
	if len(args) < 1 {
		fmt.Println("No base URL specified. Stopping smoke test.")
		os.Exit(0)
	}

	// throw err if the base URL is invalid
	baseURL, err := extractBaseURL(args[0])
	if err != nil {
		fmt.Println("Invalid base URL specified. Stopping smoke test.")
		os.Exit(0)
	}

	// extract tests from file
	extractedTests, err := extractTests(*testsFile)
	if err != nil {
		fmt.Printf("Error raised while extracting tests: %s\n", err.Error())
		os.Exit(1)
	}

	// TODO: Implement certificate check for HTTPS connections
	if *secure {
		fmt.Println("secure mode.")
	} else {
		fmt.Println("insecure mode.")
	}

	// create a new HTTP client
	client := HTTPClient(*secure, *timeout)

	// create a new inspector object
	inspector := NewInspector(client, baseURL)

	// process using the object
	results, errors := inspector.Test(extractedTests)

	// print result for each test
	for _, result := range results {
		fmt.Println(prettifyResult(result))
	}

	// throw errors if raised (in tests)
	for _, err := range errors {
		fmt.Println(err)
	}

	// Show final stats (prettified)
	fmt.Println(results.countSuccess(), " tests passed out of", len(extractedTests), "tests.")

	// TODO:
	// Show time stats

	// Show failure message if tests fail
	if results.countSuccess() != len(extractedTests) {
		termColorFormat := "\033[031m%s\033[0m"
		fmt.Println(fmt.Sprintf(termColorFormat, "Some tests failed."))
		os.Exit(2)
	}
}
