package convert

import (
	"os"

	"github.com/infocus7/imp/cmd/convert/audio"
	"github.com/infocus7/imp/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert media files to different formats",
	Long:  pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("Convert media files to different formats"),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pterm.Error.Println("Please specify a conversion type.")
			return
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := utils.CheckDependencies("convert"); err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	ConvertCmd.AddCommand(audio.AudioCmd)
}
