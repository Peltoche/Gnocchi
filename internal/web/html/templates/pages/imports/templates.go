package imports

type ImportsPageTmpl struct {
	ErrorMsg string
}

func (t *ImportsPageTmpl) Template() string {
	return "pages/imports/page_imports"
}
