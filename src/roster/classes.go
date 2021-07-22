package roster

import (
	"net/http"
)

func (sv *StaffView) Profile(w http.ResponseWriter, r *http.Request) {

	teacher := sv.GetTeacher(r)
	nav, err := sv.MakeNav(teacher, "teacher", "Student Search")
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

// show class by section
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

	nav, err := sv.MakeNav(teacher, "classroom", class.ClassName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("classedit").Execute(w, TD{N: nav, C: class})
}
