package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"llamatokenizer/tokenizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli_token.go <model.gguf>")
		fmt.Println("Example: cli_token.go /path/to/model.Q8_0.gguf")
		os.Exit(1)
	}

	modelPath := os.Args[1]

	t, err := tokenizer.New(modelPath)
	if err != nil {
		log.Fatalf("Failed to load tokenizer from %s: %v", modelPath, err)
	}
	defer t.Close()

	fmt.Printf("Tokenizer loaded from: %s\n", modelPath)
	fmt.Println("Enter text to tokenize (type 'quit' or 'exit' to stop):")
	fmt.Println("========================================")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		// Exit commands
		lowerInput := strings.ToLower(input)
		if lowerInput == "quit" || lowerInput == "exit" || lowerInput == "q" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		// Encode text to tokens
		tokens, err := t.Encode(input, true, false, false) // addBOS=true
		if err != nil {
			fmt.Printf("Error encoding: %v\n", err)
			continue
		}

		// Decode tokens back to string
		decoded, err := t.Decode(tokens)
		if err != nil {
			fmt.Printf("Error decoding: %v\n", err)
			continue
		}

		// Display results
		fmt.Printf("\nInput:    %q\n", input)
		fmt.Printf("Tokens:   %d tokens\n", len(tokens))
		fmt.Printf("Token IDs: %v\n", tokens)
		fmt.Printf("Decoded:  %q\n", decoded)

		// Show roundtrip check (may differ due to BOS token)
		if strings.HasPrefix(decoded, "▁") || decoded == input {
			fmt.Printf("Roundtrip: OK\n")
		}

		// Individual token breakdown
		fmt.Println("\nToken breakdown:")
		for i, tok := range tokens {
			piece := t.TokenToPiece(tok)
			fmt.Printf("  [%d] Token ID %d → %q\n", i, tok, piece)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}
