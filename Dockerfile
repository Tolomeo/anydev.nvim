# see https://github.com/kanielrkirby/nvim-alpine/blob/master/Dockerfile
FROM alpine:latest AS neovim-builder

ARG NEOVIM_BUILDER_DEPS="autoconf automake cmake curl g++ git gettext gettext-dev libtool make ninja openssl pkgconfig unzip binutils wget"
ARG NEOVIM_VERSION=stable

RUN apk add --no-cache ${NEOVIM_BUILDER_DEPS} && \
  git --version && \
  git clone https://github.com/neovim/neovim.git /tmp/neovim

WORKDIR /tmp/neovim

RUN  git fetch --all --tags -f && \
  git checkout ${NEOVIM_VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM node:20-alpine3.17 AS config-builder

ARG CONFIG_BUILDER_DEPS="build-base linux-headers python3 git patch tree-sitter-cli ninja bash"
ARG LUA_LANGUAGE_SERVER_VERSION="3.16.4"

RUN apk add --no-cache ${CONFIG_BUILDER_DEPS} && \
  git --version && \
  git clone --depth 1 https://github.com/tree-sitter-grammars/tree-sitter-luadoc.git /tmp/tree-sitter-luadoc && \
  git clone --depth 1 https://github.com/LuaLS/lua-language-server.git /external/lua-language-server && \
  git clone --depth 1 https://github.com/folke/lazydev.nvim.git /external/lazydev.nvim

COPY ./resources /resources

WORKDIR /tmp/tree-sitter-luadoc

RUN npm install && \
  patch -p1 -d "/tmp/tree-sitter-luadoc" < "/resources/patches/tree-sitter-luadoc-grammar.patch" && \
  tree-sitter generate && \
  npx node-gyp build && \
  cc -shared -o luadoc.so -I./src src/parser.c -Os -std=c11 -fPIC && \
  cp luadoc.so /resources/config/parser/

WORKDIR /external/lua-language-server

RUN git fetch --all --tags -f && \
  git checkout ${LUA_LANGUAGE_SERVER_VERSION} && \
  chmod +x ./make.sh && \
  ./make.sh

FROM alpine:latest

COPY --from=neovim-builder /usr/local /usr/local/
COPY --from=neovim-builder /lib/ld-musl-aarch64.so.1 /lib/
COPY --from=neovim-builder /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=neovim-builder /usr/lib/libintl.so.8 /usr/lib/
COPY --from=config-builder /external /root/external
COPY --from=config-builder /resources/config /root/.config/nvim

RUN ln -s /root/external/lua-language-server/bin/lua-language-server /usr/local/bin/lua-language-server

WORKDIR /root/dev

CMD ["nvim"]
