# see https://github.com/kanielrkirby/nvim-alpine/blob/master/Dockerfile
FROM alpine:latest AS nvim-builder

ARG NVIM_BUILDER_DEPENDENCIES="build-base cmake coreutils linux-headers curl gettext-tiny-dev git"
ARG NVIM_VERSION="v0.12.5"

RUN apk add --no-cache ${NVIM_BUILDER_DEPENDENCIES} && \
  git --version && \
  git clone --depth 1 https://github.com/neovim/neovim.git /tmp/neovim

WORKDIR /tmp/neovim

RUN  git fetch --all --tags -f && \
  git checkout ${NVIM_VERSION} && \
  make CMAKE_BUILD_TYPE=RelWithDebInfo CMAKE_INSTALL_PREFIX=/usr/local/ && \
  make install && \
  strip /usr/local/bin/nvim

FROM alpine:latest AS lua-ls-builder

ARG LUA_LS_BUILDER_DEPENDENCIES="build-base linux-headers python3 git patch tree-sitter-cli ninja bash"
ARG LUA_LS_VERSION="3.19.1"

RUN apk add --no-cache ${LUA_LS_BUILDER_DEPENDENCIES} && \
  git --version && \
  git clone --depth 1 https://github.com/LuaLS/lua-language-server.git /external/lua-language-server

WORKDIR /external/lua-language-server

RUN git fetch --all --tags -f && \
  git checkout ${LUA_LS_VERSION} && \
  chmod +x ./make.sh && \
  ./make.sh

FROM alpine:latest AS runner

COPY --from=nvim-builder /usr/local /usr/local/
COPY --from=nvim-builder /usr/lib/libgcc_s.so.1 /usr/lib/

COPY --from=lua-ls-builder /external /root/external
RUN ln -s /root/external/lua-language-server/bin/lua-language-server /usr/local/bin/lua-language-server

COPY ./resources /resources
RUN ln -s /resources/export/export.lua /usr/local/share/nvim/runtime/export.lua

ENV NVIM_APPNAME=nvim-anydev
COPY ./resources/config /root/.config/${NVIM_APPNAME}

WORKDIR /root/run

CMD exec lua-language-server \
	--configpath="/resources/export/.luarc.json" \
	--doc="/usr/local/share/nvim/runtime" \
	--doc_out_path="/root/run/out"

