package notify

import (
	"fmt"
	"os/exec"
	"strconv"
)

func Send(timeout int, icon, message string) error {
	args := []string{
		"-t", strconv.Itoa(timeout),
	}
	if icon != "" {
		args = append(args, "-i", icon)
	}
	args = append(args, message)

	cmd := exec.Command("notify-send", args...)
	return cmd.Run()
}

func SendWithActions(timeout int, icon, message string, actions map[string]string) (string, error) {
	args := []string{
		"-t", strconv.Itoa(timeout),
	}
	if icon != "" {
		args = append(args, "-i", icon)
	}

	for id, label := range actions {
		args = append(args, "-A", fmt.Sprintf("%s=%s", id, label))
	}
	args = append(args, message)

	cmd := exec.Command("notify-send", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func CaptureDelay(waitSeconds int, label, icon string) error {
	if waitSeconds > 2 {
		msg := fmt.Sprintf("Capturing %s in %d seconds", label, waitSeconds)
		return Send((waitSeconds-1)*1000, icon, msg)
	}
	return nil
}
