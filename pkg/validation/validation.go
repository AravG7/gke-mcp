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

// Package validation provides input validation utilities for MCP tools.
package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// nodeNameRegex matches a standard Kubernetes node name.
	nodeNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9\-\.]*[a-z0-9]$`)
	// projectIDRegex matches a standard Google Cloud project ID.
	projectIDRegex = regexp.MustCompile(`^[a-z][a-z0-9\-]{4,28}[a-z0-9]$`)
)

// ValidateNodeName ensures the node name is valid and safe for use in commands.
func ValidateNodeName(name string) error {
	if !nodeNameRegex.MatchString(name) {
		return fmt.Errorf("invalid node name: %s", name)
	}
	return nil
}

// ValidateProjectID ensures the project ID is valid.
func ValidateProjectID(id string) error {
	if id == "" {
		return nil // Allow empty if it's going to be defaulted
	}
	if !projectIDRegex.MatchString(id) {
		return fmt.Errorf("invalid project ID: %s", id)
	}
	return nil
}

// ValidatePath ensures a path is safe (no parent directory traversal).
func ValidatePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path cannot contain '..' for security reasons: %s", path)
	}
	return nil
}
