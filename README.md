# sway-screenshot

A screenshot and screen recording utility for Sway/Wayland.

## Features

- Screenshot capture (window, screen, or selection)
- Screen recording with wf-recorder
- Clipboard integration
- Image editing with satty
- Waybar status integration
- OBS integration
- Daemon mode.

## Dependencies

**Required:**

- [grim](https://sr.ht/~emersion/grim/) - screenshot capture
- [slurp](https://github.com/emersion/slurp) - region selection
- [wf-recorder](https://github.com/ammen99/wf-recorder) - screen recording
- [wl-clipboard](https://github.com/bugaevc/wl-clipboard) - clipboard (wl-copy/wl-paste)
- [ffmpeg](https://ffmpeg.org/) - video conversion

**Optional:**

- [satty](https://github.com/gabm/satty) - screenshot annotation/editing
- [wofi](https://hg.sr.ht/~scoopta/wofi) - menu selection
- [zenity](https://gitlab.gnome.org/GNOME/zenity) - dialogs
- [nautilus](https://apps.gnome.org/Nautilus/) - file browser
- [fd](https://github.com/sharkdp/fd) - file cleanup
- [obs-cli](https://github.com/muesli/obs-cli) - OBS Studio control
- [pass](https://www.passwordstore.org/) - password store (for OBS)
- [aichat](https://github.com/sigoden/aichat) - AI-generated filenames

## Installation

```bash
go install github.com/chmouel/sway-screenshot/cmd/sway-screenshot@latest
```

Or build from source:

```bash
make build
```

## Usage

```bash
# Screenshot commands
sway-screenshot selection-clipboard
sway-screenshot selection-file
sway-screenshot selection-edit
sway-screenshot current-window-clipboard
sway-screenshot current-window-file
sway-screenshot current-screen-clipboard

# Recording commands
sway-screenshot movie-selection
sway-screenshot movie-screen
sway-screenshot movie-current-window
sway-screenshot stop-recording
sway-screenshot pause-recording
sway-screenshot toggle-record

# Waybar integration
sway-screenshot waybar-status
sway-screenshot waybar-status --follow

# OBS integration
sway-screenshot obs-toggle-recording
sway-screenshot obs-toggle-pause
```

## Waybar Configuration

```json
"custom/screenshot": {
    "exec": "sway-screenshot waybar-status --follow",
    "on-click": "sway-screenshot toggle-record -a movie-current-window",
    "return-type": "json"
}
```

## Sway Configuration

```config
```

bindsym Print exec sway-screenshot toggle-record -a movie-current-window -w 5

## Licence

Apache 2.0
