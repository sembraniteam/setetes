package main

import (
	"embed"

	"github.com/sembraniteam/setetes/cmd/setetes/cmd"
	"github.com/spf13/cobra"
)

//go:embed ASCII
var ascii embed.FS

func main() {
	c := &cobra.Command{
		Use:     "setetes",
		Short:   "CLI tool for Setetes blood donation management system.",
		Version: "0.0.1",
	}

	data, _ := ascii.ReadFile("ASCII")
	println(string(data))

	c.CompletionOptions.DisableDefaultCmd = true
	c.AddCommand(cmd.Start(), cmd.Seed())
	cobra.CheckErr(c.Execute())
}
