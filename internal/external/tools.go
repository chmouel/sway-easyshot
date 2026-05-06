package external

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/andreykaipov/goobs"
)

// Grim captures a screenshot
func Grim(ctx context.Context, geometry, output, filename string) ([]byte, error) {
	args := []string{"-t", "png"}

	if geometry != "" {
		args = append(args, "-g", geometry)
	}
	if output != "" {
		args = append(args, "-o", output)
	}

	if filename == "" {
		args = append(args, "-")
	} else {
		args = append(args, filename)
	}

	cmd := exec.CommandContext(ctx, "grim", args...)

	if filename == "" {
		return cmd.Output()
	}

	return nil, cmd.Run()
}

// Slurp performs interactive region selection
func Slurp(ctx context.Context, color string) (string, error) {
	args := []string{}
	if color != "" {
		args = append(args, "-c", color)
	}

	cmd := exec.CommandContext(ctx, "slurp", args...) //nolint:gosec
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// WlCopy copies data to clipboard
func WlCopy(ctx context.Context, data []byte, mimeType string) error {
	cmd := exec.CommandContext(ctx, "wl-copy", "-t", mimeType) //nolint:gosec
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

// WlCopyText copies text to clipboard
func WlCopyText(ctx context.Context, text string) error {
	return WlCopy(ctx, []byte(text), "text/plain")
}

// WlPaste pastes from clipboard
func WlPaste(ctx context.Context, mimeType string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "wl-paste", "--type", mimeType) //nolint:gosec
	return cmd.Output()
}

// StartWfRecorder starts video recording
func StartWfRecorder(ctx context.Context, geometry, output, filename string, audio bool) (*exec.Cmd, error) {
	args := []string{}

	if geometry != "" {
		args = append(args, "-g", geometry)
	}
	if output != "" {
		args = append(args, "-o", output)
	}
	if audio {
		args = append(args, "-a")
	}

	args = append(args, "-f", filename)

	cmd := exec.CommandContext(ctx, "wf-recorder", args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

// Satty opens the satty image editor
func Satty(ctx context.Context, inputFile, outputFile string, earlyExit bool) error {
	args := []string{
		"--filename", inputFile,
		"--output-filename", outputFile,
	}
	if earlyExit {
		args = append(args, "--early-exit")
	}

	cmd := exec.CommandContext(ctx, "satty", args...) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Zenity shows a text entry dialog
func Zenity(ctx context.Context, text, entryText string) (string, error) {
	args := []string{
		"--entry",
		"--text", text,
		"--entry-text", entryText,
	}

	cmd := exec.CommandContext(ctx, "zenity", args...) //nolint:gosec
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// AIChat uses aichat to generate a filename
func AIChat(ctx context.Context, model, imagePath, prompt string) (string, error) {
	args := []string{
		"--model", model,
		"--file", imagePath,
		prompt,
	}

	cmd := exec.CommandContext(ctx, "aichat", args...) //nolint:gosec
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// Ffmpeg converts video files
func Ffmpeg(ctx context.Context, inputFile, outputFile string) error {
	args := []string{
		"-i", fmt.Sprintf("file:%s", inputFile),
		"-vf", "scale='min(1920,iw)':-2",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-crf", "23",
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		outputFile,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func passGet(ctx context.Context, key, fallback string) string {
	out, err := exec.CommandContext(ctx, "pass", "show", key).Output() //nolint:gosec
	if err != nil {
		return fallback
	}
	return strings.TrimSpace(string(out))
}

// NewOBSClient creates a goobs client using connection settings from the pass store.
// Keys: obs/host (default localhost), obs/port (default 4455), obs/password.
func NewOBSClient(ctx context.Context) (*goobs.Client, error) {
	host := passGet(ctx, "obs/host", "127.0.0.1")
	port := passGet(ctx, "obs/port", "4455")
	password := passGet(ctx, "obs/password", "")
	return goobs.New(fmt.Sprintf("%s:%s", host, port), goobs.WithPassword(password))
}

// Wofi shows a selection menu
func Wofi(ctx context.Context, prompt string, options []string) (string, error) {
	args := []string{
		"--dmenu",
		"--prompt", prompt,
	}

	cmd := exec.CommandContext(ctx, "wofi", args...) //nolint:gosec
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// Nautilus opens a file in nautilus
func Nautilus(ctx context.Context, fileURI string) error {
	cmd := exec.CommandContext(ctx, "nautilus", fileURI) //nolint:gosec
	return cmd.Start()
}

// CleanupOldFiles removes files older than the specified duration
func CleanupOldFiles(ctx context.Context, directory string, olderThan time.Duration) error {
	beforeTime := fmt.Sprintf("%dd", int(olderThan.Hours()/24))

	cmd := exec.CommandContext(ctx, "fd", //nolint:gosec
		"-t", "f",
		"--full-path", directory,
		"--changed-before", beforeTime,
		"--exec-batch", "rm", "-vf",
	)

	return cmd.Run()
}
