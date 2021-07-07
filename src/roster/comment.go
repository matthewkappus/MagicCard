package roster

import "time"

type Comment struct {
	ID       int
	PermID   string
	Email    string
	IsMerrit bool
	Created  time.Time
	Comment  string
	IsActive bool
}
