package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Qitmeer/llama.go/api"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/model"
)

const (
	bufferSize = 1024 * 64 // 64KB buffer for downloading
)

// downloadFile downloads a file from the given URL and saves it to the output path
// It reports progress via the progress callback function
func downloadFile(ctx context.Context, url, outputPath string, fn func(api.ProgressResponse)) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Minute, // 60 minutes timeout for large models
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	totalSize := resp.ContentLength

	// Check if file already exists and has the same size
	if info, err := os.Stat(outputPath); err == nil && info.Size() == totalSize {
		fn(api.ProgressResponse{
			Status:    "success",
			Completed: totalSize,
			Total:     totalSize,
		})
		return nil
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Report initial status
	fn(api.ProgressResponse{
		Status: "downloading",
		Total:  totalSize,
	})

	// Download with progress reporting
	buffer := make([]byte, bufferSize)
	downloaded := int64(0)
	lastReportTime := time.Now()
	reportInterval := 2 * time.Second // Report every 2 seconds

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read chunk
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				// Write to file
				if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
					return fmt.Errorf("failed to write to file: %w", writeErr)
				}
				downloaded += int64(n)

				// Report progress periodically
				if time.Since(lastReportTime) >= reportInterval {
					fn(api.ProgressResponse{
						Status:    "downloading",
						Total:     totalSize,
						Completed: downloaded,
					})
					lastReportTime = time.Now()
				}
			}

			if err == io.EOF {
				// Final progress report
				fn(api.ProgressResponse{
					Status:    "success",
					Total:     totalSize,
					Completed: downloaded,
				})
				return nil
			}

			if err != nil {
				return fmt.Errorf("error reading response: %w", err)
			}
		}
	}
}

func PullModel(ctx context.Context, hf *model.HuggingFaceModel, fn func(api.ProgressResponse)) error {
	// Report initial status
	fn(api.ProgressResponse{Status: fmt.Sprintf("pulling %s", hf.String())})

	// Resolve filename if needed (auto-detect or pattern matching)
	if hf.Filename == "" {
		fn(api.ProgressResponse{Status: "resolving filename from repository..."})
		if err := hf.ResolveFilename(); err != nil {
			return fmt.Errorf("failed to resolve filename: %w", err)
		}
		fn(api.ProgressResponse{Status: fmt.Sprintf("resolved filename: %s", hf.Filename)})
	}

	// Build download URL
	downloadURL := hf.ToDownloadURL()
	if downloadURL == "" {
		return fmt.Errorf("failed to build download URL")
	}

	fn(api.ProgressResponse{Status: fmt.Sprintf("download URL: %s", downloadURL)})

	// Determine output path
	localFilename := hf.GetLocalFilename()
	outputPath := filepath.Join(config.Conf.ModelDir, localFilename)

	fn(api.ProgressResponse{Status: fmt.Sprintf("saving to: %s", outputPath)})

	// Download the file
	if err := downloadFile(ctx, downloadURL, outputPath, fn); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	fn(api.ProgressResponse{
		Status: fmt.Sprintf("successfully downloaded to %s", outputPath),
	})

	return nil
}
