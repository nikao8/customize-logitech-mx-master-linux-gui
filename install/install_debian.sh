#!/usr/bin/env bash
set -euo pipefail

echo "==> Instalando dependências para Debian/Ubuntu..."

sudo apt update

sudo apt install -y \
    golang-go \
    libgl1-mesa-dev \
    xorg-dev \
    libxrandr-dev \
    git \
    cmake \
    build-essential \
    libevdev-dev \
    libudev-dev \
    libconfig++-dev \
    libglib2.0-dev \
    libxxf86vm-dev

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
