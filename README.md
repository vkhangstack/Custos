# <img src="frontend/src/assets/images/logo-universal.png" width="48" height="48" alt="Custos Logo"> Custos

Custos is a desktop application that acts as a “guardian” for your local network, sitting between your devices and the internet to monitor and filter traffic.
Built with Go, Wails, and a React/TypeScript frontend, Custos provides a cross‑platform GUI for managing DNS/HTTP filtering rules and observing network activity.

## Features

- **Network Guardian**: Backend implements a local DNS and HTTP proxy that can block or redirect requests based on configurable blocklists.
- **Malicious Protection**: Helps protect your devices from malicious domains and unwanted endpoints.
- **Background Operation**: Integrates with the system tray to run in the background while providing quick access to protection status and controls.
- **Cross-Platform**: Built for Windows and Linux.

## Development

- **Wails**: Powers the desktop shell.
- **Vite**: Serves the React/TS frontend in dev mode, enabling hot‑reload.
- **Run Dev**: `wails dev` in the project directory.

## Production

- **Build**: Production builds are generated with `wails build`.
- **Packaging**: Packaged for Windows and Linux using platform‑specific scripts and configuration files such as `build.ps1`, `nfpm.yaml`, and `wails.json`.

