# validate-http-headers
Utility to validate HTTP headers of multiple URLs based on a JSON spec. Useful
to quickly set default request headers and automate network deployment testing.

Write a JSON spec file for the validator to use: e.g. urls.json

Pass it to the validator

	$ ./validate-http-headers urls.json

    ------------------------------------------------------------------------
    FILE: urls.json
    
    URL: https://www.google.com/
    SUCCESS: Expected 'X-Frame-Options' to have 'SAMEORIGIN' match 'SAMEORIGIN'
    SUCCESS: Expected 'Expires' to have '-1' match '-1'
    
    URL: https://drive.google.com/
    SUCCESS: Response header 'X-Frame-Options' should not exist

Multiple JSON specs can be passed

	$ ./validate-http-headers internal.json external.json mtv.json

Simple spec:

    {
      "default": {
        "requestHeaders": {
          "Referer": ["https://bad-website.com"]
        },
        "responseHeaders": {
          "X-Frame-Options": ["SAMEORIGIN"]
        }
      },
    
      "specs": [
        {
          "url": "https://www.google.com/"
        },
        {
          "url": "https://drive.google.com/"
        }
      ]
    }

Complex spec with overrides:

    {
      "default": {
        "requestHeaders": {
          "Referer": ["https://golang.org"]
        },
        "responseHeaders": {
          "X-Frame-Options": ["SAMEORIGIN"]
        }
      },
    
      "specs": [
        {
          "url": "https://www.google.com/",
          "requestHeaders": {
            "Referer": ["https://drive.google.com"]
          },
          "responseHeaders": {
            "Expires": ["-1"]
          }
        },
        {
          "url": "https://drive.google.com/",
          "responseHeaders": {
            "X-Frame-Options": [""]
          }
        }
      ]
    }

