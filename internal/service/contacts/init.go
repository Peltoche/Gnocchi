package contacts

import (
	"context"

	"github.com/Peltoche/halium/internal/tools"
	"github.com/Peltoche/halium/internal/tools/sqlstorage"
	"github.com/Peltoche/halium/internal/tools/uuid"
)

//go:generate mockery --name Service
type Service interface {
	Create(ctx context.Context, cmd *CreateCmd) (*Contact, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Contact, error)
	EditName(ctx context.Context, cmd *EditNameCmd) (*Contact, error)
	GetAll(ctx context.Context, paginateCmd *sqlstorage.PaginateCmd) ([]Contact, error)
}

func Init(
	tools tools.Tools,
	db sqlstorage.Querier,
) Service {
	store := newSqlStorage(db)

	return newService(tools, store)
}
