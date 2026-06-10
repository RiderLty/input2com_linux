# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

input2com_linux is a Linux input event capture and forwarding daemon. It reads input events from evdev devices (mouse, keyboard, joystick, touchpad) and forwards them to various output backends (serial HID dongles, UDP, Unix Domain Socket, uinput, USB gadget). Includes a macro engine with recoil compensation, auto-fire, and AI-assisted aimbot.

All Go code is in `package main` (single flat package, no sub-packages). Comments and UI text are primarily in Chinese.

## Build Commands

**Full build** (frontend must be built first — it's embedded into the Go binary via `//go:embed`):
```bash
cd server && yarn install && yarn build
cd .. && go build -ldflags="-s -w" -o input2com
```

**Cross-compile**:
```bash
CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o input2com_amd64
CGO_ENABLE=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o input2com_arm64
```

**Frontend dev server** (proxies `/api` to backend):
```bash
cd server && yarn start
```

**Run**: `sudo ./input2com` (needs root for `/dev/input/*` access)

**Deployment**: `sudo bash install.sh` (installs systemd service)

## Architecture

**Data flow**: Input Adapter → MacroInterceptor → Output Controller

**Core interface** (`controller_interface.go`): `mouseKeyboard` with 5 methods — `MouseMove`, `MouseBtnDown`, `MouseBtnUp`, `KeyDown`, `KeyUp`.

**Input adapters** (`adapter_*.go`):
- `LinuxInputs_mouseKeyboard` — evdev mouse/keyboard from `/dev/input/eventN`
- `LinuxInputs_joystick` — evdev gamepad, auto-detects device type from capabilities, maps via JSON configs in `joystickInfos/`
- `LinuxInputs_touchpad` — evdev touchpad (single-finger mouse, two-finger scroll)
- `UDP` — receives binary-packed evdev events over UDP
- `MAKCU` — serial device input

**Output controllers** (`controller_*.go`):
- `KCOM5` — USB HID serial dongle (custom binary protocol, 2Mbaud)
- `MAKCU` — text-based `km.*` command protocol, 4Mbaud
- `UDP` / `UDS` — network forwarding
- `Uinput` — virtual input device via `/dev/uinput`
- `USBGadget` — direct HID report writes to `/dev/input/hidg*`

**Middleware** (`controller_MacroInterceptor.go`): Wraps any `mouseKeyboard` controller. Intercepts button events to execute macros — recoil compensation (trajectory files in `config/`), auto-fire, AI aim assist. The `aimBot.go` receives vision detection results over UDP (port 9321). `remote.go` generates smooth recoil trajectories.

**Config server** (`config_server.go`): Embedded React SPA served on port 9264. REST API for macro configuration. The web UI (`server/`) uses React 18, MUI 5, MobX 6, Webpack 5 (ejected CRA), managed with Yarn.

## Configuration

Runtime config: `config.yaml` (Viper). Key fields: `usingDst` selects output backend, `src.*` for input settings, `dst.*` for output settings, `server.port` for web UI.

Joystick mappings: `joystickInfos/*.json`. Recoil data: `config/*.txt` (space-separated dx dy per line, 10ms slices).

## Key Files

- `main.go` — entry point, config loading, source/destination wiring
- `defines.go` — HID/Linux keycode tables, ioctl helpers, key mappings
- `controller_interface.go` — `mouseKeyboard` interface definition
- `controller_MacroInterceptor.go` — macro engine core
- `aimBot.go` — AI aim assist UDP receiver (port 9321)
- `ipc.go` — UDP and UDS reader/writer utilities
- `custom_logger.go` — logger wrapper (go-log with color)
