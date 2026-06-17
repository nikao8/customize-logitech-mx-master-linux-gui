package src

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/** bkp service
const serviceContent = `[Unit]
Description=Logitech MX Master Configuration Daemon
After=multi-user.target

[Service]
Type=simple
ExecStart=/usr/bin/logid
Restart=always
RestartSec=3
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`

*/

const serviceContent = `[Unit]
Description=Logitech MX Master Configuration Daemon

[Service]
Type=simple
ExecStart=/usr/bin/logid
Restart=always
RestartSec=3
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`

func GenerateServiceContent() string {
	return serviceContent
}

func WriteConfigFile(content string, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func WriteServiceFile(path string) error {
	return os.WriteFile(path, []byte(serviceContent), 0644)
}

func EnableAndStartService() (string, error) {
	var output strings.Builder

	cmd := exec.Command("systemctl", "daemon-reload")
	out, err := cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl daemon-reload\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("daemon-reload failed: %w\n%s", err, string(out))
	}

	cmd = exec.Command("systemctl", "enable", "logid.service")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl enable logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("enable failed: %w\n%s", err, string(out))
	}

	cmd = exec.Command("systemctl", "start", "logid.service")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl start logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("start failed: %w\n%s", err, string(out))
	}

	return output.String(), nil
}

func StopService() (string, error) {
	cmd := exec.Command("systemctl", "stop", "logid.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("stop failed: %w\n%s", err, string(out))
	}
	return string(out), nil
}

func RestartService() (string, error) {
	cmd := exec.Command("systemctl", "restart", "logid.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("restart failed: %w\n%s", err, string(out))
	}
	return string(out), nil
}

func RemoveService() (string, error) {
	var output strings.Builder

	cmd := exec.Command("systemctl", "stop", "logid.service")
	out, err := cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl stop logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("stop failed: %w\n%s", err, string(out))
	}

	cmd = exec.Command("systemctl", "disable", "logid.service")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl disable logid.service\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("disable failed: %w\n%s", err, string(out))
	}

	if err := os.Remove("/etc/systemd/system/logid.service"); err != nil && !os.IsNotExist(err) {
		output.WriteString(fmt.Sprintf("Error removing service file: %v\n", err))
		return output.String(), fmt.Errorf("remove service file failed: %w", err)
	}
	output.WriteString("$ rm /etc/systemd/system/logid.service\n")

	cmd = exec.Command("systemctl", "daemon-reload")
	out, err = cmd.CombinedOutput()
	output.WriteString(fmt.Sprintf("$ systemctl daemon-reload\n%s\n", string(out)))
	if err != nil {
		return output.String(), fmt.Errorf("daemon-reload failed: %w\n%s", err, string(out))
	}

	return output.String(), nil
}

func IsServiceRunning() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "logid.service")
	return cmd.Run() == nil
}

func EnsureConfigDir() error {
	return os.MkdirAll("/etc", 0755)
}
