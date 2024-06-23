package contacts

import (
	"context"
	"fmt"
	"testing"
	"time"

	stdlanguage "golang.org/x/text/language"

	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/errs"
	"github.com/Peltoche/gnocchi/internal/tools/language"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Contacts_Service(t *testing.T) {
	t.Parallel()

	t.Run("Create success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		// Data
		now := time.Now()
		newContact := Contact{
			id:        uuid.UUID("some-user-id"),
			name:      &Name{},
			createdAt: now,
		}

		// Mocks
		tools.UUIDMock.On("New").Return(uuid.UUID("some-user-id")).Once()
		tools.ClockMock.On("Now").Return(now).Once()

		storage.On("Save", ctx, &newContact).Return(nil)

		// Run
		res, err := service.Create(ctx, &CreateCmd{})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, &newContact, res)
	})

	t.Run("Create with a storage error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		// Data
		now := time.Now()
		newContact := Contact{
			id:        uuid.UUID("some-user-id"),
			name:      &Name{},
			createdAt: now,
		}

		// Mocks
		tools.UUIDMock.On("New").Return(uuid.UUID("some-user-id")).Once()
		tools.ClockMock.On("Now").Return(now).Once()

		storage.On("Save", ctx, &newContact).Return(fmt.Errorf("some-error"))

		// Run
		res, err := service.Create(ctx, &CreateCmd{})

		// Asserts
		require.ErrorContains(t, err, "some-error")
		require.ErrorIs(t, err, errs.ErrInternal)
		assert.Nil(t, res)
	})

	t.Run("GetAll success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = language.SetBrowserLangFromReq(ctx, stdlanguage.English)

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		contacts := []Contact{
			*NewFakeContact(t).Build(),
			*NewFakeContact(t).Build(),
			*NewFakeContact(t).Build(),
			*NewFakeContact(t).Build(),
		}

		storage.On("GetAll", ctx, (*sqlstorage.PaginateCmd)(nil)).Return(contacts, nil)

		// Run
		res, err := service.GetAll(ctx)

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, contacts, res)
	})

	t.Run("GetAll with a storage error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ctx = language.SetBrowserLangFromReq(ctx, stdlanguage.English)

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		storage.On("GetAll", ctx, (*sqlstorage.PaginateCmd)(nil)).Return(nil, fmt.Errorf("some-error"))

		// Run
		res, err := service.GetAll(ctx)

		// Asserts
		require.ErrorContains(t, err, "some-error")
		require.ErrorIs(t, err, errs.ErrInternal)
		assert.Nil(t, res)
	})

	t.Run("GetByID success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		contact := NewFakeContact(t).Build()

		storage.On("GetByID", ctx, contact.id).Return(contact, nil)

		// Run
		res, err := service.GetByID(ctx, contact.id)

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, contact, res)
	})

	t.Run("GetByID with a storage error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		contact := NewFakeContact(t).Build()

		storage.On("GetByID", ctx, contact.id).Return(nil, fmt.Errorf("some-error"))

		// Run
		res, err := service.GetByID(ctx, contact.id)

		// Asserts
		require.ErrorContains(t, err, "some-error")
		require.ErrorIs(t, err, errs.ErrInternal)
		assert.Nil(t, res)
	})

	t.Run("Delete success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		contact := NewFakeContact(t).Build()

		storage.On("Delete", ctx, contact.id).Return(nil)

		// Run
		err := service.Delete(ctx, contact)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("Delete with a storage error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		contact := NewFakeContact(t).Build()

		storage.On("Delete", ctx, contact.id).Return(fmt.Errorf("some-error"))

		// Run
		err := service.Delete(ctx, contact)

		// Asserts
		require.ErrorContains(t, err, "some-error")
		require.ErrorIs(t, err, errs.ErrInternal)
	})

	t.Run("EditName success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		// Edit the name of "contact" with the name of "contact2"
		contact := NewFakeContact(t).Build()
		contact2 := NewFakeContact(t).Build()

		storage.On("Patch", ctx, contact, map[string]any{
			"name_prefix": contact2.name.prefix,
			"first_name":  contact2.name.firstName,
			"middle_name": contact2.name.middleName,
			"surname":     contact2.name.surname,
			"name_suffix": contact2.name.suffix,
		}).Return(nil)

		// Run
		res, err := service.EditName(ctx, &EditNameCmd{
			Contact:    contact,
			Prefix:     contact2.name.prefix,
			FirstName:  contact2.name.firstName,
			MiddleName: contact2.name.middleName,
			Surname:    contact2.name.surname,
			Suffix:     contact2.name.suffix,
		})

		// Asserts
		require.NoError(t, err)
		assert.NotEqual(t, contact, res)
		assert.NotEqual(t, contact2, res)
		assert.Equal(t, contact2.name, res.name)
	})

	t.Run("EditName with a storage error", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		service := newService(tools, storage)

		// Edit the name of "contact" with the name of "contact2"
		contact := NewFakeContact(t).Build()
		contact2 := NewFakeContact(t).Build()

		storage.On("Patch", ctx, contact, map[string]any{
			"name_prefix": contact2.name.prefix,
			"first_name":  contact2.name.firstName,
			"middle_name": contact2.name.middleName,
			"surname":     contact2.name.surname,
			"name_suffix": contact2.name.suffix,
		}).Return(fmt.Errorf("some-error"))

		// Run
		res, err := service.EditName(ctx, &EditNameCmd{
			Contact:    contact,
			Prefix:     contact2.name.prefix,
			FirstName:  contact2.name.firstName,
			MiddleName: contact2.name.middleName,
			Surname:    contact2.name.surname,
			Suffix:     contact2.name.suffix,
		})

		// Asserts
		require.ErrorContains(t, err, "some-error")
		require.ErrorIs(t, err, errs.ErrInternal)
		assert.Nil(t, res)
	})
}
