package tokenizer

import (
    "fmt"
    "sync"
    "unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../llama.cpp/include -I${SRCDIR}/../llama.cpp/ggml/include
#cgo LDFLAGS: -L${SRCDIR}/../llama.cpp/build/src -L${SRCDIR}/../llama.cpp/build/ggml/src -lllama -lggml -lggml-base -lggml-cpu -lgomp -lm -lstdc++ -lpthread -ldl
#cgo CXXFLAGS: -std=c++17 -O3 -I${SRCDIR}/../llama.cpp/include -I${SRCDIR}/../llama.cpp/ggml/include

#include <stdlib.h>
#include "llama.h"
*/
import "C"

var backendInitOnce sync.Once

type Tokenizer struct {
    model *C.struct_llama_model
}

func New(modelPath string) (*Tokenizer, error) {
    backendInitOnce.Do(func() {
        C.llama_backend_init()
    })

    cpath := C.CString(modelPath)
    defer C.free(unsafe.Pointer(cpath))

    params := C.llama_model_default_params()
    params.vocab_only = C.bool(true)

    model := C.llama_model_load_from_file(cpath, params)
    if model == nil {
        return nil, fmt.Errorf("failed to load vocabulary from %s", modelPath)
    }

    return &Tokenizer{model: model}, nil
}

func (t *Tokenizer) Close() {
    if t.model != nil {
        C.llama_model_free(t.model)
        t.model = nil
    }
}

// BackendFree shuts down the llama + ggml backend.
// Call this once at the very end of the program.
func BackendFree() {
    C.llama_backend_free()
}

// Encode text → tokens
func (t *Tokenizer) Encode(text string, addBOS, addEOS, parseSpecial bool) ([]int32, error) {
    ctext := C.CString(text)
    defer C.free(unsafe.Pointer(ctext))

    const maxTokens = 65536
    tokens := make([]C.llama_token, maxTokens)

    n := C.llama_tokenize(
        C.llama_model_get_vocab(t.model),
        ctext,
        C.int32_t(len(text)),
        &tokens[0],
        C.int32_t(maxTokens),
        C.bool(addBOS),
        C.bool(parseSpecial),
    )

    if n < 0 {
        return nil, fmt.Errorf("tokenization failed")
    }

    result := make([]int32, n)
    for i := 0; i < int(n); i++ {
        result[i] = int32(tokens[i])
    }
    return result, nil
}

// Decode tokens → string
func (t *Tokenizer) Decode(tokens []int32) (string, error) {
    if len(tokens) == 0 {
        return "", nil
    }

    cTokens := make([]C.llama_token, len(tokens))
    for i, tok := range tokens {
        cTokens[i] = C.llama_token(tok)
    }

    vocab := C.llama_model_get_vocab(t.model)

    // Initial buffer estimate
    bufSize := C.int32_t(len(tokens)*8 + 512)
    buf := make([]byte, bufSize)

    n := C.llama_detokenize(
        vocab,
        &cTokens[0],
        C.int32_t(len(tokens)),
        (*C.char)(unsafe.Pointer(&buf[0])),
        bufSize,
        C.bool(true),  // remove_special
        C.bool(false), // unparse_special
    )

    if n < 0 {
        // Buffer too small, n is -required_size
        bufSize = -n
        buf = make([]byte, bufSize)
        n = C.llama_detokenize(
            vocab,
            &cTokens[0],
            C.int32_t(len(tokens)),
            (*C.char)(unsafe.Pointer(&buf[0])),
            bufSize,
            C.bool(true),
            C.bool(false),
        )
    }

    if n < 0 {
        return "", fmt.Errorf("detokenization failed")
    }

    return string(buf[:n]), nil
}

// TokenToPiece (useful for streaming)
func (t *Tokenizer) TokenToPiece(token int32) string {
    vocab := C.llama_model_get_vocab(t.model)
    var buf [128]byte
    n := C.llama_token_to_piece(
        vocab,
        C.llama_token(token),
        (*C.char)(unsafe.Pointer(&buf[0])),
        C.int32_t(len(buf)),
        0,
        C.bool(true),
    )
    if n > 0 && n <= C.int32_t(len(buf)) {
        return string(buf[:n])
    }
    if n < 0 {
        // Buffer too small, n is -required_size
        actualBuf := make([]byte, -n)
        n = C.llama_token_to_piece(
            vocab,
            C.llama_token(token),
            (*C.char)(unsafe.Pointer(&actualBuf[0])),
            C.int32_t(len(actualBuf)),
            0,
            C.bool(true),
        )
        if n > 0 {
            return string(actualBuf[:n])
        }
    }
    return ""
}
