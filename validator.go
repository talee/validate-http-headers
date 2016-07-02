// Utility to validate HTTP headers of multiple URL based on a JSON spec
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Spec struct {
	Url             string
	RequestHeaders  map[string][]string
	ResponseHeaders map[string][]string
}

type SpecContainer struct {
	Default Spec
	Specs   []Spec
}

func main() {
	var filename string
	if len(os.Args) == 2 {
		filename = os.Args[1]
	} else {
		filename = "urls.json"
	}
	// Read file
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	// Marshal from JSON
	var specContainer SpecContainer
	e = json.Unmarshal(file, &specContainer)
	if e != nil {
		fmt.Printf("Unmarshal error: %v\n", e)
	}
	specs := specContainer.Specs
	defaultSpec := specContainer.Default

	// Request each URL in specs
	for _, spec := range specs {
		fmt.Println("URL:", spec.Url)
		resp, err := http.Get(spec.Url)
		if err != nil {
			log.Fatalf("http.Get => %v", err.Error())
		}

		// Validate headers
		expectedHeaders := clone(defaultSpec.ResponseHeaders,
			spec.ResponseHeaders)

		for key, expectedHeaderValues := range expectedHeaders {
			numResponseHeaderValues := len(resp.Header[key])
			if numResponseHeaderValues == 0 ||
				numResponseHeaderValues != len(expectedHeaderValues) {
				log.Fatalf("Missing response header %v", key)
			}
			for i, expectedHeaderValue := range expectedHeaderValues {
				respHeaderValue := resp.Header[key][i]
				if expectedHeaderValue != respHeaderValue {
					log.Fatalf("Header assertion failed: Expected '%v'"+
						"to have '%v' instead of '%v'", key,
						expectedHeaderValue, respHeaderValue)
				}
			}
		}
	}
}

func clone(maps ...map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for _, m := range maps {
		for key, val := range m {
			if val != nil {
				result[key] = val
			}
		}
	}
	return result
}
