package vcard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/phonenumbers"
	"github.com/emersion/go-vcard"
)

var (
	ErrUnsupportedVCardVersion = fmt.Errorf("unsupported vcard version")
	ErrInvalidVCard            = fmt.Errorf("invalid vcard file")
)

type service struct {
	contacts     contacts.Service
	phonenumbers phonenumbers.Service
}

func newService(contacts contacts.Service, phonenumbers phonenumbers.Service) *service {
	return &service{
		contacts:     contacts,
		phonenumbers: phonenumbers,
	}
}

func (s *service) ImportVCardFile(ctx context.Context, file io.Reader) error {
	dec := vcard.NewDecoder(file)

	for {
		card, err := dec.Decode()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidVCard, err)
		}

		if card.Kind() != vcard.KindIndividual {
			continue
		}

		if !strings.HasPrefix(card.Value("VERSION"), "3.") {
			return ErrUnsupportedVCardVersion
		}

		contact, err := s.contacts.Create(ctx, &contacts.CreateCmd{})
		if err != nil {
			return fmt.Errorf("failed to create a contact: %w", err)
		}

		name := card.Name()

		_, err = s.contacts.EditName(ctx, &contacts.EditNameCmd{
			Contact:    contact,
			Prefix:     name.HonorificPrefix,
			MiddleName: name.AdditionalName,
			FirstName:  name.GivenName,
			Surname:    name.FamilyName,
			Suffix:     name.HonorificSuffix,
		})
		if err != nil {
			return fmt.Errorf("failed to edit a contact name: %w", err)
		}

		tels, ok := card["TEL"]
		if !ok {
			continue
		}

		for _, tel := range tels {
			_, err = s.phonenumbers.Create(ctx, &phonenumbers.CreateCmd{
				Contact: contact,
				Type:    tel.Params.Get("TYPE"),
				Region:  "FR",
				Input:   tel.Value,
			})
			if err != nil {
				return fmt.Errorf("failed to register a phonenumber: %w", err)
			}
		}
	}

	return nil
}
