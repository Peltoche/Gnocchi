package phonenumbers

import (
	"time"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
)

type Phone struct {
	createdAt              time.Time
	id                     uuid.UUID
	phoneType              string
	internationalFormatted string
	nationalFormatted      string
	normalized             string
	contactID              uuid.UUID
	iso2RegionCode         string
}

func (p Phone) ID() uuid.UUID                  { return p.id }
func (p Phone) Type() string                   { return p.phoneType }
func (p Phone) InternationalFormatted() string { return p.internationalFormatted }
func (p Phone) NationalFormatted() string      { return p.nationalFormatted }
func (p Phone) ISO2RegionCode() string         { return p.iso2RegionCode }
func (p Phone) ContactID() uuid.UUID           { return p.contactID }
func (p Phone) CreatedAt() time.Time           { return p.createdAt }

type CreateCmd struct {
	Contact *contacts.Contact
	Type    string
	Region  string
	Input   string
}
