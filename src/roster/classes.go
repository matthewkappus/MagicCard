package roster

import (
	"fmt"
	"net/http"
)

func (sv *StaffView) Profile(w http.ResponseWriter, r *http.Request) {

	teacher := sv.GetTeacher(r)
	nav, err := sv.MakeNav(teacher, "teacher", "Student Search", Teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c, err := sv.MakeSchoolClassroom(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sv.tmpls.Lookup("profile").Execute(w, TD{N: nav, C: c})
}

func (sv *StaffView) AddComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	// teacher must match session
	t := sv.GetTeacher(r)
	if t != r.PostFormValue("teacher") {
		fmt.Printf("form.teacher: '%s' doesn't match session.teacher '%s'", r.PostFormValue("teacher"), t)
	}

	// perm_id, teacher, comment, title, cat, isActive
	err := sv.store.AddStarStrike(r.PostFormValue("permid"), r.PostFormValue("teacher"), r.PostFormValue("comment"), r.PostFormValue("title"), r.PostFormValue("icon"), r.PostFormValue("cat"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ClassEdit by section
func (sv *StaffView) ClassEdit(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("section")) != 4 {
		http.NotFound(w, r)
		return
	}
	teacher := sv.GetTeacher(r)

	class, err := sv.MakeClassroom(teacher, r.FormValue("section"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if class.Teacher != teacher {
		http.Error(w, "Not your class", http.StatusForbidden)
		return
	}

	nav, err := sv.MakeNav(teacher, "classroom", class.ClassName, Teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sv.tmpls.Lookup("classedit").Execute(w, TD{N: nav, C: class})
}
