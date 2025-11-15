package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	DefaultHFHost      = "huggingface.co"
	DefaultHFBranch    = "main"
	DefaultHFNamespace = "llamago"
)

// HuggingFaceModel represents a parsed Hugging Face model reference
// It supports multiple input formats:
// 1. Full URL: https://huggingface.co/namespace/repo/resolve/main/file.gguf
// 2. Repo with file: namespace/repo:file.gguf
// 3. Repo with pattern: namespace/repo:Q4_K_M (will search for matching files)
// 4. Simple repo: namespace/repo (will auto-detect GGUF files)
// 5. Repo only: repo (uses default namespace "llamago")
type HuggingFaceModel struct {
	// Host is the Hugging Face host (default: huggingface.co)
	Host string

	// Namespace is the user or organization name (default: llamago)
	Namespace string

	// Repo is the repository name
	Repo string

	// Branch is the git branch or tag (default: main)
	Branch string

	// Filename is the specific file to download
	// If empty, will be auto-detected
	Filename string

	// Pattern is used to filter files when Filename is not specified
	// e.g., "Q4_K_M", "Q8_0", etc.
	Pattern string
}

// ParseHuggingFaceModel parses a Hugging Face model reference string
// Supported formats:
//   - https://huggingface.co/namespace/repo/resolve/main/file.gguf
//   - namespace/repo:file.gguf
//   - namespace/repo:Q4_K_M
//   - namespace/repo
//   - repo (uses default namespace: llamago)
func ParseHuggingFaceModel(s string) (*HuggingFaceModel, error) {
	hf := &HuggingFaceModel{
		Host:      DefaultHFHost,
		Branch:    DefaultHFBranch,
		Namespace: DefaultHFNamespace,
	}

	// Check if it's a full URL
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return parseHuggingFaceURL(s)
	}

	// Parse short format: namespace/repo:file or namespace/repo:pattern
	var repoPath string
	var suffix string

	if idx := strings.Index(s, ":"); idx >= 0 {
		repoPath = s[:idx]
		suffix = s[idx+1:]
	} else {
		repoPath = s
	}

	// Split namespace/repo
	parts := strings.Split(repoPath, "/")
	if len(parts) == 1 {
		// Only repo name provided, use default namespace
		hf.Repo = parts[0]
	} else if len(parts) >= 2 {
		// Both namespace and repo provided
		hf.Namespace = parts[0]
		hf.Repo = parts[1]
	} else {
		return nil, fmt.Errorf("invalid format: got '%s'", repoPath)
	}

	// Determine if suffix is a filename or pattern
	if suffix != "" {
		if strings.HasSuffix(suffix, ".gguf") {
			hf.Filename = suffix
		} else {
			hf.Pattern = suffix
		}
	}

	return hf, nil
}

// parseHuggingFaceURL parses a full Hugging Face URL
func parseHuggingFaceURL(s string) (*HuggingFaceModel, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	hf := &HuggingFaceModel{
		Host:   u.Host,
		Branch: DefaultHFBranch,
	}

	// Parse path: /namespace/repo/resolve/branch/path/to/file.gguf
	path := strings.Trim(u.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid URL path: expected at least namespace/repo")
	}

	hf.Namespace = parts[0]
	hf.Repo = parts[1]

	// Check for resolve/branch/filename pattern
	if len(parts) >= 4 && parts[2] == "resolve" {
		hf.Branch = parts[3]
		if len(parts) > 4 {
			hf.Filename = strings.Join(parts[4:], "/")
		}
	} else if len(parts) > 2 {
		// Direct file reference: namespace/repo/file.gguf
		hf.Filename = strings.Join(parts[2:], "/")
	}

	return hf, nil
}

// ToDownloadURL returns the full download URL for this model
func (hf *HuggingFaceModel) ToDownloadURL() string {
	if hf.Filename == "" {
		return ""
	}

	return fmt.Sprintf("https://%s/%s/%s/resolve/%s/%s",
		hf.Host,
		hf.Namespace,
		hf.Repo,
		hf.Branch,
		hf.Filename,
	)
}

// ToRepoURL returns the URL to the repository
func (hf *HuggingFaceModel) ToRepoURL() string {
	return fmt.Sprintf("https://%s/%s/%s",
		hf.Host,
		hf.Namespace,
		hf.Repo,
	)
}

// ToAPIURL returns the Hugging Face API URL for listing files
func (hf *HuggingFaceModel) ToAPIURL() string {
	return fmt.Sprintf("https://%s/api/models/%s/%s/tree/%s",
		hf.Host,
		hf.Namespace,
		hf.Repo,
		hf.Branch,
	)
}

// String returns a string representation of the model
func (hf *HuggingFaceModel) String() string {
	if hf.Filename != "" {
		return fmt.Sprintf("%s/%s:%s", hf.Namespace, hf.Repo, hf.Filename)
	}
	if hf.Pattern != "" {
		return fmt.Sprintf("%s/%s:%s", hf.Namespace, hf.Repo, hf.Pattern)
	}
	return fmt.Sprintf("%s/%s", hf.Namespace, hf.Repo)
}

// HFFileInfo represents a file in a Hugging Face repository
type HFFileInfo struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Size int64  `json:"size,omitempty"`
}

// ListGGUFFiles fetches the list of GGUF files from the repository
func (hf *HuggingFaceModel) ListGGUFFiles() ([]HFFileInfo, error) {
	apiURL := hf.ToAPIURL()

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed (%s): %s", resp.Status, string(body))
	}

	var files []HFFileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Filter for GGUF files
	var ggufFiles []HFFileInfo
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(strings.ToLower(file.Path), ".gguf") {
			ggufFiles = append(ggufFiles, file)
		}
	}

	return ggufFiles, nil
}

// ResolveFilename attempts to determine the best matching GGUF file
// based on the pattern or automatically selects one if no pattern is given
func (hf *HuggingFaceModel) ResolveFilename() error {
	if hf.Filename != "" {
		// Filename already specified
		return nil
	}

	files, err := hf.ListGGUFFiles()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no GGUF files found in repository")
	}

	// If pattern is specified, filter by pattern
	if hf.Pattern != "" {
		var matched []HFFileInfo
		pattern := strings.ToLower(hf.Pattern)
		for _, file := range files {
			filename := strings.ToLower(filepath.Base(file.Path))
			if strings.Contains(filename, pattern) {
				matched = append(matched, file)
			}
		}

		if len(matched) == 0 {
			return fmt.Errorf("no files matching pattern '%s' found", hf.Pattern)
		}

		// Use the first match
		hf.Filename = matched[0].Path
		return nil
	}

	// No pattern specified, use the first GGUF file
	hf.Filename = files[0].Path
	return nil
}

// IsValid checks if the model reference is valid
func (hf *HuggingFaceModel) IsValid() bool {
	return hf.Namespace != "" && hf.Repo != "" && hf.Host != ""
}

// GetLocalFilename returns the filename to use when saving locally
func (hf *HuggingFaceModel) GetLocalFilename() string {
	if hf.Filename != "" {
		return filepath.Base(hf.Filename)
	}
	return ""
}