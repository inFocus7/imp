package audio

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/infocus7/imp/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	toFmt   string
	outDir  string
	outFile string
)

var AudioCmd = &cobra.Command{
	Use:   "audio [file]",
	Short: "Convert audio files to different formats",
	Long:  "Convert audio files to different formats (e.g., WAV to MP3)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			pterm.Error.Printf("Input file does not exist: %s\n", inputFile)
			return
		}

		// Validate that it is a (supported) audio file
		if !utils.IsAudioFile(inputFile) {
			pterm.Error.Println("Input file is not a (supported) audio file")
			return
		}

		// Validate output format
		if toFmt == "" {
			pterm.Error.Println("Please specify an output format using --to")
			return
		}

		if outDir == "" {
			outDir = filepath.Dir(inputFile)
		}

		if outFile == "" {
			outFile = filepath.Base(inputFile)
			outFile = strings.TrimSuffix(outFile, filepath.Ext(outFile))
		}

		if err := os.MkdirAll(outDir, 0755); err != nil {
			pterm.Error.Printf("Failed to create output directory: %v\n", err)
			return
		}

		pterm.Info.Printf("Converting %s to %s format...\n", inputFile, toFmt)
		pterm.Info.Printf("Output will be saved in: %s\n", outDir)
		pterm.Info.Printf("Output file will be: %s\n", filepath.Join(outDir, outFile+"."+toFmt))

		// Convert the file (simple ffmpeg command)
		convertedFile := filepath.Join(outDir, outFile+"."+toFmt)
		convertCmd := exec.Command("ffmpeg", "-i", inputFile, "-f", toFmt, convertedFile)
		if err := convertCmd.Run(); err != nil {
			pterm.Error.Printf("Failed to convert file: %v\n", err)
			return
		}

		pterm.Success.Printf("File converted successfully: %s\n", convertedFile)
	},
}

func init() {
	AudioCmd.Flags().StringVarP(&toFmt, "to", "t", "", "The format to convert to (e.g., mp3)")
	AudioCmd.Flags().StringVarP(&outDir, "dir", "d", "", "The output directory (defaults to input file directory)")
	AudioCmd.Flags().StringVarP(&outFile, "file", "f", "", "The output file name (defaults to input file name, exclude extension)")
}
