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

package config

import (
	"strconv"
)

const (
	KeyPort string = "port"

	defaultPort = 3000
)

// Config represents configuration needed for HTTP Source connector.
type Config struct {
	// Port - HTTP server will run in this port.
	Port int
}

// Parse attempts to parse plugins.Config into a Config struct.
func Parse(cfg map[string]string) (Config, error) {
	port := defaultPort

	if cfg[KeyPort] != "" {
		result, err := strconv.Atoi(cfg[KeyPort])
		if err != nil {
			return Config{}, ErrFieldInvalidPortType
		}

		port = result
	}

	return Config{
		Port: port,
	}, nil
}
