package video

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/infocus7/imp/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func getVideoDuration(inputFile string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func getVideoDimensions(inputFile string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	dimensions := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(dimensions) != 2 {
		return 0, 0, fmt.Errorf("unexpected output format from ffprobe")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func formatDuration(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// Get crop filter (using ffmpeg crop filter)
func getCropFilter(inputFile string, square bool, width, height int) (string, error) {
	if square {
		// Get video dimensions to calculate centered square crop
		width, height, err := getVideoDimensions(inputFile)
		if err != nil {
			return "", err
		}

		// Calculate the size of the square (smaller of width or height)
		size := width
		if height < width {
			size = height
		}

		// Calculate the offset for centering
		x := (width - size) / 2
		y := (height - size) / 2

		pterm.Info.Printf("Cropping to centered square: %dx%d\n", size, size)
		return fmt.Sprintf("crop=%d:%d:%d:%d", size, size, x, y), nil
	} else if width > 0 && height > 0 {
		// Custom dimensions
		pterm.Info.Printf("Cropping to custom dimensions: %dx%d\n", width, height)
		return fmt.Sprintf("crop=%d:%d", width, height), nil
	}

	pterm.Info.Println("No crop applied")
	return "", nil
}

func parseTime(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format, expected HH:MM:SS")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hours format")
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes format")
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds format")
	}

	return float64(hours*3600 + minutes*60 + seconds), nil
}

func interactiveSetup(inputFile string) (string, string, string, bool, int, int, bool, error) {
	// Show video info
	duration, err := getVideoDuration(inputFile)
	if err != nil {
		return "", "", "", false, 0, 0, false, err
	}

	width, height, err := getVideoDimensions(inputFile)
	if err != nil {
		return "", "", "", false, 0, 0, false, err
	}

	pterm.Info.Printf("Video duration: %s\n", formatDuration(duration))
	pterm.Info.Printf("Video dimensions: %dx%d\n", width, height)

	// Frame extraction mode
	options := []string{"Extract frame at specific time", "Extract random frame"}
	selectedOption, _ := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		Show("Select frame extraction mode")

	var frameTime string
	var random bool

	switch selectedOption {
	case "Extract frame at specific time":
		for {
			frameTime, _ = pterm.DefaultInteractiveTextInput.
				WithDefaultText("00:00:00").
				Show("Enter time (HH:MM:SS)")

			// Validate time format and range
			seconds, err := parseTime(frameTime)
			if err != nil {
				pterm.Error.Println(err)
				continue
			}

			if seconds < 0 {
				pterm.Error.Println("Time cannot be negative")
				continue
			}

			if seconds > duration {
				pterm.Error.Printf("Time %s is beyond video duration %s\n", frameTime, formatDuration(duration))
				continue
			}

			break
		}
	case "Extract random frame":
		random = true
	}

	// Output directory - default to input file's directory
	defaultOutDir := filepath.Dir(inputFile)
	outDir, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText(defaultOutDir).
		Show("Enter output directory (press Enter to use default)")

	if outDir == "" {
		outDir = defaultOutDir
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", "", "", false, 0, 0, false, fmt.Errorf("failed to create output directory: %v", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	defaultFilename := baseName + "_frame"

	filename, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText(defaultFilename).
		Show("Enter output filename (without extension, press Enter to use default)")

	if filename == "" {
		filename = defaultFilename
	}

	// Crop options
	cropOptions := []string{"No crop", "Square crop", "Custom dimensions"}
	selectedCrop, _ := pterm.DefaultInteractiveSelect.
		WithOptions(cropOptions).
		Show("Select crop option")

	var square bool
	var cropWidth, cropHeight int

	switch selectedCrop {
	case "Square crop":
		square = true
	case "Custom dimensions":
		for {
			widthStr, _ := pterm.DefaultInteractiveTextInput.
				WithDefaultText(strconv.Itoa(width)).
				Show("Enter crop width")

			cropWidth, err = strconv.Atoi(widthStr)
			if err != nil {
				pterm.Error.Println("Invalid width format")
				continue
			}

			if cropWidth <= 0 {
				pterm.Error.Println("Width must be positive")
				continue
			}

			if cropWidth > width {
				pterm.Error.Printf("Width %d is larger than video width %d\n", cropWidth, width)
				continue
			}

			break
		}

		for {
			heightStr, _ := pterm.DefaultInteractiveTextInput.
				WithDefaultText(strconv.Itoa(height)).
				Show("Enter crop height")

			cropHeight, err = strconv.Atoi(heightStr)
			if err != nil {
				pterm.Error.Println("Invalid height format")
				continue
			}

			if cropHeight <= 0 {
				pterm.Error.Println("Height must be positive")
				continue
			}

			if cropHeight > height {
				pterm.Error.Printf("Height %d is larger than video height %d\n", cropHeight, height)
				continue
			}

			break
		}
	}

	return frameTime, outDir, filename, random, cropWidth, cropHeight, square, nil
}

var frameCmd = &cobra.Command{
	Use:   "frame [file]",
	Short: "Extract frames from a video",
	Long:  "Extract a frame from a video file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			pterm.Error.Printf("Input file does not exist: %s\n", inputFile)
			return
		}

		// Validate that it is a (supported) video file
		if !utils.IsVideoFile(inputFile) {
			pterm.Error.Println("Input file is not a (supported) video file")
			return
		}

		// Interactive setup
		frameTime, outDir, filename, random, width, height, square, err := interactiveSetup(inputFile)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		if random {
			duration, err := getVideoDuration(inputFile)
			if err != nil {
				pterm.Error.Printf("Failed to get video duration: %v\n", err)
				return
			}

			// Pick random timestamp
			rand.Seed(time.Now().UnixNano())
			randomSeconds := rand.Float64() * duration
			frameTime = formatDuration(randomSeconds)

			pterm.Info.Printf("Selected random time: %s\n", frameTime)
		}

		// Get crop filter if needed
		cropFilter, err := getCropFilter(inputFile, square, width, height)
		if err != nil {
			pterm.Error.Printf("Failed to calculate crop dimensions: %v\n", err)
			return
		}

		if frameTime != "" {
			// Extract single frame at specific time
			outputFile := filepath.Join(outDir, filename+".jpg")
			pterm.Info.Printf("Extracting frame at %s...\n", frameTime)

			args := []string{"-i", inputFile, "-ss", frameTime, "-vframes", "1"}
			if cropFilter != "" {
				args = append(args, "-vf", cropFilter)
			}
			args = append(args, outputFile)

			cmd := exec.Command("ffmpeg", args...)
			if err := cmd.Run(); err != nil {
				pterm.Error.Printf("Failed to extract frame: %v\n", err)
				return
			}
			pterm.Success.Printf("Frame extracted successfully: %s\n", outputFile)
		} else {
			pterm.Error.Println("No frame time specified")
		}
	},
}
