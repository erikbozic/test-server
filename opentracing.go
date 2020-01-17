package main

import (
	"net/http"
	"strings"
)

func copyXb3Headers(source *http.Request, dest *http.Request) {
	for key, values := range source.Header {
		if strings.HasPrefix(strings.ToLower(key), "x-b3-")  {
			for _, v := range values {
				dest.Header.Add(key, v)
			}
		}
	}
}
