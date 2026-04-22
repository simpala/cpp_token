# Llama.cpp Tokenizer for Go

A clean Go wrapper for the `llama.cpp` tokenizer, providing efficient text encoding and decoding using GGUF models. This implementation uses static linking, making the resulting binaries self-contained and easy to distribute.

## Features

- **Static Linking**: No need to manage `.so` or `.dylib` files at runtime.
- **Easy Build**: Integrated `Makefile` handles the C++ build and Go linking automatically.
- **Full Tokenizer API**: Support for encoding, decoding, and piece-by-piece conversion.
- **Memory Efficient**: Loads only the vocabulary (vocab-only mode) to save RAM.

## Prerequisites

- **Go**: 1.21 or later.
- **CMake**: 3.10 or later (to build `llama.cpp`).
- **C++ Compiler**: A modern compiler with C++17 support (GCC or Clang).
- **Make**: To use the provided build automation.

## Quick Start

1. **Clone the repository** (including submodules):
   ```bash
   git clone --recursive https://github.com/simpala/llama.cpp_token.git
   cd llama.cpp_token
   ```

2. **Build the project**:
   ```bash
   make
   ```
   This will build the `llama.cpp` static libraries and then compile the Go examples.

3. **Run the CLI**:
   ```bash
   ./cli-token /path/to/your/model.gguf
   ```

## Usage

```go
import "llamatokenizer/tokenizer"

// Initialize
t, err := tokenizer.New("model.gguf")
defer t.Close()

// Encode
tokens, err := t.Encode("Hello world", true, false, false)

// Decode
text, err := t.Decode(tokens)
```

## Project Structure

- `tokenizer/`: The core Go package containing the CGo wrapper.
- `llama.cpp/`: The source of truth for the underlying tokenizer logic.
- `Makefile`: Automation for the multi-stage build process.
- `main.go`: A simple example of how to use the package.
- `cli_token.go`: An interactive CLI tool for testing tokenization.

## Development

The Go package uses CGo to interface with `llama.cpp`. The build process is configured to link statically against the following libraries produced by the `llama.cpp` build:
- `libllama.a`
- `libggml.a`
- `libggml-base.a`
- `libggml-cpu.a`

The `Makefile` ensures all necessary OpenMP (`-lgomp`) and C++ standard libraries are included during the final link stage.
