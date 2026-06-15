#!/usr/bin/env bash
set -euo pipefail

echo "==> Instalando dependências para openSUSE..."

sudo zypper install -y \
    go \
    libXcursor-devel \
    libXrandr-devel \
    Mesa-libGL-devel \
    git \
    cmake \
    gcc-c++ \
    libevdev-devel \
    systemd-devel \
    libconfig-devel \
    glib2-devel \
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
