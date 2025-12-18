package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

	// Create an instance of the app structure
	app := NewApp()

	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu()) // On macOS platform, this must be done right after `NewMenu()`
	}
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		// do something
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		// `rt` is an alias of "github.com/wailsapp/wails/v2/pkg/runtime" to prevent collision with standard package
		rt.Quit(app.ctx)
	})

	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // On macOS platform, EditMenu should be appended to enable Cmd+C, Cmd+V, Cmd+Z... shortcuts
	}
	AboutMenu := AppMenu.AddSubmenu("Help")
	AboutMenu.AddText("About", keys.CmdOrCtrl("a"), func(_ *menu.CallbackData) {
		info := app.GetAppInfo()
		rt.MessageDialog(app.ctx, rt.MessageDialogOptions{
			Type:          rt.InfoDialog,
			Title:         "About " + info.Name,
			Message:       fmt.Sprintf("%s v%s\n\n%s", info.Name, info.Version, info.Description),
			Buttons:       []string{"OK"},
			DefaultButton: "OK",
		})
	})
	// Create application with options
	err = wails.Run(&options.App{
		Title:             "custos",
		Width:             1024,
		Height:            768,
		AlwaysOnTop:       false,
		HideWindowOnClose: true,
		Windows:           &windows.Options{},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Menu:             AppMenu,
		Linux: &linux.Options{
			Icon:             []byte("frontend/src/assets/images/logo-universal.png"),
			ProgramName:      "custos",
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
