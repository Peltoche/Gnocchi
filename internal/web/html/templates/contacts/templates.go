package contacts

import "github.com/Peltoche/gnocchi/internal/service/contacts"

type ListPageTmpl struct {
	Contacts []contacts.Contact
}

func (t *ListPageTmpl) Template() string { return "contacts/page_list" }

type DetailsPageTmpl struct {
	Contact *contacts.Contact
}

func (t *DetailsPageTmpl) Template() string { return "contacts/page_details" }

type ModalEditNameTmpl struct {
	Contact *contacts.Contact
}

func (t *ModalEditNameTmpl) Template() string { return "contacts/modal_edit_name" }
