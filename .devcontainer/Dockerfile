# See here for image contents:
# https://github.com/microsoft/vscode-dev-containers/tree/v0.209.6/containers/go/.devcontainer/base.Dockerfile

FROM golang:1.20-bullseye

RUN useradd --create-home --shell /bin/bash dev

USER dev

# Install Go extension dependencies.
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0 \
  && go install -v github.com/fatih/gomodifytags@v1.16.0 \
  && go install -v github.com/josharian/impl@v1.1.0 \
  && go install -v github.com/haya14busa/goplay/cmd/goplay@v1.0.0 \
  && go install -v github.com/go-delve/delve/cmd/dlv@v1.21.0 \
  && go install -v honnef.co/go/tools/cmd/staticcheck@v0.4.3 \
  && go install -v golang.org/x/tools/gopls@v0.13.1
