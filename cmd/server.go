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

package cmd

import (
	"io"

	"github.com/m3dbx/pilosa/ctl"
	"github.com/m3dbx/pilosa/server"
	"github.com/m3dbx/pilosa/tracing"
	"github.com/m3dbx/pilosa/tracing/opentracing"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Server is global so that tests can control and verify it.
var Server *server.Command

// newServeCmd creates a pilosa server and runs it with command line flags.
func newServeCmd(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	Server = server.NewCommand(stdin, stdout, stderr)
	serveCmd := &cobra.Command{
		Use:   "server",
		Short: "Run Pilosa.",
		Long: `pilosa server runs Pilosa.

It will load existing data from the configured
directory and start listening for client connections
on the configured port.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Start & run the server.
			if err := Server.Start(); err != nil {
				return errors.Wrap(err, "running server")
			}

			// Initialize tracing in the command since it is global.
			var cfg jaegercfg.Configuration
			cfg.ServiceName = "pilosa"
			cfg.Sampler = &jaegercfg.SamplerConfig{
				Type:  Server.Config.Tracing.SamplerType,
				Param: Server.Config.Tracing.SamplerParam,
			}
			cfg.Reporter = &jaegercfg.ReporterConfig{
				LocalAgentHostPort: Server.Config.Tracing.AgentHostPort,
			}
			tracer, closer, err := cfg.NewTracer()
			if err != nil {
				return errors.Wrap(err, "initializing jaeger tracer")
			}
			defer closer.Close()
			tracing.GlobalTracer = opentracing.NewTracer(tracer)

			return errors.Wrap(Server.Wait(), "waiting on Server")
		},
	}

	// Attach flags to the command.
	ctl.BuildServerFlags(serveCmd, Server)
	return serveCmd
}
