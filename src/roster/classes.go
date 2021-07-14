package roster

import (
	"fmt"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

func (sv *StudentView) ListClasses(w http.ResponseWriter, r *http.Request) {
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

	teacher := sv.GetTeacher(r)
	sbs, err := sv.store.GetStarBars(teacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	classinfo := struct {
		Stu415s  []*synergy.Stu415
		StarBars []*comment.StarBar
		Teacher  string
		Title    string
	}{
		Stu415s:  s415s,
		StarBars: sbs,
		Teacher:  teacher,
		Title:    "Class List",
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("studentlist").Execute(w, classinfo)
}
