package contacts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/halium/internal/service/contacts"
	"github.com/Peltoche/halium/internal/tools"
	"github.com/Peltoche/halium/internal/tools/errs"
	"github.com/Peltoche/halium/internal/tools/router"
	"github.com/Peltoche/halium/internal/tools/uuid"
	"github.com/Peltoche/halium/internal/web/html"
	contactstmpl "github.com/Peltoche/halium/internal/web/html/templates/contacts"
	"github.com/go-chi/chi/v5"
)

type DetailsPage struct {
	html     html.Writer
	contacts contacts.Service
	uuid     uuid.Service
}

func NewDetailsPage(tools tools.Tools, html html.Writer, contacts contacts.Service) *DetailsPage {
	return &DetailsPage{
		html:     html,
		contacts: contacts,
		uuid:     tools.UUID(),
	}
}

func (h *DetailsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/contacts/{contactID}", h.getDetails)
	r.Get("/web/contacts/{contactID}/name", h.getEditName)
	r.Post("/web/contacts/{contactID}/name", h.editName)
}

func (h *DetailsPage) getDetails(w http.ResponseWriter, r *http.Request) {
	id, err := h.uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	contact, err := h.contacts.GetByID(r.Context(), id)
	if errors.Is(err, errs.ErrNotFound) {
		http.Redirect(w, r, "/web/contacts", http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get a contact by id: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.DetailsPageTmpl{
		Contact: contact,
	})
}

func (h *DetailsPage) getEditName(w http.ResponseWriter, r *http.Request) {
	id, err := h.uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	contact, err := h.contacts.GetByID(r.Context(), id)
	if errors.Is(err, errs.ErrNotFound) {
		http.Redirect(w, r, "/web/contacts", http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get a contact by id: %w", err))
		return
	}
	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ModalEditNameTmpl{
		Contact: contact,
	})
}

func (h *DetailsPage) editName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	contact, err := h.contacts.GetByID(ctx, id)
	if errors.Is(err, errs.ErrNotFound) {
		http.Redirect(w, r, "/web/contacts", http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get a contact by id: %w", err))
		return
	}

	contact, err = h.contacts.EditName(ctx, &contacts.EditNameCmd{
		Contact:    contact,
		Prefix:     r.FormValue("namePrefix"),
		FirstName:  r.FormValue("firstName"),
		MiddleName: r.FormValue("middleName"),
		Surname:    r.FormValue("surname"),
		Suffix:     r.FormValue("nameSuffix"),
	})
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to edit the contact name: %w", err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.DetailsPageTmpl{
		Contact: contact,
	})
}
