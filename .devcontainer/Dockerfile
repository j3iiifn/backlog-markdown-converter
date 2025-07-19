FROM node:20

ARG TZ
ENV TZ="$TZ"

# Install basic development tools, Go, and iptables/ipset
RUN apt update && apt install -y less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  wget \
  curl \
  build-essential

# Install Go (latest version 1.24.5) - auto-detect architecture
RUN ARCH=$(dpkg --print-architecture) && \
  if [ "$ARCH" = "amd64" ]; then GOARCH="amd64"; else GOARCH="arm64"; fi && \
  wget -O go.tar.gz "https://go.dev/dl/go1.24.5.linux-${GOARCH}.tar.gz" && \
  tar -C /usr/local -xzf go.tar.gz && \
  rm go.tar.gz

# Set Go environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/home/node/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace, config and Go directories and set permissions
RUN mkdir -p /workspace /home/node/.claude /home/node/go && \
  chown -R node:node /workspace /home/node/.claude /home/node/go

WORKDIR /workspace

RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/0.18.2/git-delta_0.18.2_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_0.18.2_${ARCH}.deb" && \
  rm "git-delta_0.18.2_${ARCH}.deb"

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Default powerline10k theme
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install Claude
RUN npm install -g @anthropic-ai/claude-code

# Install Go tools as node user
USER node
RUN go install golang.org/x/tools/cmd/goimports@latest && \
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Switch back to root for remaining setup
USER root

# Copy and set up firewall script
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
