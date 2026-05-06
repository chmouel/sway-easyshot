package commands

import (
	"context"
	"fmt"
	"time"

	"sway-easyshot/internal/config"
	"sway-easyshot/internal/external"
	"sway-easyshot/internal/notify"
	"sway-easyshot/internal/state"
)

// OBSHandler provides methods to interact with OBS.
type OBSHandler struct {
	cfg   *config.Config
	state *state.State
}

// NewOBSHandler creates a new OBS handler instance.
func NewOBSHandler(cfg *config.Config, st *state.State) *OBSHandler {
	return &OBSHandler{
		cfg:   cfg,
		state: st,
	}
}

// ToggleRecording toggles OBS recording state (start/stop).
func (h *OBSHandler) ToggleRecording(ctx context.Context) error {
	client, err := external.NewOBSClient(ctx)
	if err != nil {
		_ = notify.Send(2000, h.cfg.ScreenshotIcon, "Failed to connect to OBS")
		return fmt.Errorf("failed to connect to OBS: %w", err)
	}
	defer client.Disconnect() //nolint:errcheck

	status, err := client.Record.GetRecordStatus()
	if err != nil {
		_ = notify.Send(2000, h.cfg.ScreenshotIcon, "Failed to get OBS status")
		return fmt.Errorf("failed to get OBS recording status: %w", err)
	}

	if !status.OutputActive {
		time.Sleep(1 * time.Second)
		if _, err := client.Record.StartRecord(); err != nil {
			return fmt.Errorf("failed to start OBS recording: %w", err)
		}
		h.state.SetOBSState(true, false)
		return nil
	}

	if _, err := client.Record.StopRecord(); err != nil {
		return fmt.Errorf("failed to stop OBS recording: %w", err)
	}
	time.Sleep(2 * time.Second)
	_ = notify.Send(2000, h.cfg.RecordingStopIcon, "Recording has stopped")
	h.state.SetOBSState(false, false)
	return nil
}

// TogglePause toggles OBS pause state (paused/resumed).
func (h *OBSHandler) TogglePause(ctx context.Context) error {
	client, err := external.NewOBSClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to OBS: %w", err)
	}
	defer client.Disconnect() //nolint:errcheck

	// Check state before toggling — OBS doesn't update instantly after the call.
	status, err := client.Record.GetRecordStatus()
	if err != nil {
		return fmt.Errorf("failed to get OBS recording status: %w", err)
	}
	wasPaused := status.OutputPaused

	if _, err := client.Record.ToggleRecordPause(); err != nil {
		return fmt.Errorf("failed to toggle OBS pause: %w", err)
	}

	if wasPaused {
		_ = notify.Send(2000, h.cfg.RecordingStartIcon, "Recording resumed")
		h.state.SetOBSState(true, false)
	} else {
		_ = notify.Send(2000, h.cfg.RecordingPauseIcon, "Recording paused")
		h.state.SetOBSState(true, true)
	}

	return nil
}
