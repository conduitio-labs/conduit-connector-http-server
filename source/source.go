// Copyright © 2022 Meroxa, Inc.
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

package source

import (
	"context"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/conduitio-labs/conduit-connector-http-server/config"
	"github.com/conduitio-labs/conduit-connector-http-server/source/server"
)

// Source connector.
type Source struct {
	sdk.UnimplementedSource

	config config.Config
	server HTTPServer
}

// New initialises a new source.
func New() sdk.Source {
	return &Source{}
}

// Configure parses and stores configurations, returns an error in case of invalid configuration.
func (s *Source) Configure(ctx context.Context, cfgRaw map[string]string) error {
	cfg, err := config.Parse(cfgRaw)
	if err != nil {
		return err
	}

	s.config = cfg

	return nil
}

// Open prepare the plugin to start server.
func (s *Source) Open(ctx context.Context, rp sdk.Position) error {
	httpServer := server.New(s.config.Port)
	s.server = httpServer

	s.server.Start()

	return nil
}

// Read gets record.
func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	return s.server.GetRecord(ctx)
}

// Teardown gracefully shutdown connector.
func (s *Source) Teardown(ctx context.Context) error {
	return s.server.Stop(ctx)
}

// Ack check if record with position was recorded.
func (s *Source) Ack(ctx context.Context, p sdk.Position) error {
	return nil
}
