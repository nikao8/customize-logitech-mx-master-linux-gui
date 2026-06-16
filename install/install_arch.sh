#!/usr/bin/env bash
set -euo pipefail

echo "==> Instalando dependências para Arch Linux..."

sudo pacman -S --needed \
    go \
    git \
    base-devel \
    cmake \
    libevdev \
    systemd-libs \
    libconfig \
    libxxf86vm

if ! command -v logid &>/dev/null; then
    if command -v yay &>/dev/null; then
        echo "==> Instalando logiops via AUR (yay)..."
        yay -S logiops
    elif command -v paru &>/dev/null; then
        echo "==> Instalando logiops via AUR (paru)..."
        paru -S logiops
    else
        echo "==> AUR helper não encontrado. Compilando do fonte..."
        cd /tmp
        git clone https://github.com/PixlOne/logiops.git
        cd logiops
        mkdir build && cd build
        cmake ..
        make -j"$(nproc)"
        sudo make install
        cd /tmp && rm -rf logiops
        echo "==> logiops instalado com sucesso!"
    fi
    sudo systemctl enable logid.service
else
    echo "==> logid já está instalado, pulando."
fi

echo "==> Dependências instaladas. Agora execute 'make build-release' na raiz do projeto."
