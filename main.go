package main

import (
    "fmt"
    "log"

    "llamatokenizer/tokenizer" // change if your go.mod module name is different
)

func main() {
    modelPath := "/media/simpala/New_Volume1/devel/OLLAMA_MODELS/bartowski/Qwen_Qwen3.5-0.8B-Q8_0.gguf" // ← CHANGE THIS

    t, err := tokenizer.New(modelPath)
    if err != nil {
        log.Fatal(err)
    }
    defer t.Close()

    //text := "The quick brown fox jumps over the lazy dog."
    text := "she said the lights turned on by themselves"

    tokens, err := t.Encode(text, true, false, false)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tokens: %d → %v\n", len(tokens), tokens)

    decoded, err := t.Decode(tokens)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Decoded: %q\n", decoded)
    fmt.Printf("Roundtrip OK: %v\n", decoded == text)
}
