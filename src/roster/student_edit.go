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

	student, err := sv.MakeStudent(permid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nav, err := sv.MakeNav(sv.GetTeacher(r), "/student", fmt.Sprintf("%s Magic Card", student.S415.StudentName))

	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, TD{S: student, N: nav})
}
