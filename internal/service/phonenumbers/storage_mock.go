// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package phonenumbers

import (
	context "context"

	contacts "github.com/Peltoche/gnocchi/internal/service/contacts"

	mock "github.com/stretchr/testify/mock"

	sqlstorage "github.com/Peltoche/gnocchi/internal/tools/sqlstorage"

	uuid "github.com/Peltoche/gnocchi/internal/tools/uuid"
)

// mockStorage is an autogenerated mock type for the storage type
type mockStorage struct {
	mock.Mock
}

// DeletePhone provides a mock function with given fields: ctx, phone
func (_m *mockStorage) DeletePhone(ctx context.Context, phone *Phone) error {
	ret := _m.Called(ctx, phone)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Phone) error); ok {
		r0 = rf(ctx, phone)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllForContact provides a mock function with given fields: ctx, contact, cmd
func (_m *mockStorage) GetAllForContact(ctx context.Context, contact *contacts.Contact, cmd *sqlstorage.PaginateCmd) ([]Phone, error) {
	ret := _m.Called(ctx, contact, cmd)

	var r0 []Phone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *contacts.Contact, *sqlstorage.PaginateCmd) ([]Phone, error)); ok {
		return rf(ctx, contact, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *contacts.Contact, *sqlstorage.PaginateCmd) []Phone); ok {
		r0 = rf(ctx, contact, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Phone)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *contacts.Contact, *sqlstorage.PaginateCmd) error); ok {
		r1 = rf(ctx, contact, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, phoneID
func (_m *mockStorage) GetByID(ctx context.Context, phoneID uuid.UUID) (*Phone, error) {
	ret := _m.Called(ctx, phoneID)

	var r0 *Phone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*Phone, error)); ok {
		return rf(ctx, phoneID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *Phone); ok {
		r0 = rf(ctx, phoneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Phone)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, phoneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, p
func (_m *mockStorage) Save(ctx context.Context, p *Phone) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Phone) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockStorage creates a new instance of mockStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockStorage {
	mock := &mockStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
