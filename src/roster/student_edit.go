package roster

import (
	"net/http"
)

// MagicCard shows students StarStrikes and allows you to edit
func (v *View) MagicCard(w http.ResponseWriter, r *http.Request) {
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	mc, err := v.MakeStudentMagicCard(permid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: add email and session info
	v.tmpls.Lookup("studentmagiccard").Execute(w, TD{M: mc, N: v.N})
}
