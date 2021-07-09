package roster

import (
	"fmt"
	"net/http"
)

func (sv *StudentView) ListClasses(w http.ResponseWriter, r *http.Request) {
	// todo: get from session email
	classes, err := sv.store.ListClasses("Susco Taylor, Kevin R.")
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

	fmt.Println("class listing", s415s)

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("studentlist.tmpl.html").Execute(w, s415s)
}
