package main

import (
	"context"
	_ "embed"
	"fmt"
	"fyne.io/systray"
	"github.com/downace/adguardvpn-desktop/internal/adguard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	"path/filepath"
)

//go:embed build/trayicon_connected.png
var trayIconConnected []byte

//go:embed build/trayicon_disconnected.png
var trayIconDisconnected []byte

type App struct {
	ctx                   context.Context
	config                *AppConfig
	hidden                bool
	trayStart             func()
	trayEnd               func()
	shouldShutdownOnClose bool

	adGuardCli adguard.Cli
}

func cwd() string {
	executablePath, _ := os.Executable()
	return filepath.Dir(executablePath)
}

func NewApp() *App {
	app := App{
		config: makeConfig(filepath.Join(cwd(), "config.yaml")),
	}

	app.trayStart, app.trayEnd = systray.RunWithExternalLoop(app.initTray(), nil)
	app.adGuardCli = adguard.Cli{}

	return &app
}

func (a *App) setWindowHidden(hidden bool) {
	a.hidden = hidden
}

func (a *App) initTray() func() {
	return func() {
		systray.SetIcon(trayIconDisconnected)
		systray.SetTitle("AdGuard VPN")
		mToggle := systray.AddMenuItem("Toggle Window", "Toggle Window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit")

		go func() {
			for range mToggle.ClickedCh {
				if a.hidden {
					runtime.WindowShow(a.ctx)
				} else {
					runtime.WindowHide(a.ctx)
				}
				a.setWindowHidden(!a.hidden)
			}
		}()

		go func() {
			for range mQuit.ClickedCh {
				a.shouldShutdownOnClose = true
				runtime.Quit(a.ctx)
			}
		}()
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	sigintCh := make(chan os.Signal, 1)
	signal.Notify(sigintCh, os.Interrupt)
	go func() {
		for range sigintCh {
			a.shouldShutdownOnClose = true
		}
	}()

	a.trayStart()

	err := a.config.load()

	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.WarningDialog,
			Title:   "Warning",
			Message: fmt.Sprintf("Unable to load settings, using defaults\n%s", err),
		})
	}

	a.adGuardCli.CliBin = a.config.AdGuardBin
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	if a.shouldShutdownOnClose {
		return false
	}
	runtime.WindowHide(ctx)
	a.setWindowHidden(true)
	return true
}

func (a *App) shutdown(_ context.Context) {
	a.trayEnd()
}

func (a *App) PickFilePath() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{})
}

func (a *App) GetAdGuardBin() string {
	return a.config.AdGuardBin
}

func (a *App) UpdateAdGuardBin(bin string) (string, error) {
	oldBin := a.adGuardCli.CliBin

	a.adGuardCli.CliBin = bin

	out, err := a.adGuardCli.Version()

	if err != nil {
		a.adGuardCli.CliBin = oldBin
		return "", err
	}

	a.config.AdGuardBin = a.adGuardCli.CliBin
	err = a.config.save()

	if err != nil {
		a.config.AdGuardBin = oldBin
		a.adGuardCli.CliBin = oldBin
		return "", err
	}

	return out, err
}

func (a *App) GetAdGuardVersion() (string, error) {
	return a.adGuardCli.Version()
}

func (a *App) GetAdGuardStatus() (*adguard.Status, error) {
	return a.adGuardCli.Status()
}
