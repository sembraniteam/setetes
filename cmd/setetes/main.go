package main

import (
	"embed"

	"github.com/sembraniteam/setetes/cmd/setetes/base"
	"github.com/spf13/cobra"
)

//go:embed ASCII
var ascii embed.FS

func main() {
	cmd := &cobra.Command{
		Use:     "setetes",
		Short:   "CLI tool for Setetes blood donation management system.",
		Version: "0.0.1",
	}

	data, _ := ascii.ReadFile("ASCII")
	println(string(data))

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.AddCommand(base.StartCmd())
	cobra.CheckErr(cmd.Execute())
}
