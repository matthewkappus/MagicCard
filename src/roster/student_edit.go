package roster

import (
	"net/http"

	"github.com/matthewkappus/Roster/src/synergy"
)

// MagicCard shows students StarStrikes and allows you to edit
func (sv *StaffView) MagicCard(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	stu415, err := sv.store.SelectStu415(permid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, _ := sv.ClassList(r)
	teacher := sv.GetTeacher(r)
	stars, strikes, _ := sv.store.GetTeacherStarStrikes(teacher)

	ci := &ClassInfo{
		ClassList: classes,
		Stu415s:   []*synergy.Stu415{stu415},
		Stars:     stars,
		Strikes:   strikes,
		Teacher:   teacher,
		Title:     stu415.StudentName + " Magic Card",
		Path:      "classes",
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, ci)
}

// depricated
// todo: remove
func (sv *StaffView) Card(w http.ResponseWriter, r *http.Request) {
	// todo: look up student in db: join with comments
	permid := r.FormValue("id")
	if len(permid) != 9 {
		http.NotFound(w, r)
		return
	}

	stu415, err := sv.store.SelectStu415(permid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classes, _ := sv.ClassList(r)
	teacher := sv.GetTeacher(r)
	stars, strikes, _ := sv.store.GetTeacherStarStrikes(teacher)

	ci := &ClassInfo{
		ClassList: classes,
		Stu415s:   []*synergy.Stu415{stu415},
		Stars:     stars,
		Strikes:   strikes,
		Teacher:   teacher,
		Title:     stu415.StudentName + " Magic Card",
		Path:      "classes",
	}
	// todo: add email and session info
	sv.tmpls.Lookup("card").Execute(w, ci)
}
