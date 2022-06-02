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
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		cfg         map[string]string
		want        Config
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid config",
			cfg: map[string]string{
				KeyPort: "8888",
			},
			want: Config{
				Port: 8888,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "config with default port",
			cfg:  map[string]string{},
			want: Config{
				Port: defaultPort,
			},
			wantErr:     false,
			expectedErr: "",
		},
		{
			name: "wrong port type",
			cfg: map[string]string{
				KeyPort: "test",
			},
			want:        Config{},
			wantErr:     true,
			expectedErr: ErrFieldInvalidPortType.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.cfg)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parse error = \"%s\", wantErr %t", err.Error(), tt.wantErr)

					return
				}

				if err.Error() != tt.expectedErr {
					t.Errorf("expected error \"%s\", got \"%s\"", tt.expectedErr, err.Error())

					return
				}

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse = %v, want %v", got, tt.want)
			}
		})
	}
}
