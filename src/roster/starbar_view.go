package roster

import (
	"net/http"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (sv *StaffView) StarBarEdit(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		teacher := sv.GetTeacher(r)

		// todo: move staff templates to a separate folder
		// sv.RenderTemplate(w, "staff/starbar_edit.html", map[string]interface{}

		ci := &ClassInfo{

			Teacher: teacher,
			Title:   "StarBar Edit",
			Path:    "profile",
		}

		sb, err := sv.store.GetStarBarByID(r.FormValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Info    *ClassInfo
			StarBar *comment.StarBar
		}{
			Info:    ci,
			StarBar: sb,
		}

		sv.tmpls.Lookup("starbaredit").Execute(w, data)
	}

}
