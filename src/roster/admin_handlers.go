package roster

import "net/http"

func (v *View) StaffEdit(w http.ResponseWriter, r *http.Request) {
	staff, err := v.store.GetStaff()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// make a map

	v.tmpls.Lookup("staff_list").Execute(w, StaffData{N: v.N, Teachers: staff})
}

func (v *View) AdminHome(w http.ResponseWriter, r *http.Request) {
	v.tmpls.Lookup("admin_home").Execute(w, TD{N: v.N})
}