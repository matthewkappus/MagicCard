package roster

import (
	"fmt"
	"net/http"
	"strconv"
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

// AddMyStarStrikeAll creates a starstrike for all teachers to assign
func (sv *StaffView) AddMyStarStrikeAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// todo: check if Session.Type == Admin

	// (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)
	teacher := "all"
	comment := r.PostFormValue("comment")
	title := r.PostFormValue("title")
	icon := r.PostFormValue("icon")

	// todo: verify cat is 0-1
	// switch r.PostFormValue("cat") {
	// case "0":
	// 	c = comment.Star
	// case "1":
	// 	c = comment.MinorStrike
	// case "2":
	// 	c = comment.MajorStrike
	// case "3":
	// 	c = comment.MajorStrike
	// }

	c, err := strconv.Atoi(r.PostFormValue("cat"))
	if c < 0 || c > 3 || err != nil {

		http.Error(w, "Invalid cat: "+r.PostFormValue("cat"), http.StatusNotAcceptable)
		return
	}

	fmt.Println("inserting new mystarstrike", teacher, comment, title, icon)
	// todo: make title unique
	err = sv.store.InsertMyStarStrike(teacher, comment, title, icon, c)
	if err != nil {
		fmt.Println("error insrting mystarstrike", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (sv *StaffView) MyStarStrikeForm(w http.ResponseWriter, r *http.Request) {

	sv.tmpls.Lookup("mystarstrikeform").Execute(w, nil)

}
