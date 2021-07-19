package comment

import "time"

// Todo: Remove card
type Card struct {
	ID int `json:"id,omitempty"`
	//
	PermID string `json:"perm_id,omitempty"`
	// staff(name)
	Teacher string `json:"teacher,omitempty"`
	Comment string `json:"comment,omitempty"`
	// Title is a catagory of the comment
	Title string `json:"title,omitempty"`

	Created time.Time `json:"created,omitempty"`
	// max comment: 280
	IsStar   bool `json:"is_star,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
}
