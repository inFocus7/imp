package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/infocus7/imp/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func interactiveAudioSetup(inputFile string) (string, string, string, error) {
	// Get base name for default filename
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	defaultOutDir := filepath.Dir(inputFile)

	// Get output directory
	outDir, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText(defaultOutDir).
		Show("Enter output directory (press Enter to use default)")

	// If user just pressed Enter, use default
	if outDir == "" {
		outDir = defaultOutDir
	}

	// Ensure the directory exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", "", "", fmt.Errorf("failed to create output directory: %v", err)
	}

	// Get output filename
	filename, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText(baseName).
		Show("Enter output filename (without extension, press Enter to use default)")

	// If user just pressed Enter, use default
	if filename == "" {
		filename = baseName
	}

	// Get output format, excluding the input file's format
	inputExt := strings.ToLower(filepath.Ext(inputFile))
	formatOptions := make([]string, 0, len(utils.SupportedAudioExtensions))
	for ext := range utils.SupportedAudioExtensions {
		if ext != inputExt {
			formatOptions = append(formatOptions, strings.TrimPrefix(ext, "."))
		}
	}

	toFmt, _ := pterm.DefaultInteractiveSelect.
		WithOptions(formatOptions).
		Show("Select output format")

	return toFmt, outDir, filename, nil
}

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

		// Interactive setup
		toFmt, outDir, filename, err := interactiveAudioSetup(inputFile)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		pterm.Info.Printf("Converting %s to %s format...\n", inputFile, toFmt)
		pterm.Info.Printf("Output will be saved in: %s\n", outDir)
		pterm.Info.Printf("Output file will be: %s\n", filepath.Join(outDir, filename+"."+toFmt))

		// Convert the file using ffmpeg
		convertedFile := filepath.Join(outDir, filename+"."+toFmt)
		convertCmd := exec.Command("ffmpeg", "-i", inputFile, "-f", toFmt, convertedFile)
		if err := convertCmd.Run(); err != nil {
			pterm.Error.Printf("Failed to convert file: %v\n", err)
			return
		}

		pterm.Success.Printf("File converted successfully: %s\n", convertedFile)
	},
}
