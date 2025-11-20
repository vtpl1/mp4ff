FROM ubuntu:20.04

# ---- Global Args ----
ARG DEBIAN_FRONTEND=noninteractive
ARG GO_VERSION=1.25.4
ARG NODE_VERSION=22
ARG PROTOC_VERSION=32.1
ARG LLVM_VERSION=21
ARG PB_REL="https://github.com/protocolbuffers/protobuf/releases"
ARG BUF_VERSION="1.59.0"
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG CMAKE_VERSION=4.1.2
ARG PYTHON_MAIN_VERSION=3.14
ARG PYTHON_VERSION=3.14.0

# ---- Base Env ----
ENV SHELL=/bin/bash \
    LANG=en_US.utf-8 \
    LC_ALL=en_US.utf-8

# ---- Base Tools ----
RUN apt update && \
    apt install -y \
        build-essential \
        ca-certificates \
        curl \
        default-jdk \
        gdb \
        git \
        iputils-ping \
        libbz2-dev \
        libffi-dev \
        libgdbm-dev \
        liblzma-dev \
        libncursesw5-dev \
        libreadline-dev \
        libsqlite3-dev \
        libssl-dev \
        libusb-0.1-4 \
        lsb-release \
        ninja-build \
        ocl-icd-opencl-dev \
        pkg-config \
        software-properties-common \
        strace \
        tar \
        tk-dev \
        unzip \
        uuid-dev \
        valgrind \
        wget \
        zip \
        zlib1g-dev && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*


# ---- Install CMake ----
RUN wget -O cmake.sh https://github.com/Kitware/CMake/releases/download/v${CMAKE_VERSION}/cmake-${CMAKE_VERSION}-linux-x86_64.sh && \
    sh cmake.sh --prefix=/usr/local/ --exclude-subdir && rm -rf cmake.sh

# ---- Install LLVM/Clang ----
RUN wget https://apt.llvm.org/llvm.sh && chmod +x llvm.sh && ./llvm.sh ${LLVM_VERSION} all && \
    apt install -y clang-tools-${LLVM_VERSION} \
    lld-${LLVM_VERSION} clang-format-${LLVM_VERSION} \
    clang-tidy-${LLVM_VERSION} && \
    apt clean && rm -rf /var/lib/apt/lists/*

# ---- Set Clang as Default ----
ENV CC=clang-${LLVM_VERSION} \
    CXX=clang++-${LLVM_VERSION} \
    LLVM_AR=llvm-ar-${LLVM_VERSION} \
    LLVM_NM=llvm-nm-${LLVM_VERSION} \
    LLVM_RANLIB=llvm-ranlib-${LLVM_VERSION} \
    LD=ld.lld-${LLVM_VERSION}

# Symlinks for convenience
RUN ln -sf /usr/bin/clang-format-${LLVM_VERSION} /usr/local/bin/clang-format && \
    ln -sf /usr/bin/clang-tidy-${LLVM_VERSION} /usr/local/bin/clang-tidy && \
    ln -sf /usr/bin/lldb-${LLVM_VERSION} /usr/local/bin/lldb && \
    ln -sf /usr/bin/llvm-ar-${LLVM_VERSION} /usr/local/bin/llvm-ar && \
    ln -sf /usr/bin/llvm-ranlib-${LLVM_VERSION} /usr/local/bin/llvm-ranlib && \
    ln -sf /usr/bin/llvm-profdata-${LLVM_VERSION} /usr/local/bin/llvm-profdata    


RUN cd /tmp && \
    wget https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz && \
    tar -xzf Python-${PYTHON_VERSION}.tgz && \
    cd Python-${PYTHON_VERSION} && \
    ./configure --enable-optimizations --with-lto --enable-shared && \
    make -j"$(nproc)" && \
    make altinstall && \
    ldconfig && \
    cd /tmp && \
    rm -rf Python-${PYTHON_VERSION} Python-${PYTHON_VERSION}.tgz

# Ensure `python3` and `pip3` point to the new version
RUN ln -sf /usr/local/bin/python${PYTHON_MAIN_VERSION} /usr/local/bin/python3 && \
    ln -sf /usr/local/bin/pip${PYTHON_MAIN_VERSION} /usr/local/bin/pip3 && \
    ln -sf /usr/local/bin/python${PYTHON_MAIN_VERSION} /usr/local/bin/python && \
    ln -sf /usr/local/bin/pip${PYTHON_MAIN_VERSION} /usr/local/bin/pip


# ---- Install Go ----
RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz" && \
    rm "go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH="/home/${USERNAME}/go"
ENV GOBIN="${GOPATH}/bin"
ENV PATH="/usr/local/go/bin:${GOBIN}:${PATH}"

RUN mkdir -p /etc/OpenCL/vendors && \
    echo "/usr/lib/x86_64-linux-gnu/libnvidia-opencl.so.1" > /etc/OpenCL/vendors/nvidia.icd

# ---- Create User ----
RUN groupadd --gid ${USER_GID} ${USERNAME} && \
    useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME}

USER ${USERNAME}

# ---- Install Node ----
RUN wget -q -O - https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.3/install.sh | bash && \
    export NVM_DIR="$HOME/.nvm" && \
    . "$NVM_DIR/nvm.sh" && \
    nvm install ${NODE_VERSION}

# ---- Install Protoc & Buf ----
RUN mkdir -p "$HOME/tmp" && \
    wget -q -O "$HOME/tmp/protoc.zip" \
        "${PB_REL}/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip" && \
    unzip "$HOME/tmp/protoc.zip" -d "$HOME/.local" && \
    rm -rf "$HOME/tmp"

RUN BIN="$HOME/.local/bin" && \    
    curl -sSL \
        "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" \
        -o "${BIN}/buf" && \
    chmod +x "${BIN}/buf"

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.2

RUN go install github.com/4meepo/tagalign/cmd/tagalign@latest && \
    go install github.com/client9/misspell/cmd/misspell@latest && \
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest && \
    go install github.com/google/pprof@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install github.com/securego/gosec/v2/cmd/gosec@latest && \
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install golang.org/x/lint/golint@latest && \
    go install golang.org/x/vuln/cmd/govulncheck@latest && \
    go install gotest.tools/gotestsum@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install mvdan.cc/gofumpt@latest


ENV PATH="$PATH:/home/$USERNAME/go/bin"

RUN python -m pip install -U pip poetry bump2version black yapf isort setuptools wheel