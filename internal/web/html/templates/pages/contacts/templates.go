package contacts

import (
	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/phonenumbers"
	"github.com/Peltoche/gnocchi/internal/web/phones"
)

type ListPageTmpl struct {
	Contacts []contacts.Contact
}

func (t *ListPageTmpl) Template() string { return "pages/contacts/page_list" }

type DetailsPageTmpl struct {
	Contact *contacts.Contact
	Phones  []phonenumbers.Phone
}

func (t *DetailsPageTmpl) Template() string { return "pages/contacts/page_details" }

type ModalEditNameTmpl struct {
	Contact *contacts.Contact
}

func (t *ModalEditNameTmpl) Template() string { return "pages/contacts/modal_edit_name" }

type ModalRegisterPhoneNumberTmpl struct {
	Error    error
	Input    string
	Selected phones.PhoneData
	Contact  *contacts.Contact
	Phones   []phones.PhoneData
}

func (t *ModalRegisterPhoneNumberTmpl) Template() string {
	return "pages/contacts/modal_register_phonenumber"
}

type ModalImportsTmpl struct {
	ErrorMsg string
}

func (t *ModalImportsTmpl) Template() string {
	return "pages/contacts/modal_imports"
}
