package roster

import (
	"fmt"
	"net/http"
)

func (sv *StudentView) ListClasses(w http.ResponseWriter, r *http.Request) {
	teacherCookie, err := r.Cookie("teacher")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("getting classes for", teacherCookie.Value)
	classes, err := sv.store.ListClasses(teacherCookie.Value)
	if err != nil {
		http.NotFound(w, r)
	}
	sv.tmpls.Lookup("classlist.tmpl.html").Execute(w, classes)
}

// show class by section
func (sv *StudentView) Class(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("section")) != 4 {
		http.NotFound(w, r)
		return
	}

	s415s, err := sv.store.ListStudents(r.FormValue("section"))
	if err != nil {
		http.NotFound(w, r)
		fmt.Printf("error looking for %s:\n%v ", r.FormValue("section"), err)
		return
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("studentlist.tmpl.html").Execute(w, s415s)
}
