# HuggingFaceModel 使用示例

`HuggingFaceModel` 是一个专门用于解析和处理 Hugging Face 模型引用的工具类，支持多种输入格式并能自动查找和下载 GGUF 格式的模型文件。

## 功能特性

- ✅ 支持多种输入格式（完整 URL、简化格式、模式匹配）
- ✅ 自动解析 Hugging Face 仓库结构
- ✅ 支持文件名模式匹配（如 Q4_K_M、Q8_0 等）
- ✅ 自动查找仓库中的 GGUF 文件
- ✅ 完整的单元测试覆盖

## 支持的输入格式

### 1. 简化格式（推荐）

```go
// 指定完整文件名
hf, _ := model.ParseHuggingFaceModel("unsloth/llama-3-8b:llama-3-8b-Q4_K_M.gguf")

// 使用模式匹配（自动查找包含 "Q4_K_M" 的文件）
hf, _ := model.ParseHuggingFaceModel("unsloth/llama-3-8b:Q4_K_M")

// 自动检测（使用仓库中的第一个 GGUF 文件）
hf, _ := model.ParseHuggingFaceModel("unsloth/llama-3-8b")
```

### 2. 完整 URL 格式

```go
// 标准 resolve 格式
hf, _ := model.ParseHuggingFaceModel(
    "https://huggingface.co/unsloth/llama-3-8b/resolve/main/llama-3-8b-Q4_K_M.gguf",
)

// 指定分支版本
hf, _ := model.ParseHuggingFaceModel(
    "https://huggingface.co/microsoft/phi-2/resolve/v1.0/phi-2.gguf",
)

// 子目录中的文件
hf, _ := model.ParseHuggingFaceModel(
    "https://huggingface.co/microsoft/phi-2/resolve/main/gguf/phi-2-Q4_K_M.gguf",
)
```

## 完整使用示例

### 示例 1: 基本用法

```go
package main

import (
    "fmt"
    "github.com/Qitmeer/llama.go/model"
)

func main() {
    // 解析模型引用
    hf, err := model.ParseHuggingFaceModel("unsloth/llama-3-8b:Q4_K_M")
    if err != nil {
        panic(err)
    }

    // 检查有效性
    if !hf.IsValid() {
        panic("invalid model reference")
    }

    // 打印模型信息
    fmt.Printf("Namespace: %s\n", hf.Namespace)  // unsloth
    fmt.Printf("Repo: %s\n", hf.Repo)            // llama-3-8b
    fmt.Printf("Pattern: %s\n", hf.Pattern)      // Q4_K_M
    fmt.Printf("Branch: %s\n", hf.Branch)        // main
}
```

### 示例 2: 自动解析文件名

```go
package main

import (
    "fmt"
    "github.com/Qitmeer/llama.go/model"
)

func main() {
    // 使用模式匹配
    hf, err := model.ParseHuggingFaceModel("unsloth/llama-3-8b:Q4_K_M")
    if err != nil {
        panic(err)
    }

    // 从仓库中自动查找匹配的文件
    if err := hf.ResolveFilename(); err != nil {
        panic(err)
    }

    // 获取下载 URL
    downloadURL := hf.ToDownloadURL()
    fmt.Printf("Download URL: %s\n", downloadURL)

    // 获取本地保存的文件名
    localFilename := hf.GetLocalFilename()
    fmt.Printf("Local filename: %s\n", localFilename)
}
```

### 示例 3: 列出仓库中的所有 GGUF 文件

```go
package main

import (
    "fmt"
    "github.com/Qitmeer/llama.go/model"
)

func main() {
    hf, err := model.ParseHuggingFaceModel("unsloth/llama-3-8b")
    if err != nil {
        panic(err)
    }

    // 获取所有 GGUF 文件
    files, err := hf.ListGGUFFiles()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d GGUF files:\n", len(files))
    for _, file := range files {
        fmt.Printf("  - %s (%.2f MB)\n", file.Path, float64(file.Size)/(1024*1024))
    }
}
```

### 示例 4: 使用 PullModel 下载模型

```go
package main

import (
    "context"
    "fmt"
    "github.com/Qitmeer/llama.go/api"
    "github.com/Qitmeer/llama.go/model"
    "github.com/Qitmeer/llama.go/server/routes"
)

func main() {
    // 解析模型名称
    name := model.ParseName("hf.co/unsloth/llama-3-8b:Q4_K_M")

    // 定义进度回调
    progressFn := func(resp api.ProgressResponse) {
        if resp.Total > 0 {
            percent := float64(resp.Completed) / float64(resp.Total) * 100
            fmt.Printf("\r%.2f%% - %s", percent, resp.Status)
        } else {
            fmt.Println(resp.Status)
        }
    }

    // 下载模型
    ctx := context.Background()
    if err := routes.PullModel(ctx, name, progressFn); err != nil {
        panic(err)
    }

    fmt.Println("\nDownload complete!")
}
```

## API 端点使用示例

### 使用 curl 下载模型

```bash
# 方式 1: 使用模式匹配
curl -X POST http://localhost:8081/api/pull \
  -H "Content-Type: application/json" \
  -d '{"model": "unsloth/llama-3-8b:Q4_K_M"}'

# 方式 2: 指定完整文件名
curl -X POST http://localhost:8081/api/pull \
  -H "Content-Type: application/json" \
  -d '{"model": "unsloth/llama-3-8b:llama-3-8b-Q4_K_M.gguf"}'

# 方式 3: 使用完整 URL
curl -X POST http://localhost:8081/api/pull \
  -H "Content-Type: application/json" \
  -d '{"model": "https://huggingface.co/unsloth/llama-3-8b/resolve/main/llama-3-8b-Q4_K_M.gguf"}'

# 方式 4: 兼容原有格式
curl -X POST http://localhost:8081/api/pull \
  -H "Content-Type: application/json" \
  -d '{"model": "hf.co/unsloth/llama-3-8b:Q4_K_M"}'
```

## 结构体字段说明

```go
type HuggingFaceModel struct {
    // Host: Hugging Face 主机地址（默认: huggingface.co）
    Host string

    // Namespace: 用户或组织名称
    Namespace string

    // Repo: 仓库名称
    Repo string

    // Branch: Git 分支或标签（默认: main）
    Branch string

    // Filename: 要下载的具体文件
    // 如果为空，将使用 Pattern 进行匹配或自动检测
    Filename string

    // Pattern: 文件名匹配模式
    // 用于在仓库中查找包含该模式的 GGUF 文件
    // 例如: "Q4_K_M", "Q8_0", "f16" 等
    Pattern string
}
```

## 常用方法

| 方法 | 说明 |
|------|------|
| `ParseHuggingFaceModel(s string)` | 解析模型引用字符串 |
| `ToDownloadURL()` | 生成完整的下载 URL |
| `ToRepoURL()` | 生成仓库首页 URL |
| `ToAPIURL()` | 生成 Hugging Face API URL |
| `String()` | 返回简化的字符串表示 |
| `IsValid()` | 检查模型引用是否有效 |
| `GetLocalFilename()` | 获取本地保存的文件名 |
| `ListGGUFFiles()` | 列出仓库中的所有 GGUF 文件 |
| `ResolveFilename()` | 自动解析文件名（基于 Pattern 或自动检测）|

## 错误处理

```go
hf, err := model.ParseHuggingFaceModel("invalid/format")
if err != nil {
    // 处理解析错误
    fmt.Printf("Parse error: %v\n", err)
}

if !hf.IsValid() {
    // 处理无效的模型引用
    fmt.Println("Invalid model reference")
}

// 文件名解析可能失败
if err := hf.ResolveFilename(); err != nil {
    // 可能的原因：
    // - 仓库中没有 GGUF 文件
    // - 找不到匹配 Pattern 的文件
    // - 网络错误
    fmt.Printf("Failed to resolve filename: %v\n", err)
}
```

## 测试

运行单元测试：

```bash
go test ./model -v
```

运行特定测试：

```bash
go test ./model -v -run TestParseHuggingFaceModel
```

## 常见问题

### Q: 如何知道仓库中有哪些 GGUF 文件？

使用 `ListGGUFFiles()` 方法：

```go
hf, _ := model.ParseHuggingFaceModel("unsloth/llama-3-8b")
files, err := hf.ListGGUFFiles()
for _, file := range files {
    fmt.Println(file.Path)
}
```

### Q: Pattern 匹配是区分大小写的吗？

不区分。Pattern 匹配会将文件名和模式都转换为小写进行比较。

### Q: 如何指定不同的分支？

使用完整 URL 格式：

```go
hf, _ := model.ParseHuggingFaceModel(
    "https://huggingface.co/owner/repo/resolve/v1.0/model.gguf",
)
```

### Q: 是否支持私有仓库？

目前的实现使用公开 API，不支持需要认证的私有仓库。


## 完整工作流程示例

```go
package main

import (
    "context"
    "fmt"
    "github.com/Qitmeer/llama.go/model"
    "github.com/Qitmeer/llama.go/server/routes"
)

func downloadModel(modelRef string) error {
    // 1. 解析模型引用
    fmt.Printf("Parsing model reference: %s\n", modelRef)
    hf, err := model.ParseHuggingFaceModel(modelRef)
    if err != nil {
        return fmt.Errorf("parse error: %w", err)
    }

    // 2. 验证
    if !hf.IsValid() {
        return fmt.Errorf("invalid model reference")
    }

    // 3. 列出可用文件（可选）
    fmt.Println("Available GGUF files in repository:")
    files, err := hf.ListGGUFFiles()
    if err != nil {
        return fmt.Errorf("failed to list files: %w", err)
    }
    for _, file := range files {
        fmt.Printf("  - %s\n", file.Path)
    }

    // 4. 解析文件名
    if hf.Filename == "" {
        fmt.Println("Resolving filename...")
        if err := hf.ResolveFilename(); err != nil {
            return fmt.Errorf("failed to resolve filename: %w", err)
        }
        fmt.Printf("Resolved to: %s\n", hf.Filename)
    }

    // 5. 显示下载信息
    fmt.Printf("Download URL: %s\n", hf.ToDownloadURL())
    fmt.Printf("Local filename: %s\n", hf.GetLocalFilename())

    // 6. 执行下载（使用 PullModel）
    // ... 下载逻辑 ...

    return nil
}

func main() {
    if err := downloadModel("unsloth/llama-3-8b:Q4_K_M"); err != nil {
        panic(err)
    }
}
```