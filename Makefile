BINARY_NAME=anydev
TREE_SITTER_LUA_SRC_DIR=external/tree-sitter-lua

ifeq ($(shell uname),Darwin)
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.dylib
else
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.so
endif

TREE_SITTER_LUA_TARGET_DIR := .config/nvim/parsers
TREE_SITTER_LUA_TARGET_FILE := lua.so

.PHONY=tree-sitter-install
tree-sitter-install:
	@(cd $(TREE_SITTER_LUA_SRC_DIR) && npm install)

.PHONY=tree-sitter-build
tree-sitter-build: tree-sitter-install
	$(MAKE) -C $(TREE_SITTER_LUA_SRC_DIR) all
	cp "$(TREE_SITTER_LUA_SRC_DIR)/$(TREE_SITTER_LUA_SRC_FILE)" "$(TREE_SITTER_LUA_TARGET_DIR)/$(TREE_SITTER_LUA_TARGET_FILE)"

.PHONY=generate
generate:
	@go generate ./...

.PHONY=run
run:
	@go run ./cmd/extract

.PHONY=install
install:
	@go mod download
