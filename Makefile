BINARY_NAME=anydev
TREE_SITTER_TARGET_DIR := .config/nvim/parsers

TREE_SITTER_LUA_SRC_DIR=external/tree-sitter-lua

ifeq ($(shell uname),Darwin)
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.dylib
else
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.so
endif

TREE_SITTER_LUA_TARGET_FILE := lua.so

TREE_SITTER_LUADOC_SRC_DIR=external/tree-sitter-luadoc
TREE_SITTER_LUADOC_SRC_FILE=libtree-sitter-luadoc.so
TREE_SITTER_LUADOC_TARGET_FILE=luadoc.so

.PHONY=tree-sitter-install
tree-sitter-install:
	@echo "Installing tree-sitter-lua dependencies"
	@(cd $(TREE_SITTER_LUA_SRC_DIR) && npm install)
	@echo "Installing tree-sitter-luadoc dependencies"
	@(cd $(TREE_SITTER_LUADOC_SRC_DIR) && npm install)

.PHONY=tree-sitter-build
# tree-sitter-build: tree-sitter-install
tree-sitter-build:
	@echo "Building tree-sitter-lua"
	$(MAKE) -C $(TREE_SITTER_LUA_SRC_DIR) all
	cp "$(TREE_SITTER_LUA_SRC_DIR)/$(TREE_SITTER_LUA_SRC_FILE)" "$(TREE_SITTER_TARGET_DIR)/$(TREE_SITTER_LUA_TARGET_FILE)"
	@echo "Building tree-sitter-luadoc"
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && cc -o $(TREE_SITTER_LUADOC_SRC_FILE) -I./src src/parser.c -Os -std=c11 -bundle -fPIC
	cp "$(TREE_SITTER_LUADOC_SRC_DIR)/$(TREE_SITTER_LUADOC_SRC_FILE)" "$(TREE_SITTER_TARGET_DIR)/$(TREE_SITTER_LUADOC_TARGET_FILE)"

.PHONY=generate
generate:
	@go generate ./...

.PHONY=extract
extract:
	@go run ./cmd/extract

.PHONY=install
install:
	@go mod download
