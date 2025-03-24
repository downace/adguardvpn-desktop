package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	// https://github.com/wailsapp/wails/issues/2977
	_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")

	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "AdGuard VPN",
		Width:         800,
		Height:        600,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:     app.startup,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,
		Bind: []interface{}{
			app,
		},
		EnumBind: []interface{}{
			allExclusionModes,
		},
		Linux: &linux.Options{
			Icon: appIcon,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
