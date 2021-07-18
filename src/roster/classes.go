package roster

import (
	"fmt"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

// ClassInfo is teachers class info
type ClassInfo struct {
	Stu415s   []*synergy.Stu415
	Stars     []*comment.StarStrike
	Strikes   []*comment.StarStrike
	ClassList []*synergy.Stu415
	Teacher   string
	ClassName string
	Title     string
	StatusMsg string
	// shoud be a css class
	Path string
}

func (sv *StaffView) Profile(w http.ResponseWriter, r *http.Request) {
	// todo: render template, add class info

	teacher := sv.GetTeacher(r)
	stars, strikes, _ := sv.store.GetTeacherStarStrikes(teacher)
	list, _ := sv.store.ListClasses(teacher)

	ci := &ClassInfo{
		Path:      "profile",
		ClassList: list,
		Stars:     stars,
		Strikes:   strikes,
		Title:     teacher + " Profile",
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
	classinfo := &ClassInfo{
		Stu415s:   s415s,
		Stars:     stars,
		Strikes:   stikes,
		Teacher:   teacher,
		ClassList: list,
		ClassName: s415s[0].CourseIDAndTitle,
		Title:     s415s[0].CourseIDAndTitle + " Class List",
	}

	// todo: wrap s415s in struct with class info and tags
	sv.tmpls.Lookup("classedit").Execute(w, classinfo)
}
