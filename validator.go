// Utility to validate HTTP headers of multiple URL based on a JSON spec
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const (
	FileError = iota
	UnmarshalError
	InvalidRequest
	FailedRequest
	MissingResponseHeader
	FailAssertResponseHeaderValue
)

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
		os.Exit(FileError)
	}

	// Marshal from JSON
	var specContainer SpecContainer
	e = json.Unmarshal(file, &specContainer)
	if e != nil {
		fmt.Printf("Unmarshal error: %v\n", e)
		os.Exit(UnmarshalError)
	}
	specs := specContainer.Specs
	defaultSpec := specContainer.Default

	// Request each URL in specs
	for _, spec := range specs {
		// Setup request headers
		req, err := http.NewRequest("GET", spec.Url, nil)
		if err != nil {
			fmt.Printf("Failed to create request:", err)
			os.Exit(InvalidRequest)
		}
		for key, headerValues := range spec.RequestHeaders {
			for _, val := range headerValues {
				req.Header.Add(key, val)
			}
		}

		// Send request
		fmt.Println("\nURL:", req.URL)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("http.Get => %v", err.Error())
			os.Exit(FailedRequest)
		}

		// Validate response headers
		expectedHeaders := clone(defaultSpec.ResponseHeaders,
			spec.ResponseHeaders)

		for key, expectedHeaderValues := range expectedHeaders {
			// Assert exist/nonexistence
			numResponseHeaderValues := len(resp.Header[key])
			expectedNumHeaderValues := len(expectedHeaderValues)
			if expectedNumHeaderValues == 1 && expectedHeaderValues[0] == "" {
				expectedNumHeaderValues = 0
			}
			if numResponseHeaderValues != expectedNumHeaderValues {
				fmt.Printf("FAIL: Mismatch of existence for response header "+
					"'%v': Expected %v response headers with array values"+
					"'%v'. Got header values '%v'\n",
					key,
					expectedNumHeaderValues,
					expectedHeaderValues,
					resp.Header[key])
				continue
			} else if expectedNumHeaderValues == 0 {
				fmt.Println("SUCCESS: Response header '"+key+"' should not",
					"exist")
				continue
			}

			// Assert equality
			for i, expectedHeaderValue := range expectedHeaderValues {
				respHeaderValue := resp.Header[key][i]
				if expectedHeaderValue != respHeaderValue {
					fmt.Printf("Header assertion failed: Expected '%v'"+
						"to have '%v' instead of '%v'\n", key,
						expectedHeaderValue, respHeaderValue)
				} else {
					fmt.Println("SUCCESS: Expected", "'"+key+"'",
						"to have", "'"+expectedHeaderValue+"'", "match",
						"'"+respHeaderValue+"'")
				}
			}
		}
	}
}

// Returns a merge of the maps. Nil values ignored
func clone(maps ...map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for _, m := range maps {
		for key, values := range m {
			if values != nil {
				result[key] = values
			}
		}
	}
	return result
}
