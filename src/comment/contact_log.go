package comment

import "time"

// todo: move from db
type Teacher struct {
	Teacher    string
	FullName   string
	StaffEmail string
}

// Contact logs teacher-parent communications regarding starstrikes
type Contact struct {
	ID int

	Sender Teacher

	StudentName string

	// when contacted
	Sent time.Time

	// Respondent usually is parent of student
	Respondent string

	StarStrike *StarStrike

	Message string
}

func NewContact(ss *StarStrike, sender Teacher, resp, msg string) *Contact {
	return &Contact{
		StarStrike: ss,
		Sender:     sender,
		Respondent: resp,
		Message:    msg,
	}
}
