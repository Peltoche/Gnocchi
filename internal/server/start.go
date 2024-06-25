package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Peltoche/gnocchi/assets"
	"github.com/Peltoche/gnocchi/internal/migrations"
	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/phonenumbers"
	"github.com/Peltoche/gnocchi/internal/service/utilities"
	"github.com/Peltoche/gnocchi/internal/service/vcard"
	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/logger"
	"github.com/Peltoche/gnocchi/internal/tools/router"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	contactsweb "github.com/Peltoche/gnocchi/internal/web/contacts"
	"github.com/Peltoche/gnocchi/internal/web/html"
	importsweb "github.com/Peltoche/gnocchi/internal/web/imports"
	"github.com/spf13/afero"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type Folder string

type Config struct {
	fx.Out
	Tools    tools.Config
	FS       afero.Fs
	Storage  sqlstorage.Config
	Folder   Folder
	Listener router.Config
	HTML     html.Config
	Assets   assets.Config
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(router.Registerer)),
		fx.ResultTags(`group:"routes"`),
	)
}

func start(ctx context.Context, cfg Config, invoke fx.Option) *fx.App {
	app := fx.New(
		fx.WithLogger(func(tools tools.Tools) fxevent.Logger { return logger.NewFxLogger(tools.Logger()) }),
		fx.Provide(
			func() context.Context { return ctx },
			func() Config { return cfg },

			func(folder Folder, fs afero.Fs, tools tools.Tools) (string, error) {
				folderPath, err := filepath.Abs(string(folder))
				if err != nil {
					return "", fmt.Errorf("invalid path: %q: %w", folderPath, err)
				}

				err = fs.MkdirAll(string(folder), 0o755)
				if err != nil && !errors.Is(err, os.ErrExist) {
					return "", fmt.Errorf("failed to create the %s: %w", folderPath, err)
				}

				if fs.Name() == afero.NewMemMapFs().Name() {
					tools.Logger().Info("Load data from memory")
				} else {
					tools.Logger().Info(fmt.Sprintf("Load data from %s", folder))
				}

				return folderPath, nil
			},

			// Tools
			fx.Annotate(tools.NewToolbox, fx.As(new(tools.Tools))),
			fx.Annotate(html.NewRenderer, fx.As(new(html.Writer))),
			sqlstorage.Init,

			// Services
			fx.Annotate(contacts.Init, fx.As(new(contacts.Service))),
			fx.Annotate(phonenumbers.Init, fx.As(new(phonenumbers.Service))),
			fx.Annotate(vcard.Init, fx.As(new(vcard.Service))),

			// HTTP handlers
			AsRoute(assets.NewHTTPHandler),
			AsRoute(utilities.NewHTTPHandler),

			// Web Pages
			AsRoute(contactsweb.NewListPage),
			AsRoute(contactsweb.NewDetailsPage),
			AsRoute(importsweb.NewImportsPage),

			// HTTP Router / HTTP Server
			router.InitMiddlewares,
			fx.Annotate(router.NewServer, fx.ParamTags(`group:"routes"`)),
		),

		fx.Invoke(migrations.Run),

		invoke,
	)

	return app
}
