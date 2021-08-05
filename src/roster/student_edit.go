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

// StudentCard shows magic card belonging to student user by their cookie
func (v *View) StudentCard(w http.ResponseWriter, r *http.Request) {

	_, user, _, scope, err := sessionCookies(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if scope != Student {
		http.Error(w, "User Not in Student Scope", http.StatusUnauthorized)
		return
	}

	mc, err := v.MakeStudentMagicCard(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: add email and session info
	v.tmpls.Lookup("studentmagiccard").Execute(w, TD{M: mc, N: v.N})
}
