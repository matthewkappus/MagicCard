package roster

import (
	"net/http"
	"strconv"
	"time"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

// AddContact parses Post data and adds a contact to the database
func (v *View) AddContact(w http.ResponseWriter, r *http.Request) {
	// sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message

	s := &synergy.Staff{
		Name: r.PostFormValue("teacher"),
	}
	s.Email, _ = v.store.GetTeacherEmail(s.Name)

	// todo add validation
	// sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message
	ssid := r.PostFormValue("ss_id")

	// todo add validation

	if ssid == "" {
		http.Error(w, "No starstrike id provided", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(r.PostFormValue("ss_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ss, err := v.store.GetStarStrike(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	c := &comment.Contact{
		Sender:     s,
		StarStrike: ss,
		Sent:       time.Now(),
		Respondent: r.PostFormValue("respondent"),
		Message:    r.PostFormValue("message"),
		IsClosed:   r.PostFormValue("is_closed") == "true",
	}
	// todo add validation
	if err := v.store.InsertContact(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	v.SendAlert(w, &Alert{Message: "Contact message saved", Type: "success"})

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

// ContactLog shows the contact log for a student
func (v *View) ContactLog(w http.ResponseWriter, r *http.Request) {
	perm := r.URL.Query().Get("id")
	name, err := v.store.GetStudentName(perm)
	if err != nil || name == "" {
		http.Error(w, "No student found for "+perm, http.StatusBadRequest)
		return
	}
	c, err := v.makeContactMap(perm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the student's name

	v.tmpls.ExecuteTemplate(w, "contact_log", ContactData{N: v.N, C: c, StudentName: formatName(name)})
}

// makeContactMap
func (v *View) makeContactMap(perm string) (map[*comment.StarStrike]*comment.Contact, error) {
	cm := make(map[*comment.StarStrike]*comment.Contact)

	// get student strikes
	ss, err := v.store.GetStudentStrikes(perm)
	if err != nil {
		return nil, err
	}

	// for each strike, get the contact
	for _, s := range ss {
		c, err := v.store.GetContact(s.ID)
		if err != nil {
			cm[s] = nil
			continue
		}
		cm[s] = c
	}

	return cm, nil
}



