package contacts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactModel(t *testing.T) {
	contact := NewFakeContact(t).Build()

	t.Run("Name", func(t *testing.T) {
		name := contact.name

		assert.Equal(t, name.prefix, name.Prefix())
		assert.Equal(t, name.firstName, name.FirstName())
		assert.Equal(t, name.middleName, name.MiddleName())
		assert.Equal(t, name.surname, name.Surname())
		assert.Equal(t, name.suffix, name.Suffix())
	})

	t.Run("DisplayName", func(t *testing.T) {
		name := contact.name

		assert.Contains(t, name.DisplayName(), name.prefix)
		assert.Contains(t, name.DisplayName(), name.firstName)
		assert.Contains(t, name.DisplayName(), name.middleName)
		assert.Contains(t, name.DisplayName(), name.surname)
		assert.Contains(t, name.DisplayName(), name.suffix)
	})

	t.Run("DisplayName with no name", func(t *testing.T) {
		name := &Name{}

		assert.Equal(t, "(No name)", name.DisplayName())
	})
}
