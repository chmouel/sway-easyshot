package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sway-screenshot/internal/config"
	"sway-screenshot/internal/external"
	"sway-screenshot/internal/notify"
	"sway-screenshot/internal/sway"
)

type ScreenshotHandler struct {
	cfg *config.Config
}

func NewScreenshotHandler(cfg *config.Config) *ScreenshotHandler {
	return &ScreenshotHandler{cfg: cfg}
}

func (h *ScreenshotHandler) CurrentWindowClipboard(ctx context.Context, delay int) error {
	if err := notify.CaptureDelay(delay, "window to clipboard", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	geom, err := sway.GetFocusedWindowGeometry(ctx)
	if err != nil {
		return fmt.Errorf("failed to get window geometry: %w", err)
	}

	time.Sleep(time.Duration(delay) * time.Second)

	data, err := external.Grim(ctx, geom, "", "")
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return external.WlCopy(ctx, data, "image/png")
}

func (h *ScreenshotHandler) CurrentWindowFile(ctx context.Context, delay int) error {
	if err := notify.CaptureDelay(delay, "window to file", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	geom, err := sway.GetFocusedWindowGeometry(ctx)
	if err != nil {
		return fmt.Errorf("failed to get window geometry: %w", err)
	}

	file := h.cfg.GenerateFilename()
	time.Sleep(time.Duration(delay) * time.Second)

	_, err = external.Grim(ctx, geom, "", file)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return notify.Send(3000, h.cfg.ScreenshotIcon, fmt.Sprintf("Screenshot saved: %s", filepath.Base(file)))
}

func (h *ScreenshotHandler) CurrentScreenClipboard(ctx context.Context, delay int, useCurrentScreen bool) error {
	output, err := sway.SelectOutput(ctx, useCurrentScreen)
	if err != nil || output == "" {
		return fmt.Errorf("failed to select output: %w", err)
	}

	if err := notify.CaptureDelay(delay, "screen to clipboard", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	time.Sleep(time.Duration(delay) * time.Second)

	data, err := external.Grim(ctx, "", output, "")
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return external.WlCopy(ctx, data, "image/png")
}

func (h *ScreenshotHandler) SelectionFile(ctx context.Context, delay int) error {
	if err := notify.CaptureDelay(delay, "selection to file", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	geom, err := external.Slurp(ctx, "")
	if err != nil || geom == "" {
		return fmt.Errorf("selection cancelled or failed: %w", err)
	}

	file := h.cfg.GenerateFilename()
	time.Sleep(time.Duration(delay) * time.Second)

	_, err = external.Grim(ctx, geom, "", file)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Show notification with actions
	actions := map[string]string{
		"copyclip": "Copy image",
		"rename":   "Rename",
		"copypath": "Copy path",
		"edit":     "Edit",
	}

	action, err := notify.SendWithActions(30000, h.cfg.ScreenshotIcon, filepath.Base(file), actions)
	if err != nil {
		// Action selection failed, but screenshot was saved
		return notify.Send(5000, h.cfg.ScreenshotIcon, fmt.Sprintf("Screenshot saved: %s", filepath.Base(file)))
	}

	action = strings.TrimSpace(action)

	switch action {
	case "copyclip":
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		return external.WlCopy(ctx, data, "image/png")

	case "copypath":
		return external.WlCopyText(ctx, file)

	case "rename", "edit":
		newname, err := external.Zenity(ctx, "Rename file", filepath.Base(file))
		if err != nil || newname == "" {
			return nil
		}

		ext := filepath.Ext(file)
		if !strings.HasSuffix(newname, ext) {
			newname = newname + ext
		}

		if action == "edit" {
			outputFile := filepath.Join(h.cfg.SaveLocation, newname)
			return external.Satty(ctx, file, outputFile, true)
		}

		newPath := filepath.Join(h.cfg.SaveLocation, newname)
		return os.Rename(file, newPath)
	}

	return nil
}

func (h *ScreenshotHandler) SelectionEdit(ctx context.Context, delay int) error {
	if err := notify.CaptureDelay(delay, "selection edit", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	geom, err := external.Slurp(ctx, "#ff0000ff")
	if err != nil || geom == "" {
		return fmt.Errorf("selection cancelled or failed: %w", err)
	}

	time.Sleep(time.Duration(delay) * time.Second)

	data, err := external.Grim(ctx, geom, "", "")
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Write to temporary file for satty
	tmpFile := fmt.Sprintf("/tmp/screenshot-%d.png", time.Now().Unix())
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	outputFile := filepath.Join(h.cfg.SaveLocation, fmt.Sprintf("screenshot-%s.png", time.Now().Format("20060102-15:04:05")))
	return external.Satty(ctx, tmpFile, outputFile, true)
}

func (h *ScreenshotHandler) SelectionClipboard(ctx context.Context, delay int) error {
	if err := notify.CaptureDelay(delay, "selection to clipboard", h.cfg.ScreenshotIcon); err != nil {
		return err
	}

	geom, err := external.Slurp(ctx, "")
	if err != nil || geom == "" {
		return fmt.Errorf("selection cancelled or failed: %w", err)
	}

	time.Sleep(time.Duration(delay) * time.Second)

	data, err := external.Grim(ctx, geom, "", "")
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	if err := external.WlCopy(ctx, data, "image/png"); err != nil {
		return err
	}

	// Show notification with actions
	actions := map[string]string{
		"save":   "Save",
		"saveai": "Save with AI",
		"edit":   "Edit",
	}

	action, err := notify.SendWithActions(30000, h.cfg.ScreenshotIcon, "Screenshot captured to clipboard", actions)
	if err != nil {
		return nil // Clipboard copy succeeded, ignore action error
	}

	action = strings.TrimSpace(action)

	if action == "" || (action != "save" && action != "saveai" && action != "edit") {
		return nil
	}

	defaultName := filepath.Base(h.cfg.GenerateFilename())

	if action == "saveai" {
		tmpFile := fmt.Sprintf("/tmp/screenshot-%d.png", time.Now().Unix())
		clipData, err := external.WlPaste(ctx, "image/png")
		if err != nil {
			return err
		}

		if err := os.WriteFile(tmpFile, clipData, 0644); err != nil {
			return err
		}
		defer os.Remove(tmpFile)

		aiName, err := external.AIChat(ctx, h.cfg.AIModelImage, tmpFile,
			"identify a filename for that image and return only the slug of the filename, nothing else")

		if err == nil && aiName != "" {
			defaultName = aiName
			if !strings.HasSuffix(defaultName, ".png") {
				defaultName = defaultName + ".png"
			}
		}
	}

	newname, err := external.Zenity(ctx, "File Name", defaultName)
	if err != nil || newname == "" {
		return nil
	}

	if !strings.HasSuffix(newname, ".png") {
		newname = newname + ".png"
	}

	outputFile := filepath.Join(h.cfg.SaveLocation, newname)

	if action == "edit" {
		clipData, err := external.WlPaste(ctx, "image/png")
		if err != nil {
			return err
		}

		tmpFile := fmt.Sprintf("/tmp/screenshot-%d.png", time.Now().Unix())
		if err := os.WriteFile(tmpFile, clipData, 0644); err != nil {
			return err
		}
		defer os.Remove(tmpFile)

		return external.Satty(ctx, tmpFile, outputFile, true)
	}

	// Save action
	clipData, err := external.WlPaste(ctx, "image/png")
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, clipData, 0644); err != nil {
		return err
	}

	// Open in file manager
	return external.Nautilus(ctx, "file://"+outputFile)
}
