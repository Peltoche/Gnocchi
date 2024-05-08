package contacts

import (
	"fmt"
	"net/http"
	"path"

	"github.com/Peltoche/halium/internal/service/contacts"
	"github.com/Peltoche/halium/internal/tools/router"
	"github.com/Peltoche/halium/internal/tools/sqlstorage"
	"github.com/Peltoche/halium/internal/web/html"
	contactstmpl "github.com/Peltoche/halium/internal/web/html/templates/contacts"
	"github.com/go-chi/chi/v5"
)

type ListPage struct {
	html     html.Writer
	contacts contacts.Service
}

func NewListPage(html html.Writer, contacts contacts.Service) *ListPage {
	return &ListPage{
		html:     html,
		contacts: contacts,
	}
}

func (h *ListPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/contacts", h.getList)
	r.Get("/web/contacts/more", h.getMoreContacts)
	r.Post("/web/contacts", h.createNewContact)
}

func (h *ListPage) getList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contacts, err := h.contacts.GetAll(ctx, &sqlstorage.PaginateCmd{Limit: 20})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the contact list: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ListPageTmpl{
		Contacts: contacts,
	})
}

func (h *ListPage) getMoreContacts(w http.ResponseWriter, r *http.Request) {
}

func (h *ListPage) createNewContact(w http.ResponseWriter, r *http.Request) {
	contact, err := h.contacts.Create(r.Context(), &contacts.CreateCmd{})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to create a new contact: %w", err))
		return
	}

	http.Redirect(w, r, path.Join("/web/contacts/", string(contact.ID())), http.StatusFound)
}
