package contacts

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/vcard"
	"github.com/Peltoche/gnocchi/internal/tools/router"
	"github.com/Peltoche/gnocchi/internal/web/html"
	contactstmpl "github.com/Peltoche/gnocchi/internal/web/html/templates/pages/contacts"
	"github.com/go-chi/chi/v5"
)

type ListPage struct {
	html     html.Writer
	contacts contacts.Service
	vcard    vcard.Service
}

func NewListPage(html html.Writer, contacts contacts.Service, vcard vcard.Service) *ListPage {
	return &ListPage{
		html:     html,
		vcard:    vcard,
		contacts: contacts,
	}
}

func (h *ListPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/contacts", h.getList)
	r.Post("/web/contacts", h.createNewContact)
	r.Get("/web/contacts/imports", h.getImportsModal)
	r.Post("/web/contacts/imports", h.importFile)
}

func (h *ListPage) getList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contacts, err := h.contacts.GetAll(ctx)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the contact list: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ListPageTmpl{
		Contacts: contacts,
	})
}

func (h *ListPage) createNewContact(w http.ResponseWriter, r *http.Request) {
	contact, err := h.contacts.Create(r.Context(), &contacts.CreateCmd{})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to create a new contact: %w", err))
		return
	}

	http.Redirect(w, r, path.Join("/web/contacts/", string(contact.ID())), http.StatusFound)
}

func (h *ListPage) getImportsModal(w http.ResponseWriter, r *http.Request) {
	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ModalImportsTmpl{})
}

func (h *ListPage) importFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to fetch the file: %w", err))
		return
	}
	defer file.Close()

	err = h.vcard.ImportVCardFile(r.Context(), file)
	switch {
	case err == nil:
		h.getList(w, r)
	case errors.Is(err, vcard.ErrUnsupportedVCardVersion) ||
		errors.Is(err, vcard.ErrInvalidVCard):
		h.html.WriteHTMLTemplate(w, r, http.StatusUnprocessableEntity, &contactstmpl.ModalImportsTmpl{
			ErrorMsg: err.Error(),
		})
	default:
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to import the vcard file: %w", err))
	}
}
