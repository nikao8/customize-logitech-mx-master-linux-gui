# Logitech MX Master Configuration GUI

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/nikao8)

**Português (Brasil)** | [English](#english)

---

## Português (Brasil)

Interface gráfica para gerar e gerenciar a configuração e personalização dos mouses da linha Logitech MX Master no Linux através do [logiops](https://github.com/PixlOne/logiops) (logid) para mouses.

### Instalação

Baixe o `.tar.gz` da sua distribuição na [página de releases](https://github.com/nikao8/logitech-mx-master-customization-linux-gui/releases):

```bash
tar -xzf logid-config-gui-<distro>-<versao>.tar.gz
cd logid-config-gui-<distro>-<versao>
chmod +x install.sh
sudo ./install.sh
sudo ./logid-config-gui
```

O script `install.sh` instala as dependências (Go, bibliotecas gráficas e o logiops se necessário).

### Build manual

```bash
make build-release        # compila o binário
sudo make install         # instala em /usr/local/bin
sudo ./logid-config-gui
```

### Uso

1. Descubra o nome do dispositivo: `sudo logid -v`
2. Abra o GUI, selecione o nome exato e configure botões, DPI, SmartShift e roda lateral
3. Na aba **Service**, clique **Save Configuration** e depois **Install & Start Service**
4. A cada alteração que fizer, lembre de clicar em **Save Configuration** e **Restart Service**

---

## English

Graphical interface to generate and manage configuration and customization for Logitech MX Master series mice on Linux via [logiops](https://github.com/PixlOne/logiops) (logid).

### Installation

Download the `.tar.gz` for your distribution from the [releases page](https://github.com/nikao8/logitech-mx-master-customization-linux-gui/releases):

```bash
tar -xzf logid-config-gui-<distro>-<version>.tar.gz
cd logid-config-gui-<distro>-<version>
chmod +x install.sh
sudo ./install.sh
sudo ./logid-config-gui
```

The `install.sh` script installs dependencies (Go, graphical libraries, and logiops if needed).

### Manual build

```bash
make build-release        # builds the binary
sudo make install         # installs to /usr/local/bin
sudo ./logid-config-gui
```

### Usage

1. Find your device name: `sudo logid -v`
2. Open the GUI, select the exact device name, and configure buttons, DPI, SmartShift, and the side wheel
3. In the **Service** tab, click **Save Configuration** and then **Install & Start Service**
4. After any change, remember to click **Save Configuration** and **Restart Service**

## Languages

Interface available in **English** and **Português (Brasil)**.

## References

- [logiops - GitHub](https://github.com/PixlOne/logiops)
- [logiops Configuration Wiki](https://github.com/PixlOne/logiops/wiki/Configuration)
- [Arch Wiki - Logitech MX Master](https://wiki.archlinux.org/title/Logitech_MX_Master)
- [Linux input-event-codes.h](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h)
