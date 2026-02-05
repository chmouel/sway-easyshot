package state

import (
	"sync"

	"sway-screenshot/pkg/protocol"
)

type State struct {
	mu            sync.RWMutex
	recording     bool
	paused        bool
	recordingFile string
	recordingPID  int
	obsRecording  bool
	obsPaused     bool
}

func NewState() *State {
	return &State{}
}

func (s *State) GetState() *protocol.State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &protocol.State{
		Recording:     s.recording,
		Paused:        s.paused,
		RecordingFile: s.recordingFile,
		OBSRecording:  s.obsRecording,
		OBSPaused:     s.obsPaused,
	}
}

func (s *State) SetRecording(recording bool, file string, pid int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.recording = recording
	s.recordingFile = file
	s.recordingPID = pid
}

func (s *State) SetOBSState(recording bool, paused bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.obsRecording = recording
	s.obsPaused = paused
}

func (s *State) GetRecordingPID() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.recordingPID
}

func (s *State) SetPaused(paused bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paused = paused
}

func (s *State) GetWaybarStatus() *protocol.WaybarStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Priority: wf-recorder > OBS
	if s.recording {
		if s.paused {
			return &protocol.WaybarStatus{
				Text:    "󰏤",
				Tooltip: "Recording paused",
				Class:   "paused",
				Alt:     "paused",
			}
		}
		return &protocol.WaybarStatus{
			Text:    "󰑊",
			Tooltip: "Recording in progress",
			Class:   "recording",
			Alt:     "recording",
		}
	}

	if s.obsRecording {
		if s.obsPaused {
			return &protocol.WaybarStatus{
				Text:    "󰏤",
				Tooltip: "OBS recording paused",
				Class:   "paused",
				Alt:     "paused",
			}
		}
		return &protocol.WaybarStatus{
			Text:    "󰑊",
			Tooltip: "OBS recording in progress",
			Class:   "recording",
			Alt:     "recording",
		}
	}

	return &protocol.WaybarStatus{
		Text:    "󰕧",
		Tooltip: "Ready for screenshot/recording",
		Class:   "idle",
		Alt:     "idle",
	}
}
