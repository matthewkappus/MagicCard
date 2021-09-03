package roster

import (
	"net/http"
	"strconv"
	"time"

	"github.com/matthewkappus/MagicCard/src/comment"
	"github.com/matthewkappus/Roster/src/synergy"
)

func (v *View) AddContact(w http.ResponseWriter, r *http.Request) {
	// sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message

	s := &synergy.Staff{
		Name: r.PostFormValue("teacher"),
	}
	s.Email, _ = v.store.GetTeacherEmail(s.Name)

	// todo add validation
	// sender_name, sender_fullname, sender_email, student_name, sent, respondent, starstrike, message
	ssid := r.PostFormValue("ss_id")
	ss := new(comment.StarStrike)
	if ssid != "" {

		i, err := strconv.Atoi(r.PostFormValue("ss_id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ss, err = v.store.GetStarStrike(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	}

	c := &comment.Contact{
		Sender:     s,
		StarStrike: ss,
		Sent:       time.Now(),
		Respondent: r.PostFormValue("respondent"),
		Message:    r.PostFormValue("message"),
	}
	// todo add validation
	if err := v.store.InsertContact(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v.SendAlert(w, &Alert{Message: "Contact message saved", Type: "success"})

	http.Redirect(w, r, "/", http.StatusFound)
}
func (v *View) ContactLog(w http.ResponseWriter, r *http.Request) {

	c, err := v.store.GetContacts(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	v.tmpls.ExecuteTemplate(w, "contact_log", ContactData{N: v.N, C: c})
}
