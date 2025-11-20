package base

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/megalodev/setetes/internal/bootstrap"
	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start the Setetes server",
		Example: "setetes start --config ./path/to/config.yml",
		Version: "0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				fmt.Printf("failed to get absolute path of setetes: %v\n", err)
				os.Exit(1)
			}

			boots := bootstrap.New(absPath)
			if err := boots.Init(); err != nil {
				fmt.Printf("failed to initialize: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&path, "config", "", "path to the Setetes config file")
	if err := cmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}

	return cmd
}
