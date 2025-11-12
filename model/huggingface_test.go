package model

import (
	"testing"
)

func TestParseHuggingFaceModel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNS    string
		wantRepo  string
		wantFile  string
		wantPat   string
		wantBranch string
		wantErr   bool
	}{
		{
			name:      "Simple repo",
			input:     "unsloth/llama-3-8b",
			wantNS:    "unsloth",
			wantRepo:  "llama-3-8b",
			wantFile:  "",
			wantPat:   "",
			wantBranch: "main",
			wantErr:   false,
		},
		{
			name:      "Repo with file",
			input:     "unsloth/llama-3-8b:llama-3-8b-Q4_K_M.gguf",
			wantNS:    "unsloth",
			wantRepo:  "llama-3-8b",
			wantFile:  "llama-3-8b-Q4_K_M.gguf",
			wantPat:   "",
			wantBranch: "main",
			wantErr:   false,
		},
		{
			name:      "Repo with pattern",
			input:     "unsloth/llama-3-8b:Q4_K_M",
			wantNS:    "unsloth",
			wantRepo:  "llama-3-8b",
			wantFile:  "",
			wantPat:   "Q4_K_M",
			wantBranch: "main",
			wantErr:   false,
		},
		{
			name:      "Full URL with resolve",
			input:     "https://huggingface.co/unsloth/llama-3-8b/resolve/main/llama-3-8b-Q4_K_M.gguf",
			wantNS:    "unsloth",
			wantRepo:  "llama-3-8b",
			wantFile:  "llama-3-8b-Q4_K_M.gguf",
			wantPat:   "",
			wantBranch: "main",
			wantErr:   false,
		},
		{
			name:      "Full URL with different branch",
			input:     "https://huggingface.co/microsoft/phi-2/resolve/v1.0/phi-2.gguf",
			wantNS:    "microsoft",
			wantRepo:  "phi-2",
			wantFile:  "phi-2.gguf",
			wantPat:   "",
			wantBranch: "v1.0",
			wantErr:   false,
		},
		{
			name:      "Invalid format - no slash",
			input:     "invalidmodel",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hf, err := ParseHuggingFaceModel(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHuggingFaceModel() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseHuggingFaceModel() unexpected error: %v", err)
				return
			}

			if hf.Namespace != tt.wantNS {
				t.Errorf("Namespace = %v, want %v", hf.Namespace, tt.wantNS)
			}
			if hf.Repo != tt.wantRepo {
				t.Errorf("Repo = %v, want %v", hf.Repo, tt.wantRepo)
			}
			if hf.Filename != tt.wantFile {
				t.Errorf("Filename = %v, want %v", hf.Filename, tt.wantFile)
			}
			if hf.Pattern != tt.wantPat {
				t.Errorf("Pattern = %v, want %v", hf.Pattern, tt.wantPat)
			}
			if hf.Branch != tt.wantBranch {
				t.Errorf("Branch = %v, want %v", hf.Branch, tt.wantBranch)
			}
		})
	}
}

func TestHuggingFaceModel_ToDownloadURL(t *testing.T) {
	tests := []struct {
		name string
		hf   *HuggingFaceModel
		want string
	}{
		{
			name: "Standard model",
			hf: &HuggingFaceModel{
				Host:      "huggingface.co",
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
				Branch:    "main",
				Filename:  "llama-3-8b-Q4_K_M.gguf",
			},
			want: "https://huggingface.co/unsloth/llama-3-8b/resolve/main/llama-3-8b-Q4_K_M.gguf",
		},
		{
			name: "File in subdirectory",
			hf: &HuggingFaceModel{
				Host:      "huggingface.co",
				Namespace: "microsoft",
				Repo:      "phi-2",
				Branch:    "main",
				Filename:  "gguf/phi-2-Q4_K_M.gguf",
			},
			want: "https://huggingface.co/microsoft/phi-2/resolve/main/gguf/phi-2-Q4_K_M.gguf",
		},
		{
			name: "No filename",
			hf: &HuggingFaceModel{
				Host:      "huggingface.co",
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
				Branch:    "main",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hf.ToDownloadURL()
			if got != tt.want {
				t.Errorf("ToDownloadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuggingFaceModel_String(t *testing.T) {
	tests := []struct {
		name string
		hf   *HuggingFaceModel
		want string
	}{
		{
			name: "With filename",
			hf: &HuggingFaceModel{
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
				Filename:  "llama-3-8b-Q4_K_M.gguf",
			},
			want: "unsloth/llama-3-8b:llama-3-8b-Q4_K_M.gguf",
		},
		{
			name: "With pattern",
			hf: &HuggingFaceModel{
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
				Pattern:   "Q4_K_M",
			},
			want: "unsloth/llama-3-8b:Q4_K_M",
		},
		{
			name: "Simple repo",
			hf: &HuggingFaceModel{
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
			},
			want: "unsloth/llama-3-8b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hf.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuggingFaceModel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		hf   *HuggingFaceModel
		want bool
	}{
		{
			name: "Valid model",
			hf: &HuggingFaceModel{
				Host:      "huggingface.co",
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
			},
			want: true,
		},
		{
			name: "Missing namespace",
			hf: &HuggingFaceModel{
				Host: "huggingface.co",
				Repo: "llama-3-8b",
			},
			want: false,
		},
		{
			name: "Missing repo",
			hf: &HuggingFaceModel{
				Host:      "huggingface.co",
				Namespace: "unsloth",
			},
			want: false,
		},
		{
			name: "Missing host",
			hf: &HuggingFaceModel{
				Namespace: "unsloth",
				Repo:      "llama-3-8b",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hf.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHuggingFaceModel_GetLocalFilename(t *testing.T) {
	tests := []struct {
		name string
		hf   *HuggingFaceModel
		want string
	}{
		{
			name: "Simple filename",
			hf: &HuggingFaceModel{
				Filename: "model.gguf",
			},
			want: "model.gguf",
		},
		{
			name: "Filename with path",
			hf: &HuggingFaceModel{
				Filename: "gguf/subfolder/model-Q4_K_M.gguf",
			},
			want: "model-Q4_K_M.gguf",
		},
		{
			name: "No filename",
			hf:   &HuggingFaceModel{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hf.GetLocalFilename()
			if got != tt.want {
				t.Errorf("GetLocalFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
