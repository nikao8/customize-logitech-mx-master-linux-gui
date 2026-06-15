#!/usr/bin/env bash
set -euo pipefail

echo "==> Instalando dependências para Void Linux..."

sudo xbps-install -S \
    go \
    libXcursor-devel \
    libXrandr-devel \
    Mesa-devel \
    git \
    cmake \
    gcc \
    libevdev-devel \
    libconfig-devel \
    glib-devel \
    pkg-config \
    libXxf86vm-devel

if ! command -v logid &>/dev/null; then
    echo "==> Compilando e instalando logiops (logid)..."
    cd /tmp
    git clone https://github.com/PixlOne/logiops.git
    cd logiops
    mkdir build && cd build
    cmake ..
    make -j"$(nproc)"
    sudo make install
    sudo systemctl enable logid.service
    cd /tmp && rm -rf logiops
    echo "==> logiops instalado com sucesso!"
else
    echo "==> logid já está instalado, pulando."
fi

echo "==> Dependências instaladas. Agora execute 'make build-release' na raiz do projeto."
