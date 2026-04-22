.PHONY: all build-llama build clean run

LLAMA_DIR = llama.cpp
LLAMA_BUILD_DIR = $(LLAMA_DIR)/build
NPROC = $(shell nproc 2>/dev/null || echo 4)

all: build

$(LLAMA_BUILD_DIR)/CMakeCache.txt:
	cmake -S $(LLAMA_DIR) -B $(LLAMA_BUILD_DIR) \
		-DLLAMA_BUILD_EXAMPLES=OFF \
		-DLLAMA_BUILD_TESTS=OFF \
		-DLLAMA_BUILD_SERVER=OFF \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=OFF

build-llama: $(LLAMA_BUILD_DIR)/CMakeCache.txt
	cmake --build $(LLAMA_BUILD_DIR) --config Release -j $(NPROC)

build: build-llama
	go build -o llama-token main.go
	go build -o cli-token cli_token.go

run: build
	@echo "Usage: ./cli-token <model_path>"
	@echo "Example: ./cli-token models/qwen.gguf"

clean:
	rm -rf $(LLAMA_BUILD_DIR)
	rm -f llama-token cli-token
