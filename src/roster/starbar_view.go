package roster

import (
	"fmt"
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (sv *StaffView) StarBarEdit(w http.ResponseWriter, r *http.Request) {

	// get starbar shows edit form based on form data id
	teacher := sv.GetTeacher(r)
	list, _ := sv.store.ListClasses(teacher)
	ci := &ClassInfo{
		Teacher:   teacher,
		ClassList: list,
		Title:     "StarBar Edit",
		Path:      "profile",
	}

	sb := new(comment.StarBar)
	if r.FormValue("id") != "" {
		sb, _ = sv.store.GetStarBarByID(r.FormValue("id"))
		if sb.Teacher != teacher {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
	}

	data := struct {
		Info    *ClassInfo
		StarBar *comment.StarBar
	}{
		Info:    ci,
		StarBar: sb,
	}

	fmt.Println("getting editform")
	sv.tmpls.Lookup("starbaredit").Execute(w, data)

}

func (sv *StaffView) StarBarCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("parsed form:\n %+v\n", r.Form)

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

func (sv *StaffView) updateStarBar(w http.ResponseWriter, r *http.Request) {

	// sb := new(comment.StarBar)
	// // parse posted data
	// idInt := 0
	// if r.PostFormValue("id") != "" {
	// 	idInt, _ = strconv.Atoi(r.PostFormValue("id"))

	// }

	// sb.ID = idInt
	// sb.Title = r.PostFormValue("title")
	// sb.Comment = r.PostFormValue("comment")
	// sb.IsStar = r.PostFormValue("isStar") == "true"

	// // if !sb.IsValid() {
	// // 	fmt.Fprint(w, "Not a valid starbar")
	// // }

	// // update starbar
	// if idInt > 0 {
	// 	fmt.Printf("updating sb %d", idInt)
	// 	sv.store.UpdateStarBar(sb.ID, sb.Teacher, sb.Teacher, sb.Comment, sb.IsStar)
	// } else {
	// 	fmt.Printf("creating sb %s", sb.Title)

	// 	sv.store.AddStarBar(sb.Teacher, sb.Teacher, sb.Comment, sb.IsStar)

	// }

	// http.Redirect(w, r, "/profile", http.StatusFound)
}
