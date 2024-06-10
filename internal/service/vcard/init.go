package vcard

import (
	"context"
	"io"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/phonenumbers"
)

//go:generate mockery --name Service
type Service interface {
	ImportVCardFile(ctx context.Context, file io.Reader) error
}

func Init(
	contacts contacts.Service,
	phonenumbers phonenumbers.Service,
) Service {
	return newService(contacts, phonenumbers)
}
