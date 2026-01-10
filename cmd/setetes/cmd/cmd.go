package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sembraniteam/setetes/internal/bootstrap"
	"github.com/spf13/cobra"
)

func Start() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:     "start",
		Short:   "Start the Setetes server",
		Example: "setetes start --config ./path/to/config.yml",
		Version: "0.0.1",
		Run: func(_ *cobra.Command, _ []string) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				fmt.Printf("failed to get absolute path of Setetes: %v\n", err)
				os.Exit(1)
			}

			bts := bootstrap.New(absPath)
			if err = bts.Init(); err != nil {
				fmt.Printf("failed to initialize: %v\n", err)
				os.Exit(1)
			}
		},
	}

	c.Flags().
		StringVar(&path, "config", "", "path to the Setetes config file. Must be '.yml' or '.yaml' file.")
	err := c.MarkFlagRequired("config")
	if err != nil {
		panic(err)
	}

	return c
}

func Seed() *cobra.Command {
	var path string
	c := &cobra.Command{
		Use:     "seed",
		Short:   "Insert seed data into the database",
		Example: "setetes seed --config ./path/to/config.yml",
		Version: "0.0.1",
		Run: func(_ *cobra.Command, _ []string) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				fmt.Printf("failed to get absolute path of Setetes: %v\n", err)
				os.Exit(1)
			}

			bts := bootstrap.New(absPath)
			if err = bts.Seeder(); err != nil {
				fmt.Printf("failed to seeding: %v\n", err)
				os.Exit(1)
			}
		},
	}

	c.Flags().
		StringVar(&path, "config", "", "path to the Setetes config file. Must be '.yml' or '.yaml' file.")
	err := c.MarkFlagRequired("config")
	if err != nil {
		panic(err)
	}

	return c
}
