package db

import "github.com/matthewkappus/MagicCard/src/comment"

const (
	createStarBar = `CREATE TABLE IF NOT EXISTS starbar(id PRIMARY KEY, teacher, title, comment TEXT, isStar BOOLEAN, FOREIGN KEY(teacher) REFERENCES staff(name))`
	createStarBarSave = `INSERT OR REPLACE INTO starbar(id, teacher, title, comment, isStar) VALUES(?, ?, ?, ?, ?)`
	selectStarBarByTeacher = `SELECT * FROM starbar WHERE teacher = ?`
)

// func (s *Store) AddStarBar(teacher, title, comment string, isStar bool) (*comment.StarBar, error) {

// }

// GetStarBars takes a staff.name and returns their StarBars or an error
func (s  *Store) GetStarBars(teacher string)([]*comment.StarBar, error)	  {
	rows, err := s.db.Query(selectStarBarByTeacher, teacher)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	starBars := make([]*comment.StarBar, 0)
	for rows.Next() {
		sb := new(comment.StarBar)
		if err := rows.Scan(&sb.ID, &sb.Teacher, &sb.Title, &sb.Comment, &sb.IsStar); err != nil {
			continue
		}
		starBars = append(starBars, sb)
	}
	return starBars, nil
}