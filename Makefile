BINARY_NAME=anydev
EXTERNAL_DIR=external
RESOURCES_DIR=resources
PATCHES_DIR=$(RESOURCES_DIR)/patches

.PHONY=all
all: install

.PHONY=install
install:
	git submodule update --init --recursive
	$(MAKE) lua-parser
	$(MAKE) luadoc-parser
	go mod download

.PHONY=clean
clean:
	cd $(TREE_SITTER_LUA_SRC_DIR) && git checkout . && git clean -fd
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && git checkout . && git clean -fd

# Treesitter lua

TREE_SITTER_LUA_SRC_DIR=$(EXTERNAL_DIR)/tree-sitter-lua
ifeq ($(shell uname),Darwin)
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.dylib
else
    TREE_SITTER_LUA_SRC_FILE := libtree-sitter-lua.so
endif
TREE_SITTER_LUA_SRC_PATH=$(TREE_SITTER_LUA_SRC_DIR)/$(TREE_SITTER_LUA_SRC_FILE)

$(TREE_SITTER_LUA_SRC_PATH):
	@echo "Installing tree-sitter-lua dependencies"
	cd $(TREE_SITTER_LUA_SRC_DIR) && npm install
	@echo "Building tree-sitter-lua"
	$(MAKE) -C $(TREE_SITTER_LUA_SRC_DIR) all

TREE_SITTER_LUA_TARGET_PATH := $(RESOURCES_DIR)/config/parser/lua.so

$(TREE_SITTER_LUA_TARGET_PATH): $(TREE_SITTER_LUA_SRC_PATH)
	cp "$(TREE_SITTER_LUA_SRC_PATH)" "$(TREE_SITTER_LUA_TARGET_PATH)"

.PHONY=lua-parser
lua-parser: $(TREE_SITTER_LUA_TARGET_PATH)

# Treesitter luadoc

TREE_SITTER_LUADOC_SRC_DIR=$(EXTERNAL_DIR)/tree-sitter-luadoc

.PHONY=luadoc-parser-patch
luadoc-parser-patch:
	@echo "Patching luadoc gramar"
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && git checkout . && git clean -fd
	patch -p1 -d "$(TREE_SITTER_LUADOC_SRC_DIR)" < "$(PATCHES_DIR)/tree-sitter-luadoc-grammar.patch"
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && npm run build

TREE_SITTER_LUADOC_SRC_FILE=libtree-sitter-luadoc.so
TREE_SITTER_LUADOC_SRC_PATH=$(TREE_SITTER_LUADOC_SRC_DIR)/$(TREE_SITTER_LUADOC_SRC_FILE)

$(TREE_SITTER_LUADOC_SRC_PATH):
	@echo "Installing tree-sitter-luadoc dependencies"
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && npm install
	$(MAKE) luadoc-parser-patch
	@echo "Building tree-sitter-luadoc"
	cd $(TREE_SITTER_LUADOC_SRC_DIR) && cc -o $(TREE_SITTER_LUADOC_SRC_FILE) -I./src src/parser.c -Os -std=c11 -bundle -fPIC

TREE_SITTER_LUADOC_TARGET_PATH=$(RESOURCES_DIR)/config/parser/luadoc.so

$(TREE_SITTER_LUADOC_TARGET_PATH): $(TREE_SITTER_LUADOC_SRC_PATH)
	cp "$(TREE_SITTER_LUADOC_SRC_PATH)" "$(TREE_SITTER_LUADOC_TARGET_PATH)"

.PHONY=luadoc-parser
luadoc-parser: $(TREE_SITTER_LUADOC_TARGET_PATH)

# Dev

.PHONY=generate
generate:
	go generate ./...

.PHONY=extract
lex:
	go run ./cmd/lex
