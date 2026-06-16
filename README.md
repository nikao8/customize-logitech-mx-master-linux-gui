# Logitech MX Master Configuration GUI

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/nikao8)

Interface gráfica para gerar e gerenciar a configuração e personalização dos mouses da linha Logitech MX Master no Linux através do [logiops](https://github.com/PixlOne/logiops) (logid) para mouses.

## Instalação

Baixe o `.tar.gz` da sua distribuição na [página de releases](https://github.com/nikao8/logitech-mx-master-customization-linux/releases):

```bash
tar -xzf logid-config-gui-<distro>-<versao>.tar.gz
cd logid-config-gui-<distro>-<versao>
chmod +x install.sh
sudo ./install.sh
sudo ./logid-config-gui
```

O script `install.sh` instala as dependências (Go, bibliotecas gráficas e o logiops se necessário).

## Build manual

```bash
make build-release        # compila o binário
sudo make install         # instala em /usr/local/bin
sudo ./logid-config-gui
```

## Uso

1. Descubra o nome do dispositivo: `sudo logid -v`
2. Abra o GUI, selecione o nome exato e configure botões, DPI, SmartShift e roda lateral
3. Na aba **Service**, clique **Save Configuration** e depois **Install & Start Service**

## Idiomas

Interface disponível em **English** e **Português (Brasil)**.

## Estrutura

```
├── app.go        # Interface Fyne
├── config.go     # Modelo e gerador de /etc/logid.cfg
├── mapping.go    # Mapas de CID, keycodes, traduções
├── service.go    # Gerenciamento do serviço systemd
├── main.go       # Ponto de entrada
├── Makefile      # build-debian, build-arch, build-fedora
└── README.md
```

## Referências

- [logiops - GitHub](https://github.com/PixlOne/logiops)
- [logiops Configuration Wiki](https://github.com/PixlOne/logiops/wiki/Configuration)
- [Arch Wiki - Logitech MX Master](https://wiki.archlinux.org/title/Logitech_MX_Master)
- [Linux input-event-codes.h](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h)
