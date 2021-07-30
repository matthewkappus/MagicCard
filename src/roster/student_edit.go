package roster

import (
	"fmt"
	"net/http"
)

// MagicCard shows students StarStrikes and allows you to edit
func (sv *StaffView) MagicCard(w http.ResponseWriter, r *http.Request) {
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	mc, err := sv.MakeStudent(permid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nav, err := sv.MakeNav(sv.GetTeacher(r), "/student", fmt.Sprintf("%s Magic Card", mc.S415.StudentName), Teacher)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, TD{M: mc, N: nav})
}
