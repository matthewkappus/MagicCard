package roster

import (
	"fmt"
	"net/http"
	"strconv"
)

func (v *View) Profile(w http.ResponseWriter, r *http.Request) {

	c, err := v.MakeSchoolClassroom(v.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.tmpls.Lookup("profile").Execute(w, TD{N: v.N, C: c})
}

func (v *View) AddComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	if v.User != r.PostFormValue("teacher") {
		fmt.Printf("form.teacher: '%s' doesn't match session.teacher '%s'", r.PostFormValue("teacher"), v.User)
		http.Error(w, "Not your class", http.StatusForbidden)
		return
	}

	// perm_id, teacher, comment, title, cat, isActive
	err := v.store.AddStarStrike(r.PostFormValue("permid"), r.PostFormValue("teacher"), r.PostFormValue("comment"), r.PostFormValue("title"), r.PostFormValue("icon"), r.PostFormValue("cat"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// SendAlert corresponds to ReadAlert in the referrer handler
	v.SendAlert(w, &Alert{Message: "StarStrike added", Type: "success"})
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// ClassEdit by section
func (v *View) ClassEdit(w http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("section")) != 4 {
		http.NotFound(w, r)
		return
	}

	class, err := v.MakeClassroom(v.User, r.FormValue("section"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if class.Teacher != v.User {
		http.Error(w, "Not your class", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.tmpls.Lookup("classedit").Execute(w, TD{N: v.N, C: class})
}

// AddMyStarStrikeAll creates a starstrike for all teachers to assign
func (v *View) AddMyStarStrikeAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// (id INTEGER PRIMARY KEY, teacher TEXT, comment TEXT, title TEXT, icon TEXT, created DATETIME DEFAULT CURRENT_TIMESTAMP, cat INTEGER, isActive BOOLEAN DEFAULT true)
	teacher := "all"
	comment := r.PostFormValue("comment")
	title := r.PostFormValue("title")
	icon := r.PostFormValue("icon")

	c, err := strconv.Atoi(r.PostFormValue("cat"))
	if c < 0 || c > 3 || err != nil {

		http.Error(w, "Invalid cat: "+r.PostFormValue("cat"), http.StatusNotAcceptable)
		return
	}

	// todo: make title unique
	err = v.store.InsertMyStarStrike(teacher, comment, title, icon, c)
	if err != nil {
		fmt.Println("error insrting mystarstrike", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

// SendAlert corresponds to ReadAlert in the referrer handler
v.SendAlert(w, &Alert{Message: title +" added", Type: "success"})
http.Redirect(w, r, r.Referer(), http.StatusSeeOther)

}

func (v *View) MyStarStrikeForm(w http.ResponseWriter, r *http.Request) {

	v.tmpls.Lookup("mystarstrikeform").Execute(w, nil)

}
