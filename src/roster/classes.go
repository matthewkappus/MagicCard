package roster

import (
	"fmt"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

type ClassInfo struct {
	Stu415s   []*synergy.Stu415
	Stars     []*comment.StarBar
	Bars      []*comment.StarBar
	ClassList []*synergy.Stu415
	Teacher   string
	ClassName string
	Title     string
}

func (sv *StaffView) ListClasses(w http.ResponseWriter, r *http.Request) {
	// todo: allow all staff to magiccard
	teacherCookie, err := r.Cookie("teacher")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, err := sv.store.ListClasses(teacherCookie.Value)
	if err != nil {
		http.NotFound(w, r)
	}

	sv.tmpls.Lookup("classlist").Execute(w, classes)
}

func (sv *StaffView) Profile(w http.ResponseWriter, r *http.Request) {
	// todo: render template, add class info

	fmt.Fprintf(w, "welcome %s", sv.GetTeacher(r))
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
	stars, bars, err := sv.store.GetStarBars(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// todo: put in helper function
	list, _ := sv.store.ListClasses(teacher)
	classinfo := &ClassInfo{
		Stu415s:   s415s,
		Stars:     stars,
		Bars:      bars,
		Teacher:   teacher,
		ClassList: list,
		ClassName: s415s[0].CourseIDAndTitle,
		Title:     s415s[0].CourseIDAndTitle + " Class List",
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("classedit").Execute(w, classinfo)
}
