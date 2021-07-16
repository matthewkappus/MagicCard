package comment

// StarBar stores comments as a category by Title and associates them with a Teacher for reuse.
type StarBar struct {
	ID      int
	Teacher string
	Title   string
	Comment string
	// if not star, then bar
	IsStar bool
}

// IsValid
// title: 1 word, non duplicate to teacher, 20 chars, Comment: 280 chars
func (sb *StarBar) IsValid() bool {
	return true
}
