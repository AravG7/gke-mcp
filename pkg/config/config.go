// Copyright 2025 Google LLC
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

// Package config loads configuration derived from local gcloud defaults.
package config

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// Config contains runtime configuration derived from the environment.
type Config struct {
	userAgent        string
	defaultProjectID string
	defaultLocation  string
	logger           *slog.Logger
}

// Logger returns the logger instance.
func (c *Config) Logger() *slog.Logger {
	return c.logger
}

// UserAgent returns the user agent string for outbound API calls.
func (c *Config) UserAgent() string {
	return c.userAgent
}

// DefaultProjectID returns the default GCP project ID, if set.
func (c *Config) DefaultProjectID() string {
	return c.defaultProjectID
}

// DefaultLocation returns the default GCP region or zone, if set.
func (c *Config) DefaultLocation() string {
	return c.defaultLocation
}

// New constructs a Config populated from gcloud and build version.
func New(version string) *Config {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	return &Config{
		userAgent:        "gke-mcp/" + version,
		defaultProjectID: getDefaultProjectID(logger),
		defaultLocation:  getDefaultLocation(logger),
		logger:           logger,
	}
}

func getDefaultProjectID(logger *slog.Logger) string {
	projectID, err := getGcloudConfig("core/project")
	if err != nil {
		logger.Error("Failed to get default project", "error", err)
		return ""
	}
	return projectID
}

func getDefaultLocation(logger *slog.Logger) string {
	region, err := getGcloudConfig("compute/region")
	if err == nil {
		return region
	}
	zone, err := getGcloudConfig("compute/zone")
	if err == nil {
		return zone
	}
	return ""
}

func getGcloudConfig(key string) (string, error) {
	// #nosec G204
	out, err := exec.Command("gcloud", "config", "get", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
