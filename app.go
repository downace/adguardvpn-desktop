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
	mToggleWindow      *systray.MenuItem
	mConnect           *systray.MenuItem
	mLocations         *systray.MenuItem
	mLocationsSubItems []*systray.MenuItem
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
	sudoAskpassFile       *os.File
}

const trayMenuMaxLocations = 5

func cwd() string {
	executablePath, _ := os.Executable()
	return filepath.Dir(executablePath)
}

func NewApp() *App {
	app := App{
		config: makeConfig(filepath.Join(cwd(), "config.yaml")),
	}

	app.adGuardCli = adguard.Cli{
		OnStatusChange:    app.handleStatusChange,
		OnLocationsLoaded: app.handleLocationsLoaded,
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
		a.trayMenuItems.mLocations = systray.AddMenuItem("Locations", "Locations")
		a.trayMenuItems.mLocations.Hide()
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

func (a *App) initSudoAskpassCommand() (string, error) {
	file, err := os.CreateTemp(os.TempDir(), "sudo-askpass-*.sh")

	if err != nil {
		return "", err
	}

	bin, _ := os.Executable()

	_, err = file.WriteString(fmt.Sprintf(`#!/bin/sh
%s --sudo-askpass
`, bin))
	if err != nil {
		return "", err
	}
	err = file.Chmod(0o700)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	a.sudoAskpassFile = file

	return file.Name(), nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	sudoAskpassCommand, err := a.initSudoAskpassCommand()

	if err != nil {
		a.shouldShutdownOnClose = true
		runtime.WindowHide(ctx)
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: fmt.Sprintf("Failed to initialize sudo prompt:\n%s", err),
		})
		runtime.Quit(ctx)
		return
	}

	sigintCh := make(chan os.Signal, 1)
	signal.Notify(sigintCh, os.Interrupt)
	go func() {
		for range sigintCh {
			a.shouldShutdownOnClose = true
		}
	}()

	a.trayStart, a.trayEnd = systray.RunWithExternalLoop(a.initTray(), nil)
	a.trayStart()

	err = a.config.load()

	if err != nil {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.WarningDialog,
			Title:   "Warning",
			Message: fmt.Sprintf("Unable to load settings, using defaults\n%s", err),
		})
	}

	a.adGuardCli.CliBin = a.config.AdGuardBin
	a.adGuardCli.SudoAskpassCommand = sudoAskpassCommand

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
	if a.trayEnd != nil {
		a.trayEnd()
	}
	if a.sudoAskpassFile != nil {
		_ = os.Remove(a.sudoAskpassFile.Name())
	}
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

func (a *App) handleLocationsLoaded(locations []adguard.Location) {
	for _, item := range a.trayMenuItems.mLocationsSubItems {
		item.Remove()
	}
	a.trayMenuItems.mLocationsSubItems = make([]*systray.MenuItem, 0)

	if len(locations) > 0 {
		a.trayMenuItems.mLocations.Show()
		for _, location := range locations[:min(len(locations), trayMenuMaxLocations)] {
			label := fmt.Sprintf("%s, %s (%dms)", location.City, location.Country, location.Ping)
			item := a.trayMenuItems.mLocations.AddSubMenuItem(label, label)
			go func() {
				for range item.ClickedCh {
					_ = a.adGuardCli.Connect(location.City)
				}
			}()
			a.trayMenuItems.mLocationsSubItems = append(a.trayMenuItems.mLocationsSubItems, item)
		}
	} else {
		a.trayMenuItems.mLocations.Hide()
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
