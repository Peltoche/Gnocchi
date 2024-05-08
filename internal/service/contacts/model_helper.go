package contacts

import (
	"testing"
	"time"

	"github.com/Peltoche/halium/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
)

type FakeContactBuilder struct {
	t       testing.TB
	contact *Contact
}

func NewFakeContact(t testing.TB) *FakeContactBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakeContactBuilder{
		t: t,
		contact: &Contact{
			id: uuidProvider.New(),
			name: &Name{
				prefix:     gofakeit.NamePrefix(),
				firstName:  gofakeit.FirstName(),
				middleName: gofakeit.MiddleName(),
				surname:    gofakeit.LastName(),
				suffix:     gofakeit.NameSuffix(),
			},
			createdAt: createdAt,
		},
	}
}

func (f *FakeContactBuilder) Build() *Contact {
	return f.contact
}
