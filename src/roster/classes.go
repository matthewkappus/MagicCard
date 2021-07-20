package roster

import (
	"fmt"
	"net/http"
)

func (sv *StaffView) Profile(w http.ResponseWriter, r *http.Request) {
	// todo: render template, add class info

	teacher := sv.GetTeacher(r)
	stars, strikes, _ := sv.store.GetTeacherStarStrikes(teacher)
	list, _ := sv.store.ListClasses(teacher)

	ci := &Classroom{
		ClassList: list,
		MyStars:   stars,
		MyStrikes: strikes,
	}

	sv.tmpls.Lookup("profile").Execute(w, ci)
}

// show class by section
func (sv *StaffView) ClassEdit(w http.ResponseWriter, r *http.Request) {
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

	teacher := sv.GetTeacher(r)
	stars, stikes, err := sv.store.GetTeacherStarStrikes(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// todo: put in helper function
	list, _ := sv.store.ListClasses(teacher)
	classinfo := &Classroom{
		Stu415s:   s415s,
		MyStars:   stars,
		MyStrikes: stikes,
		Teacher:   teacher,
		ClassList: list,
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("classedit").Execute(w, classinfo)
}
