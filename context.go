package main

import (
	"context"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/justwatchcom/gopass/pkg/backend/crypto/gpg"
	"github.com/justwatchcom/gopass/pkg/config"
	"github.com/justwatchcom/gopass/pkg/ctxutil"
	"github.com/justwatchcom/gopass/pkg/store/sub"
	"golang.org/x/crypto/ssh/terminal"
)

func initContext(ctx context.Context, cfg *config.Config) context.Context {
	// always trust
	ctx = gpg.WithAlwaysTrust(ctx, true)

	// check recipients conflicts with always trust, make sure it's not enabled
	// when always trust is
	if gpg.IsAlwaysTrust(ctx) {
		ctx = sub.WithCheckRecipients(ctx, false)
	}

	// debug flag
	if gdb := os.Getenv("GOPASS_DEBUG"); gdb == "true" {
		ctx = ctxutil.WithDebug(ctx, true)
	}

	// need this override for our integration tests
	if nc := os.Getenv("GOPASS_NOCOLOR"); nc == "true" || ctxutil.IsNoColor(ctx) {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
	}

	// only emit color codes when stdout is a terminal
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
		ctx = ctxutil.WithTerminal(ctx, false)
		ctx = ctxutil.WithInteractive(ctx, false)
	}

	// reading from stdin?
	if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeCharDevice == 0 {
		ctx = ctxutil.WithInteractive(ctx, false)
		ctx = ctxutil.WithStdin(ctx, true)
	}

	// disable colored output on windows since cmd.exe doesn't support ANSI color
	// codes. Other terminal may do, but until we can figure that out better
	// disable this for all terms on this platform
	if runtime.GOOS == "windows" {
		color.NoColor = true
		ctx = ctxutil.WithColor(ctx, false)
	}

	return ctx
}
