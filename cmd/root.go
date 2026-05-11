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

// Package cmd wires up the CLI entrypoints for the GKE MCP server.
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/GoogleCloudPlatform/gke-mcp/pkg/install"
	"github.com/GoogleCloudPlatform/gke-mcp/pkg/server"
	"github.com/spf13/cobra"
)

var (
	version = "(unknown)"

	// command flags
	serverMode     string
	serverHost     string
	serverPort     int
	allowedOrigins []string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gke-mcp",
		Short: "An MCP Server for Google Kubernetes Engine",
		Run:   runRootCmd,
	}

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install the GKE MCP Server into your AI tool settings.",
	}

	installGeminiCLICmd = &cobra.Command{
		Use:   "gemini-cli",
		Short: "Install the GKE MCP Server into your Gemini CLI settings.",
		Run:   runInstallGeminiCLICmd,
	}

	installCursorCmd = &cobra.Command{
		Use:   "cursor",
		Short: "Install the GKE MCP Server into your Cursor settings.",
		Run:   runInstallCursorCmd,
	}

	installClaudeDesktopCmd = &cobra.Command{
		Use:   "claude-desktop",
		Short: "Install the GKE MCP Server into your Claude Desktop settings.",
		Run:   runInstallClaudeDesktopCmd,
	}

	installClaudeCodeCmd = &cobra.Command{
		Use:   "claude-code",
		Short: "Install the GKE MCP Server into your Claude Code CLI settings.",
		Run:   runInstallClaudeCodeCmd,
	}

	installDeveloper   bool
	installProjectOnly bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Only attempt to read build info if version hasn't been set via ldflags
	if version == "(unknown)" || version == "dev" || version == "(devel)" {
		if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			version = bi.Main.Version
		}
	}

	rootCmd.Flags().StringVar(&serverMode, "server-mode", "stdio", "transport to use for the server: stdio (default) or http")
	rootCmd.Flags().StringVar(&serverHost, "server-host", "127.0.0.1", "server host to use when server-mode is http; defaults to 127.0.0.1")
	rootCmd.Flags().IntVar(&serverPort, "server-port", 8080, "server port to use when server-mode is http; defaults to 8080")
	rootCmd.Flags().StringSliceVar(&allowedOrigins, "allowed-origins", []string{"http://localhost"}, "comma-separated list of allowed Origin headers")
	rootCmd.AddCommand(installCmd)

	installCmd.AddCommand(installGeminiCLICmd)
	installCmd.AddCommand(installCursorCmd)
	installCmd.AddCommand(installClaudeDesktopCmd)
	installCmd.AddCommand(installClaudeCodeCmd)

	installGeminiCLICmd.Flags().BoolVarP(&installDeveloper, "developer", "d", false, "Install the MCP Server in developer mode for Gemini CLI")
	installGeminiCLICmd.Flags().BoolVarP(&installProjectOnly, "project-only", "p", false, "Install the MCP Server only for the current project. Please run this in the root directory of your project")

	installCursorCmd.Flags().BoolVarP(&installProjectOnly, "project-only", "p", false, "Install the MCP Server only for the current project. Please run this in the root directory of your project")
	installClaudeCodeCmd.Flags().BoolVarP(&installProjectOnly, "project-only", "p", false, "Install the MCP Server only for the current project. Please run this in the root directory of your project")
}

func runRootCmd(cmd *cobra.Command, _ []string) {
	opts := server.Options{
		ServerMode:     serverMode,
		ServerHost:     serverHost,
		ServerPort:     serverPort,
		AllowedOrigins: allowedOrigins,
	}
	if err := server.StartMCPServer(cmd.Context(), opts, version); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func installOptions() (*install.Options, error) {
	return install.NewInstallOptions(
		version,
		installProjectOnly,
		installDeveloper,
	)
}

func runInstallGeminiCLICmd(_ *cobra.Command, _ []string) {
	opts, err := installOptions()
	if err != nil {
		log.Fatalf("Failed to get install options: %v", err)
	}

	if err := install.GeminiCLIExtension(opts); err != nil {
		log.Fatalf("Failed to install for gemini-cli: %v", err)
	}
	fmt.Println("Successfully installed GKE MCP server as a gemini-cli extension.")
}

func runInstallCursorCmd(_ *cobra.Command, _ []string) {
	opts, err := installOptions()
	if err != nil {
		log.Fatalf("Failed to get install options: %v", err)
	}

	if err := install.CursorMCPExtension(opts); err != nil {
		log.Fatalf("Failed to install for cursor: %v", err)
	}
	fmt.Println("Successfully installed GKE MCP server as a cursor MCP server.")
}

func runInstallClaudeDesktopCmd(_ *cobra.Command, _ []string) {
	opts, err := installOptions()
	if err != nil {
		log.Fatalf("Failed to get install options: %v", err)
	}

	if err := install.ClaudeDesktopExtension(opts); err != nil {
		log.Fatalf("Failed to install for Claude Desktop: %v", err)
	}
	fmt.Println("Successfully installed GKE MCP server in Claude Desktop configuration.")
}

func runInstallClaudeCodeCmd(_ *cobra.Command, _ []string) {
	opts, err := installOptions()
	if err != nil {
		log.Fatalf("Failed to get install options: %v", err)
	}

	if err := install.ClaudeCodeExtension(opts); err != nil {
		log.Fatalf("Failed to install for Claude Code: %v", err)
	}

	fmt.Println("Successfully installed GKE MCP server for Claude Code.")
}
