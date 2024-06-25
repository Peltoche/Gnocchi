package imports

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peltoche/gnocchi/internal/service/vcard"
	"github.com/Peltoche/gnocchi/internal/tools/router"
	"github.com/Peltoche/gnocchi/internal/web/html"
	importsmpl "github.com/Peltoche/gnocchi/internal/web/html/templates/pages/imports"
	"github.com/go-chi/chi/v5"
)

type ImportsPage struct {
	html  html.Writer
	vcard vcard.Service
}

func NewImportsPage(html html.Writer, vcard vcard.Service) *ImportsPage {
	return &ImportsPage{
		html:  html,
		vcard: vcard,
	}
}

func (h *ImportsPage) Register(r chi.Router, mids *router.Middlewares) {
	if mids != nil {
		r = r.With(mids.Defaults()...)
	}

	r.Get("/web/imports", h.getImportsPage)
	r.Post("/web/imports/vcs", h.importFile)
}

func (h *ImportsPage) getImportsPage(w http.ResponseWriter, r *http.Request) {
	h.html.WriteHTMLTemplate(w, r, http.StatusOK, &importsmpl.ImportsPageTmpl{})
}

func (h *ImportsPage) importFile(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		h.html.WriteHTMLTemplate(w, r, http.StatusUnprocessableEntity, &importsmpl.ImportsPageTmpl{
			ErrorMsg: "missing file",
		})
		return
	}
	defer file.Close()

	err = h.vcard.ImportVCardFile(r.Context(), file)
	switch {
	case err == nil:
		http.Redirect(w, r, "/web/contacts", http.StatusFound)
		return
	case errors.Is(err, vcard.ErrUnsupportedVCardVersion) ||
		errors.Is(err, vcard.ErrInvalidVCard):
		h.html.WriteHTMLTemplate(w, r, http.StatusUnprocessableEntity, &importsmpl.ImportsPageTmpl{
			ErrorMsg: err.Error(),
		})
	default:
		h.html.WriteHTMLErrorPage(w, r, fmt.Errorf("failed to import the vcard file: %w", err))
	}
}
