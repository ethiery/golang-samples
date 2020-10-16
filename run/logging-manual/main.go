// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample logging-manual shows how to leverage Cloud Run structured logging without a client library.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("Defaulting to port", port)
	}

	// Start HTTP server.
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// [START run_manual_logging_object]

// Entry defines a log entry.
type Entry struct {
	Severity       string                 `json:"severity,omitempty"`
	Message        string                 `json:"message"`
	HTTPRequest    HTTPRequest            `json:"httpRequest"`
	Timestamp      time.Time              `json:"time"`
	InsertID       string                 `json:"logging.googleapis.com/insertId"`
	Labels         map[string]interface{} `json:"logging.googleapis.com/labels"`
	Operation      string                 `json:"logging.googleapis.com/operation"`
	SourceLocation SourceLocation         `json:"logging.googleapis.com/sourceLocation"`
	SpanID         string                 `json:"logging.googleapis.com/spanId"`
	Trace          string                 `json:"logging.googleapis.com/trace"`
	TraceSampled   bool                   `json:"logging.googleapis.com/trace_sampled"`
}

type HTTPRequest struct {
	Method                         string `json:"requestMethod"`
	URL                            string `json:"requestUrl"`
	Size                           string `json:"requestSize"`
	Status                         int    `json:"status"`
	ResponseSize                   string `json:"responseSize"`
	UserAgent                      string `json:"userAgent"`
	RemoteIP                       string `json:"remoteIp"`
	ServerIP                       string `json:"serverIp"`
	Referer                        string `json:"referer"`
	Latency                        string `json:"latency"`
	CacheLookup                    bool   `json:"cacheLookup"`
	CacheHit                       bool   `json:"cacheHit"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer"`
	CacheFillBytes                 string `json:"cacheFillBytes"`
	Protocol                       string `json:"protocol"`
}

type SourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}

// String renders an entry structure to the JSON format expected by Cloud Logging.
func (e Entry) String() string {
	if e.Severity == "" {
		e.Severity = "INFO"
	}
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}

// [END run_manual_logging_object]
// [START run_manual_logging]

func init() {
	// Disable log prefixes such as the default timestamp.
	// Prefix text prevents the message from being parsed as JSON.
	// A timestamp is added when shipping logs to Cloud Logging.
	log.SetFlags(0)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(Entry{
		Severity: "NOTICE",
		Message:  "This is the default display field.",
		HTTPRequest: HTTPRequest{
			Method:                         "POST",
			URL:                            "https://myapi.com",
			Size:                           "1234",
			Status:                         200,
			ResponseSize:                   "5678",
			UserAgent:                      "UserAgent",
			RemoteIP:                       "192.168.1.1",
			ServerIP:                       "192.168.1.1",
			Referer:                        "https://referer.com",
			Latency:                        "3.5s",
			CacheLookup:                    true,
			CacheHit:                       true,
			CacheValidatedWithOriginServer: true,
			CacheFillBytes:                 "31415",
			Protocol:                       "HTTP/2",
		},
		Timestamp: time.Date(2020, time.October, 16, 21, 22, 23, 24, time.UTC),
		InsertID:  "123456",
		Labels: map[string]interface{}{
			"key1": "value",
			"key2": "42",
		},
		Operation: "operation",
		SourceLocation: SourceLocation{
			File:     "main.go",
			Line:     "132",
			Function: "indexHandler",
		},
		SpanID:       "000000000000004a",
		Trace:        "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824",
		TraceSampled: true,
	})

	fmt.Fprintln(w, "Hello Logger!")
}

// [END run_manual_logging]
