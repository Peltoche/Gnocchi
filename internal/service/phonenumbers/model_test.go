package phonenumbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Phone_model_getters(t *testing.T) {
	phone := NewFakePhone(t).Build()

	assert.Equal(t, phone.id, phone.ID())
	assert.Equal(t, phone.phoneType, phone.Type())
	assert.Equal(t, phone.internationalFormatted, phone.InternationalFormatted())
	assert.Equal(t, phone.nationalFormatted, phone.NationalFormatted())
	assert.Equal(t, phone.iso2RegionCode, phone.ISO2RegionCode())
	assert.Equal(t, phone.contactID, phone.ContactID())
	assert.Equal(t, phone.createdAt, phone.CreatedAt())
}
