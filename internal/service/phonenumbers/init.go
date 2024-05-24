package phonenumbers

import (
	"context"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
)

//go:generate mockery --name Service
type Service interface {
	Create(ctx context.Context, cmd *CreateCmd) (*Phone, error)
	GetAllForContact(ctx context.Context, contact *contacts.Contact, cmd *sqlstorage.PaginateCmd) ([]Phone, error)
	DeleteContactPhone(ctx context.Context, contact *contacts.Contact, phoneID uuid.UUID) error
}

func Init(
	tools tools.Tools,
	db sqlstorage.Querier,
) Service {
	store := newSqlStorage(db)

	return newService(tools, store)
}
