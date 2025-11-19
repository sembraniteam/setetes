package main

import (
	"github.com/megalodev/setetes/cmd/setetes/base"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:     "setetes",
		Short:   "CLI tool for Setetes blood donation management system",
		Version: "0.0.1",
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(base.StartCmd())
	cobra.CheckErr(cmd.Execute())
}
