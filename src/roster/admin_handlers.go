package roster

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/matthewkappus/MagicCard/src/comment"
)

func (v *View) StaffEdit(w http.ResponseWriter, r *http.Request) {
	staff, err := v.store.GetStaff()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.tmpls.Lookup("staff_list").Execute(w, StaffData{N: v.N, Teachers: staff})
}

func (v *View) AdminHome(w http.ResponseWriter, r *http.Request) {
	v.tmpls.Lookup("admin_home").Execute(w, TD{N: v.N})
}

func (v *View) AddPerfectAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("stu401")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// todo: check if file is a csv
	// https://freshman.tech/file-upload-golang/

	perms, err := Get401Perms(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := &comment.StarStrike{
		Teacher:  "Madison Attendance",
		Title:    "Perfect Attendance",
		Created:  time.Now(),
		Cat:      comment.Star,
		Icon:     "award",
		IsActive: true,
	}

	ss := comment.BatchCreateStarstrike(s, perms)

	if err = v.store.BatchAddStarStrikes(ss); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, " Award %d Badges", len(perms))
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)

}

func Get401Perms(r io.Reader) ([]string, error) {

	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}

	var perms []string

	for _, row := range rows {
		// len()
		if len(row[9]) == 9 {
			perms = append(perms, row[9])
		}
	}
	return perms, nil
}

func ParseAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("stu401")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// todo: check if file is a csv
	// https://freshman.tech/file-upload-golang/

	perms, err := Get401Perms(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Perms: %v", perms)
}
