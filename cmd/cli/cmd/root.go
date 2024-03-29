/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/shoet/blog/internal/logging"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Blog application management CLI",
	Long:  "Blog application management CLI",
}

func Execute() {
	ctx := context.Background()
	logger := logging.NewLogger(os.Stdout, "debug")
	ctx = context.WithValue(ctx, logging.LoggerKey{}, logger)
	rootCmd.SetContext(ctx)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
