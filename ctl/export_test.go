// Copyright 2017 Pilosa Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctl

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/m3dbx/pilosa"
	"github.com/m3dbx/pilosa/test"
)

func TestExportCommand_Validation(t *testing.T) {
	buf := bytes.Buffer{}
	stdin, stdout, stderr := GetIO(buf)

	cm := NewExportCommand(stdin, stdout, stderr)

	err := cm.Run(context.Background())
	if err != pilosa.ErrIndexRequired {
		t.Fatalf("Command not working, expect: %s, actual: '%s'", pilosa.ErrIndexRequired, err)
	}

	cm.Index = "i"
	err = cm.Run(context.Background())
	if err != pilosa.ErrFieldRequired {
		t.Fatalf("Command not working, expect: %s, actual: '%s'", pilosa.ErrFieldRequired, err)
	}
}

func TestExportCommand_Run(t *testing.T) {
	cmd := test.MustRunCluster(t, 1)[0]

	buf := bytes.Buffer{}
	stdin, stdout, stderr := GetIO(buf)
	cm := NewExportCommand(stdin, stdout, stderr)
	hostport := cmd.API.Node().URI.HostPort()
	cm.Host = hostport

	http.DefaultClient.Do(test.MustNewHTTPRequest("POST", "http://"+hostport+"/index/i", strings.NewReader("")))
	http.DefaultClient.Do(test.MustNewHTTPRequest("POST", "http://"+hostport+"/index/i/field/f", strings.NewReader("")))

	cm.Index = "i"
	cm.Field = "f"
	if err := cm.Run(context.Background()); err != nil {
		t.Fatalf("Export Run doesn't work: %s", err)
	}
}
