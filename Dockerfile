# see https://github.com/kanielrkirby/nvim-alpine/blob/master/Dockerfile
FROM alpine:latest AS neovim-builder

ARG NEOVIM_DEPS="autoconf automake cmake curl g++ git gettext gettext-dev libtool make ninja openssl pkgconfig unzip binutils wget"
ARG NEOVIM_TARGET=stable

RUN apk add --no-cache ${NEOVIM_DEPS} && \
  git --version && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim && \
  cd /tmp/neovim && \
  git fetch --all --tags -f && \
  git checkout ${NEOVIM_TARGET} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM node:20-alpine3.17 AS config-builder

ARG TREESITTER_GRAMMAR_DEPS="build-base python3 git patch tree-sitter-cli"

RUN apk add --no-cache ${TREESITTER_GRAMMAR_DEPS} && \
  git --version && \
  git clone --depth 1 https://github.com/tree-sitter-grammars/tree-sitter-luadoc.git /external/tree-sitter-luadoc && \
  git clone --depth 1 https://github.com/folke/lazydev.nvim.git /external/lazydev.nvim && \
  mkdir -p /tmp/.config/nvim

COPY ./resources /resources

WORKDIR /external/tree-sitter-luadoc

RUN npm install && \
  patch -p1 -d "/external/tree-sitter-luadoc" < "/resources/patches/tree-sitter-luadoc-grammar.patch" && \
	tree-sitter generate && \
	npx node-gyp build && \
  cc -shared -o luadoc.so -I./src src/parser.c -Os -std=c11 -fPIC && \
  cp luadoc.so /resources/config/parser/

FROM alpine:latest

COPY --from=neovim-builder /usr/local /usr/local/
RUN true
COPY --from=neovim-builder /lib/ld-musl-aarch64.so.1 /lib/
RUN true
COPY --from=neovim-builder /usr/lib/libgcc_s.so.1 /usr/lib/
RUN true
COPY --from=neovim-builder /usr/lib/libintl.so.8 /usr/lib/

RUN mkdir -p /root/.config/nvim

COPY  --from=config-builder /external /root/external
COPY --from=config-builder /resources/config /root/.config/nvim

WORKDIR /home/dev

CMD ["nvim"]
