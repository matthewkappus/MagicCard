package roster

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (sv *StaffView) StarBarEdit(w http.ResponseWriter, r *http.Request) {

	teacher := sv.GetTeacher(r)
	list, _ := sv.store.ListClasses(teacher)
	ci := &ClassInfo{

		Teacher:   teacher,
		ClassList: list,
		Title:     "StarBar Edit",
		Path:      "profile",
	}

	data := struct {
		Info    *ClassInfo
		StarBar *comment.StarBar
	}{
		Info: ci,
	}

	if r.FormValue("id") != "" {
		data.StarBar, _ = sv.store.GetStarBarByID(r.FormValue("id"))
	}
	if r.Method == http.MethodPost {

		idInt := 0
		if r.PostFormValue("id") != "" {
			idInt, _ = strconv.Atoi(r.PostFormValue("id"))

		}

		sb := comment.StarBar{
			ID:      idInt,
			Teacher: teacher,
			Title:   r.PostFormValue("title"),
			Comment: r.PostFormValue("comment"),
			IsStar:  r.PostFormValue("isStar") == "true",
		}
		if !sb.IsValid() {
			fmt.Fprint(w, "Not a valid starbar")
		}

		// update starbar
		if idInt > 0 {
			err := sv.store.UpdateStarBar(sb.ID, sb.Teacher, sb.Teacher, sb.Comment, sb.IsStar)
			ci.StatusMsg = err.Error()
		}
		// create starbar
		_, err := sv.store.AddStarBar(sb.Teacher, sb.Teacher, sb.Comment, sb.IsStar)
		if err != nil {
			ci.StatusMsg = err.Error()
		}

		http.Redirect(w, r, "/profile", http.StatusFound)
	} else {
		// get starbar

		sv.tmpls.Lookup("starbaredit").Execute(w, data)

	}

}

//     type StarBar struct {
// 	ID      int
// 	Teacher string
// 	Title   string
// 	Comment string
// 	// if not star, then bar
// 	IsStar bool
// }
