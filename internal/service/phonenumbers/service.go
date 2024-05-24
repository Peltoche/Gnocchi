package phonenumbers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/clock"
	"github.com/Peltoche/gnocchi/internal/tools/errs"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/nyaruka/phonenumbers"
)

//go:generate mockery --name storage
type storage interface {
	Save(ctx context.Context, p *Phone) error
	GetAllForContact(ctx context.Context, contact *contacts.Contact, cmd *sqlstorage.PaginateCmd) ([]Phone, error)
	GetByID(ctx context.Context, phoneID uuid.UUID) (*Phone, error)
	DeletePhone(ctx context.Context, phone *Phone) error
}

type service struct {
	storage storage
	uuid    uuid.Service
	clock   clock.Clock
}

// newService create a new user service.
func newService(tools tools.Tools, storage storage) *service {
	return &service{
		storage: storage,
		uuid:    tools.UUID(),
		clock:   tools.Clock(),
	}
}

func (s *service) Create(ctx context.Context, cmd *CreateCmd) (*Phone, error) {
	number, err := phonenumbers.Parse(cmd.Input, strings.ToUpper(cmd.Region))
	if err != nil {
		return nil, errs.Validation(fmt.Errorf("invalid phone number %q: %w", cmd.Input, err))
	}

	internationalFormated := phonenumbers.Format(number, phonenumbers.INTERNATIONAL)

	phoneType := strings.TrimSpace(cmd.Type)
	if len(phoneType) == 0 {
		return nil, errs.Validation(errors.New("invalid/missing phone type"))
	}

	phone := Phone{
		createdAt:              s.clock.Now(),
		id:                     s.uuid.New(),
		phoneType:              phoneType,
		internationalFormatted: internationalFormated,
		nationalFormatted:      phonenumbers.Format(number, phonenumbers.NATIONAL),
		normalized:             strings.NewReplacer(" ", "", "+", "").Replace(internationalFormated),
		contactID:              cmd.Contact.ID(),
		iso2RegionCode:         strings.ToUpper(cmd.Region),
	}

	err = s.storage.Save(ctx, &phone)
	if err != nil {
		return nil, fmt.Errorf("failed to save into storage: %w", err)
	}

	return &phone, nil
}

func (s *service) GetAllForContact(ctx context.Context, contact *contacts.Contact, cmd *sqlstorage.PaginateCmd) ([]Phone, error) {
	res, err := s.storage.GetAllForContact(ctx, contact, cmd)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(err)
	}

	if err != nil {
		return nil, errs.Internal(err)
	}

	return res, nil
}

func (s *service) DeleteContactPhone(ctx context.Context, contact *contacts.Contact, phoneID uuid.UUID) error {
	phone, err := s.storage.GetByID(ctx, phoneID)
	if errors.Is(err, errNotFound) {
		return nil
	}

	if err != nil {
		return errs.Internal(err)
	}

	if phone.contactID != contact.ID() {
		return errs.Unauthorized(fmt.Errorf("phone number %q is not owned by contact %q", phoneID, contact.ID()))
	}

	err = s.storage.DeletePhone(ctx, phone)
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to delete the phone number %q: %w", phoneID, err))
	}

	return nil
}
