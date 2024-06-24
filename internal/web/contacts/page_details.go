package contacts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/gnocchi/internal/service/contacts"
	"github.com/Peltoche/gnocchi/internal/service/phonenumbers"
	"github.com/Peltoche/gnocchi/internal/tools"
	"github.com/Peltoche/gnocchi/internal/tools/errs"
	"github.com/Peltoche/gnocchi/internal/tools/logger"
	"github.com/Peltoche/gnocchi/internal/tools/router"
	"github.com/Peltoche/gnocchi/internal/tools/uuid"
	"github.com/Peltoche/gnocchi/internal/web/html"
	contactstmpl "github.com/Peltoche/gnocchi/internal/web/html/templates/pages/contacts"
	"github.com/Peltoche/gnocchi/internal/web/phones"
	"github.com/go-chi/chi/v5"
)

type DetailsPage struct {
	html         html.Writer
	contacts     contacts.Service
	phonenumbers phonenumbers.Service
	tools        tools.Tools
	uuid         uuid.Service
	phones       []phones.PhoneData
}

func NewDetailsPage(tools tools.Tools,
	html html.Writer,
	contacts contacts.Service,
	phonenumbers phonenumbers.Service,
) *DetailsPage {
	return &DetailsPage{
		html:         html,
		tools:        tools,
		uuid:         tools.UUID(),
		contacts:     contacts,
		phonenumbers: phonenumbers,
		phones:       phones.GetData(),
	}
}

func (h *DetailsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/contacts/{contactID}", h.getDetails)
	r.Delete("/web/contacts/{contactID}", h.deleteContact)
	r.Get("/web/contacts/{contactID}/name", h.getEditName)
	r.Post("/web/contacts/{contactID}/name", h.editName)
	r.Get("/web/contacts/{contactID}/phones", h.getRegisterPhoneModal)
	r.Post("/web/contacts/{contactID}/phones", h.registerNewPhone)
	r.Delete("/web/contacts/{contactID}/phones/{phoneID}", h.deletePhone)
}

func (h *DetailsPage) getDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	phones, err := h.phonenumbers.GetAllForContact(ctx, contact, nil)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to fetch the phonenumbers for contact %q: %w", contact.ID(), err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.DetailsPageTmpl{
		Contact: contact,
		Phones:  phones,
	})
}

func (h *DetailsPage) deleteContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	err := h.contacts.Delete(ctx, contact)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to delete contact %q: %w", contact.ID(), err))
		return
	}

	w.Header().Add("HX-Redirect", "/web/contacts")
	w.WriteHeader(http.StatusNoContent)
}

func (h *DetailsPage) getEditName(w http.ResponseWriter, r *http.Request) {
	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ModalEditNameTmpl{
		Contact: contact,
	})
}

func (h *DetailsPage) editName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	contact, err := h.contacts.EditName(ctx, &contacts.EditNameCmd{
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

	h.renderDetailsPage(w, r, contact)
}

func (h *DetailsPage) getRegisterPhoneModal(w http.ResponseWriter, r *http.Request) {
	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.ModalRegisterPhoneNumberTmpl{
		Error:    nil,
		Input:    "",
		Selected: h.phones[0],
		Phones:   h.phones,
		Contact:  contact,
	})
}

func (h *DetailsPage) registerNewPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	region := r.FormValue("region")
	ptype := r.FormValue("type")
	input := r.FormValue("input")

	_, err := h.phonenumbers.Create(ctx, &phonenumbers.CreateCmd{
		Contact: contact,
		Type:    ptype,
		Region:  region,
		Input:   input,
	})
	if errors.Is(err, errs.ErrValidation) {
		selected := h.phones[0]
		for _, elem := range h.phones {
			if elem.Iso2Code == region {
				selected = elem
			}
		}

		h.html.WriteHTMLTemplate(w, r, http.StatusUnprocessableEntity, &contactstmpl.ModalRegisterPhoneNumberTmpl{
			Error:    err,
			Input:    input,
			Selected: selected,
			Phones:   h.phones,
			Contact:  contact,
		})
		return
	}

	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to register the phone: %w", err))
		return
	}

	h.renderDetailsPage(w, r, contact)
}

func (h *DetailsPage) deletePhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := h.getContactFromURL(w, r)
	if contact == nil {
		return
	}

	phoneID, err := h.uuid.Parse(chi.URLParam(r, "phoneID"))
	if err != nil {
		logger.LogEntrySetError(ctx, fmt.Errorf("failed to parse the phoneID from the url: %w", err))
		w.WriteHeader(http.StatusNotFound)
	}

	h.phonenumbers.DeleteContactPhone(ctx, contact, phoneID)
}

func (h *DetailsPage) getContactFromURL(w http.ResponseWriter, r *http.Request) *contacts.Contact {
	contactID, err := h.uuid.Parse(chi.URLParam(r, "contactID"))
	if err != nil {
		logger.LogEntrySetError(r.Context(), fmt.Errorf("failed to parse the contactID from the url: %w", err))
		http.Redirect(w, r, "/web/contacts", http.StatusTemporaryRedirect)
		return nil
	}

	contact, err := h.contacts.GetByID(r.Context(), contactID)
	if errors.Is(err, errs.ErrNotFound) {
		http.Redirect(w, r, "/web/contacts", http.StatusTemporaryRedirect)
		return nil
	}

	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("contact not found: %w", err))
		return nil
	}

	return contact
}

func (h *DetailsPage) renderDetailsPage(w http.ResponseWriter, r *http.Request, contact *contacts.Contact) {
	ctx := r.Context()

	phones, err := h.phonenumbers.GetAllForContact(ctx, contact, nil)
	if err != nil {
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to get the phones numbers for contact %q: %w", contact.ID(), err))
		return
	}

	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &contactstmpl.DetailsPageTmpl{
		Contact: contact,
		Phones:  phones,
	})
}
