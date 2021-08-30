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
