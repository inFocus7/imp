package cmd

import (
	"github.com/infocus7/imp/cmd/convert"
	"github.com/infocus7/imp/cmd/version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "imp",
		Short: "imp is a tool for media processing",
		Long:  "imp (inf0's media pipeline) is a tool for media processing",
	}

	cmd.AddCommand(convert.ConvertCmd)

	return cmd
}

func Execute() error {
	cmd := Command()

	cmd.Version = "n/a"
	cmd.SetVersionTemplate(version.VersionTemplate())

	return cmd.Execute()
}
