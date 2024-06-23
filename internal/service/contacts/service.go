package contacts

import (
	"context"
	"errors"
	"fmt"

	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/clock"
	"github.com/Peltoche/gnocchi/internal/tools/errs"
	"github.com/Peltoche/gnocchi/internal/tools/language"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"golang.org/x/text/collate"
)

//go:generate mockery --name storage
type storage interface {
	Save(ctx context.Context, c *Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contact, error)
	Patch(ctx context.Context, contact *Contact, fields map[string]any) error
	GetAll(ctx context.Context, cmd *sqlstorage.PaginateCmd) ([]Contact, error)
	Delete(ctx context.Context, contactID uuid.UUID) error
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

func (s *service) GetAll(ctx context.Context) ([]Contact, error) {
	res, err := s.storage.GetAll(ctx, nil)
	if err != nil {
		return nil, errs.Internal(err)
	}

	browserLang := language.GetBrowserLangFromReq(ctx)

	collator := collate.New(browserLang, collate.IgnoreCase, collate.IgnoreDiacritics)

	contactList := &contactList{res}

	collator.Sort(contactList)

	return contactList.contacts, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Contact, error) {
	res, err := s.storage.GetByID(ctx, id)
	if errors.Is(err, errNotFound) {
		return nil, errs.NotFound(err)
	}

	if err != nil {
		return nil, errs.Internal(err)
	}

	return res, nil
}

func (s *service) Delete(ctx context.Context, contact *Contact) error {
	err := s.storage.Delete(ctx, contact.id)
	if err != nil {
		return errs.Internal(err)
	}

	return nil
}

// Create will create and register a new user.
func (s *service) Create(ctx context.Context, cmd *CreateCmd) (*Contact, error) {
	contact := Contact{
		id:        s.uuid.New(),
		name:      &Name{},
		createdAt: s.clock.Now(),
	}

	err := s.storage.Save(ctx, &contact)
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("storage error: %w", err))
	}

	return &contact, nil
}

func (s *service) EditName(ctx context.Context, cmd *EditNameCmd) (*Contact, error) {
	newName := Name{
		prefix:     cmd.Prefix,
		firstName:  cmd.FirstName,
		middleName: cmd.MiddleName,
		surname:    cmd.Surname,
		suffix:     cmd.Suffix,
	}

	err := s.storage.Patch(ctx, cmd.Contact, map[string]any{
		"name_prefix": newName.prefix,
		"first_name":  newName.firstName,
		"middle_name": newName.middleName,
		"surname":     newName.surname,
		"name_suffix": newName.suffix,
	})
	if err != nil {
		return nil, errs.Internal(fmt.Errorf("failed to patch the contact name: %w", err))
	}

	updatedContact := *cmd.Contact
	updatedContact.name = &newName

	return &updatedContact, nil
}

type contactList struct {
	contacts []Contact
}

func (l *contactList) Len() int {
	return len(l.contacts)
}

func (l *contactList) Swap(i, j int) {
	l.contacts[i], l.contacts[j] = l.contacts[j], l.contacts[i]
}

func (l *contactList) Bytes(i int) []byte {
	return []byte(l.contacts[i].name.DisplayName())
}
