package main

import (
	"context"
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var iconData []byte

//go:embed build/windows/icon.ico
var iconIco []byte

func main() {
	// Resolve log path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		println("Error getting home directory:", err.Error())
		return
	}
	logPath := filepath.Join(homeDir, ".custos", "logs")
	if err := os.MkdirAll(logPath, 0755); err != nil {
		println("Error creating log directory:", err.Error())
		return
	}

	// Single Instance Lock
	lockFile := filepath.Join(homeDir, ".custos", "app.lock")
	if _, err := os.Stat(lockFile); err == nil {
		// On Windows, if the app is running, removing will fail or be locked.
		// We try to remove it; if it fails, another instance is likely using it.
		err = os.Remove(lockFile)
		if err != nil {
			println("Another instance of Custos is already running.")
			return
		}
	}
	// Create the lock file
	f, err := os.Create(lockFile)
	if err != nil {
		println("Error creating lock file:", err.Error())
		return
	}
	f.Close()
	defer os.Remove(lockFile)

	// Create an instance of the app structure
	app := NewApp()

	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu()) // On macOS platform, this must be done right after `NewMenu()`
	}
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Reload", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		rt.WindowReloadApp(app.ctx)
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		rt.Quit(app.ctx)
	})

	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // On macOS platform, EditMenu should be appended to enable Cmd+C, Cmd+V, Cmd+Z... shortcuts
	}
	AboutMenu := AppMenu.AddSubmenu("Help")
	AboutMenu.AddText("About", keys.CmdOrCtrl("a"), func(_ *menu.CallbackData) {
		rt.EventsEmit(app.ctx, "navigate-to", "/about")
	})
	AboutMenu.AddText("Open Source", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		rt.EventsEmit(app.ctx, "navigate-to", "/opensource")
	})
	AboutMenu.AddText("Reset Data", keys.Combo("c", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
		app.store.ResetData()
		rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
			Type:          rt.InfoDialog,
			Title:         "Reset Data",
			Message:       "Data reset successfully",
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
	})

	// Create application with options
	err = wails.Run(&options.App{
		Title:             strings.ToTitle(app.GetAppInfo().Name[:1]) + strings.ToLower(app.GetAppInfo().Name[1:]),
		Width:             1024,
		Height:            768,
		AlwaysOnTop:       false,
		HideWindowOnClose: true,
		Windows:           &windows.Options{},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			go setupTray(app)
		},
		OnShutdown: func(ctx context.Context) {
			app.shutdown(ctx)
			systray.Quit()
		},
		Menu: AppMenu,
		Linux: &linux.Options{
			Icon:             iconData,
			ProgramName:      strings.ToLower(app.GetAppInfo().Name),
			WebviewGpuPolicy: linux.WebviewGpuPolicyAlways,
		},
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.WARNING,
		Logger:             logger.NewFileLogger(filepath.Join(logPath, "app.log")),
		StartHidden:        false,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func setupTray(app *App) {
	runtime.LockOSThread()
	systray.Run(func() {
		if runtime.GOOS == "windows" {
			systray.SetIcon(iconIco)
		} else {
			systray.SetIcon(iconData)
		}
		systray.SetTitle(app.GetAppInfo().Name)
		systray.SetTooltip(app.GetAppInfo().Name)

		mShow := systray.AddMenuItem("Show", "Show the main window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the application")

		for {
			select {
			case <-mShow.ClickedCh:
				rt.WindowUnminimise(app.ctx)
				rt.WindowShow(app.ctx)
			case <-mQuit.ClickedCh:
				systray.Quit()
				rt.Quit(app.ctx)
				return
			}
		}
	}, func() {
		// OnExit
	})
}
