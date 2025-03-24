package main

import (
	"context"
	_ "embed"
	"fmt"
	"fyne.io/systray"
	"github.com/downace/adguardvpn-desktop/internal/adguard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
)

//go:embed build/trayicon_default.png
var trayIconDefault []byte

//go:embed build/trayicon_connecting.png
var trayIconConnecting []byte

//go:embed build/trayicon_connected.png
var trayIconConnected []byte

//go:embed build/trayicon_disconnected.png
var trayIconDisconnected []byte

type AppTrayMenuItems struct {
	mToggleWindow *systray.MenuItem
	mConnect      *systray.MenuItem
}

type App struct {
	ctx                   context.Context
	config                *AppConfig
	hidden                bool
	trayStart             func()
	trayEnd               func()
	shouldShutdownOnClose bool
	trayMenuItems         AppTrayMenuItems
	adGuardCli            adguard.Cli
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
	app.adGuardCli = adguard.Cli{
		OnStatusChange: app.handleStatusChange,
	}

	return &app
}

func (a *App) setWindowHidden(hidden bool) {
	a.hidden = hidden
	if hidden {
		runtime.WindowHide(a.ctx)
		a.trayMenuItems.mToggleWindow.SetTitle("Show Window")
	} else {
		runtime.WindowShow(a.ctx)
		a.trayMenuItems.mToggleWindow.SetTitle("Hide Window")
	}
}

func (a *App) initTray() func() {
	return func() {
		systray.SetIcon(trayIconDefault)
		systray.SetTitle("AdGuard VPN")

		a.trayMenuItems.mToggleWindow = systray.AddMenuItem("Hide Window", "Hide Window")
		systray.AddSeparator()
		a.trayMenuItems.mConnect = systray.AddMenuItem("Connect", "Connect")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit")

		go func() {
			for range a.trayMenuItems.mToggleWindow.ClickedCh {
				a.setWindowHidden(!a.hidden)
			}
		}()

		go func() {
			for range a.trayMenuItems.mConnect.ClickedCh {
				_ = a.adGuardCli.ToggleConnection()
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
	_ = a.adGuardCli.RefreshStatus()
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
	return a.adGuardCli.GetStatus()
}

func (a *App) GetAdGuardAccount() (*adguard.Account, error) {
	return a.adGuardCli.Account()
}

func (a *App) AdGuardGetLocations() ([]adguard.Location, error) {
	return a.adGuardCli.GetLocations()
}

func (a *App) AdGuardConnect(location string) error {
	return a.adGuardCli.Connect(location)
}

func (a *App) AdGuardDisconnect() error {
	return a.adGuardCli.Disconnect()
}

func (a *App) handleStatusChange(status *adguard.Status) {
	runtime.EventsEmit(a.ctx, "status-changed", status)

	if status.Connecting {
		systray.SetIcon(trayIconConnecting)
		systray.SetTitle("AdGuard VPN - Connecting...")
		a.trayMenuItems.mConnect.SetTitle("Connecting...")
		a.trayMenuItems.mConnect.Disable()
	} else if status.Connected {
		systray.SetIcon(trayIconConnected)
		systray.SetTitle("AdGuard VPN - Connected to " + status.Location.City)
		a.trayMenuItems.mConnect.SetTitle("Disconnect")
		a.trayMenuItems.mConnect.Enable()
	} else {
		systray.SetIcon(trayIconDisconnected)
		systray.SetTitle("AdGuard VPN - Disconnected")
		a.trayMenuItems.mConnect.SetTitle("Connect")
		a.trayMenuItems.mConnect.Enable()
	}
}

func (a *App) GetFavoriteLocations() []string {
	return slices.AppendSeq([]string{}, maps.Values(a.config.FavoriteLocations))
}

func (a *App) AddFavoriteLocation(location string) error {
	a.config.FavoriteLocations[location] = location
	return a.config.save()
}

func (a *App) RemoveFavoriteLocation(location string) error {
	delete(a.config.FavoriteLocations, location)
	return a.config.save()
}

func (a *App) AdGuardGetExclusionMode() (adguard.ExclusionMode, error) {
	return a.adGuardCli.GetExclusionMode()
}

func (a *App) AdGuardSetExclusionMode(mode adguard.ExclusionMode) error {
	return a.adGuardCli.SetExclusionMode(mode)
}

func (a *App) AdGuardExclusionsShow() ([]string, error) {
	return a.adGuardCli.ExclusionsShow()
}

func (a *App) AdGuardExclusionsAdd(exclusions []string) error {
	return a.adGuardCli.ExclusionsAdd(exclusions)
}

func (a *App) AdGuardExclusionsRemove(exclusion string) error {
	return a.adGuardCli.ExclusionsRemove(exclusion)
}
