package tools

import (
	"log/slog"

	"github.com/Peltoche/halium/internal/tools/clock"
	"github.com/Peltoche/halium/internal/tools/password"
	"github.com/Peltoche/halium/internal/tools/response"
	"github.com/Peltoche/halium/internal/tools/uuid"
)

// Tools regroup all the utilities required for a working server.
type Tools interface {
	Clock() clock.Clock
	UUID() uuid.Service
	Logger() *slog.Logger
	ResWriter() response.Writer
	Password() password.Password
}
