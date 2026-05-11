// Copyright 2026 Google LLC
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

package apps

import (
	"context"
	"testing"

	"github.com/GoogleCloudPlatform/gke-mcp/pkg/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestInstallApps(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		s       *mcp.Server
		c       *config.Config
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			s: mcp.NewServer(
				&mcp.Implementation{
					Name:    "Test Server",
					Version: "1.0.0",
				},
				&mcp.ServerOptions{},
			),
			c:       &config.Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InstallApps(tt.ctx, tt.s, tt.c); (err != nil) != tt.wantErr {
				t.Errorf("InstallApps() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
