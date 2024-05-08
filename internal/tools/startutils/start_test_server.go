package startutils

import (
	"testing"

	"github.com/Peltoche/halium/internal/tools"
	"github.com/Peltoche/halium/internal/tools/sqlstorage"
	"github.com/spf13/afero"
)

type Server struct {
	// Main tools
	Tools *tools.Toolbox
	DB    sqlstorage.Querier
	FS    afero.Fs

	// Services
}

func NewServer(t *testing.T) *Server {
	t.Helper()

	tools := tools.NewToolboxForTest(t)
	db := sqlstorage.NewTestStorage(t)
	afs := afero.NewMemMapFs()

	return &Server{
		Tools: tools,
		DB:    db,
		FS:    afs,
	}
}
