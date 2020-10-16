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

package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "no project, no trace",
			want: `{` +
				`"severity":"NOTICE",` +
				`"message":"This is the default display field.",` +
				`"httpRequest":{` +
				`"requestMethod":"POST",` +
				`"requestUrl":"https://myapi.com",` +
				`"requestSize":"1234",` +
				`"status":200,` +
				`"responseSize":"5678",` +
				`"userAgent":"UserAgent",` +
				`"remoteIp":"192.168.1.1",` +
				`"serverIp":"192.168.1.1",` +
				`"referer":"https://referer.com",` +
				`"latency":"3.5s",` +
				`"cacheLookup":true,` +
				`"cacheHit":true,` +
				`"cacheValidatedWithOriginServer":true,` +
				`"cacheFillBytes":"31415",` +
				`"protocol":"HTTP/2"` +
				`},` +
				`"time":"2020-10-16T21:22:23.000000024Z",` +
				`"logging.googleapis.com/insertId":"123456",` +
				`"logging.googleapis.com/labels":{"key1":"value","key2":"42"},` +
				`"logging.googleapis.com/operation":"operation",` +
				`"logging.googleapis.com/sourceLocation":{` +
				`"file":"main.go",` +
				`"line":"132",` +
				`"function":"indexHandler"` +
				`},` +
				`"logging.googleapis.com/spanId":"000000000000004a",` +
				`"logging.googleapis.com/trace":"projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824",` +
				`"logging.googleapis.com/trace_sampled":true` +
				"}\n",
		},
	}
	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		b := callHandler(indexHandler, rr, req)

		if b.String() != test.want {
			t.Errorf("log entry:\nwant %s\ngot  %s", test.want, b.String())
		}
	}
}

// callHandler calls an HTTP handler with the provided request and returns the log output.
func callHandler(h func(w http.ResponseWriter, r *http.Request), rr http.ResponseWriter, req *http.Request) bytes.Buffer {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	originalWriter := os.Stderr
	log.SetOutput(writer)
	defer log.SetOutput(originalWriter)

	h(rr, req)
	writer.Flush()
	return buf
}
