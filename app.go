package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/downace/adguardvpn-desktop/internal/adguard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
)

type App struct {
	ctx    context.Context
	config *AppConfig

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

	app.adGuardCli = adguard.Cli{}

	return &app
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

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
