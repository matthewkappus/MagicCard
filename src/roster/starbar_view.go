package roster

import (
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (sv *StaffView) StarBarCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		teacher := sv.GetTeacher(r)
		list, _ := sv.store.ListClasses(teacher)
		ci := &ClassInfo{
			Teacher:   teacher,
			ClassList: list,
			Title:     "StarBar Edit",
			Path:      "profile",
		}
		sv.tmpls.Lookup("starbarcreate").Execute(w, ci)
	}

	if r.Method == http.MethodPost {
		sb := new(comment.StarBar)
		sb.Teacher = sv.GetTeacher(r)
		sb.Title = r.PostFormValue("title")
		sb.Comment = r.PostFormValue("comment")
		sb.IsStar = r.PostFormValue("isStar") == "true"

		if !sb.IsValid() {

			http.Error(w, "Invalid Form", http.StatusBadRequest)
		} else {
			if _, err := sv.store.AddStarBar(sb.Teacher, sb.Title, sb.Comment, sb.IsStar); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/profile", http.StatusFound)

}

func (sv *StaffView) StarBarDelete(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("id") == "" {
		http.Error(w, "No ID", http.StatusBadRequest)
		return
	}
	// check that teacher is owner before deleting
	sb, err := sv.store.GetStarBarByID(r.FormValue("id"))
	if err != nil || sb.Teacher != sv.GetTeacher(r) {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	if err := sv.store.DeleteStarBarByID(r.FormValue("id")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusFound)
}

func (sv *StaffView) StarBarEdit(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		if r.FormValue("id") == "" {
			http.Error(w, "No ID", http.StatusBadRequest)
			return
		}
		teacher := sv.GetTeacher(r)
		list, _ := sv.store.ListClasses(teacher)
		ci := &ClassInfo{
			Teacher:   teacher,
			ClassList: list,
			Title:     "StarBar Edit",
			Path:      "profile",
		}
		sb, _ := sv.store.GetStarBarByID(r.FormValue("id"))

		data := struct {
			Info    *ClassInfo
			StarBar *comment.StarBar
		}{
			StarBar: sb,
			Info:    ci,
		}
		sv.tmpls.Lookup("starbaredit").Execute(w, data)
	}

	if r.Method == http.MethodPost {
		sb := new(comment.StarBar)
		sb.Teacher = sv.GetTeacher(r)
		sb.Title = r.PostFormValue("title")
		sb.Comment = r.PostFormValue("comment")
		sb.IsStar = r.PostFormValue("isStar") == "true"

		if !sb.IsValid() {

			http.Error(w, "Invalid Form", http.StatusBadRequest)
		} else {
			if _, err := sv.store.AddStarBar(sb.Teacher, sb.Title, sb.Comment, sb.IsStar); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}
