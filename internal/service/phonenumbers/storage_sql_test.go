package phonenumbers

import (
	"context"
	"testing"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhoneSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)

	contact := contacts.NewFakeContact(t).BuildAndStore(ctx, db)
	phone := NewFakePhone(t).WithContact(contact).Build()

	t.Run("Create success", func(t *testing.T) {
		// Run
		err := store.Save(ctx, phone)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetAllForContact success", func(t *testing.T) {
		// Run
		res, err := store.GetAllForContact(ctx, contact, &sqlstorage.PaginateCmd{})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, []Phone{*phone}, res)
	})
}
