package phonenumbers

import (
	"testing"
	"time"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/nyaruka/phonenumbers"
	"github.com/stretchr/testify/require"
)

type FakePhoneBuilder struct {
	t     testing.TB
	phone *Phone
}

func NewFakePhone(t testing.TB) *FakePhoneBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	iso2RegionCode := "US"

	num, err := phonenumbers.Parse(gofakeit.Phone(), iso2RegionCode)
	require.NoError(t, err)

	return &FakePhoneBuilder{
		t: t,
		phone: &Phone{
			id:                     uuidProvider.New(),
			phoneType:              "Home",
			iso2RegionCode:         iso2RegionCode,
			internationalFormatted: phonenumbers.Format(num, phonenumbers.INTERNATIONAL),
			nationalFormatted:      phonenumbers.Format(num, phonenumbers.NATIONAL),
			contactID:              uuidProvider.New(),
			createdAt:              createdAt,
		},
	}
}

func (f *FakePhoneBuilder) WithContact(contact *contacts.Contact) *FakePhoneBuilder {
	f.phone.contactID = contact.ID()

	return f
}

func (f *FakePhoneBuilder) Build() *Phone {
	return f.phone
}
