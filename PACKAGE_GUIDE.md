# Packaging Guide: Creating a Reusable Go Module

This document explains how to turn this repository into a reusable Go package that other projects can import.

## Option 1: Use as a Local Module (Recommended for Development)

If you have another Go project on your local machine and want to use this tokenizer, you can use the `replace` directive in your project's `go.mod`.

1. In your **other** project:
   ```bash
   go mod edit -replace llamatokenizer=/path/to/llama.cpp_token
   go get llamatokenizer/tokenizer
   ```
2. Import it in your code:
   ```go
   import "llamatokenizer/tokenizer"
   ```

## Option 2: Turn into a Standalone Library

To make this a clean library that others can `go get` easily, follow these steps:

### 1. Rename the Module
Choose a canonical name (usually your GitHub URL).
```bash
go mod edit -module github.com/youruser/llama-tokenizer-go
```

### 2. Update Internal Imports
If you rename the module, update `main.go` and `cli_token.go` to use the new module name:
```go
import "github.com/youruser/llama-tokenizer-go/tokenizer"
```

### 3. Handling the C++ Dependency
This is the most challenging part of CGo packaging. You have three choices:

#### A. The "Build First" Approach (Current Setup)
Users must run `make build-llama` before they can `go build` their own project. This is the cleanest for developers but requires a manual step.

#### B. The "Static Vendoring" Approach
You can pre-compile the `.a` files for different architectures (linux-amd64, darwin-arm64, etc.) and use conditional CGo flags:
```go
// #cgo linux,amd64 LDFLAGS: ${SRCDIR}/libs/linux_amd64/libllama.a ...
```
*Note: This makes the repository very large but makes `go get` work instantly.*

#### C. The "Header-Only" / "Source" Approach
You can try to include the `llama.cpp` source files directly in the `tokenizer/` directory. CGo will then compile the C++ code automatically when someone runs `go build`.
*Note: This is difficult with `llama.cpp` due to its complex CMake-based build system and many files.*

## Recommended Strategy for Distribution

1. **Keep the Submodule**: Maintain `llama.cpp` as a git submodule.
2. **Provide a Helper Script**: Include a script that automates the `cmake` and `make` steps.
3. **Use Environment Variables**: Allow users to point to a custom `llama.cpp` installation:
   ```go
   // tokenizer.go
   #cgo LDFLAGS: -L${LLAMA_PATH} -lllama
   ```

## Best Practices
- **Versioning**: Use git tags (e.g., `v1.0.0`) so Go's proxy can cache your module.
- **Minimal API**: Keep the `tokenizer` package focused. Avoid putting main applications (`main.go`) inside the package directory.
- **Cross-Compilation**: Remember that CGo makes cross-compilation difficult. Users on Windows or Mac will need their own build of `llama.cpp`.
