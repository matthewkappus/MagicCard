package db

import "github.com/matthewkappus/MagicCard/src/comment"

const (
	createStarBar          = `CREATE TABLE IF NOT EXISTS starbar(id INTEGER PRIMARY KEY AUTOINCREMENT, teacher, title, comment TEXT, isStar BOOLEAN, FOREIGN KEY(teacher) REFERENCES staff(name))`
	createStarBarSave      = `INSERT OR REPLACE INTO starbar(id, teacher, title, comment, isStar) VALUES(?, ?, ?, ?, ?)`
	selectStarBarByTeacher = `SELECT * FROM starbar WHERE teacher = ?`
	insertStarBar          = `INSERT INTO starbar(teacher, title, comment, isStar) VALUES(?, ?, ?, ?)`
)

func (s *Store) AddStarBar(teacher, title, comments string, isStar bool) (*comment.StarBar, error) {
	sb := &comment.StarBar{
		Teacher: teacher,
		Title:   title,
		Comment: comments,
		IsStar:  isStar,
	}

	// sb.IsValid()
	res, err := s.db.Exec(insertStarBar, teacher, title, comments, isStar)
	if err != nil {
		return nil, err
	}
	i64, err := res.LastInsertId()
	sb.ID = int(i64)
	if err != nil {
		return nil, err
	}
	return sb, nil

}

// GetStarBars takes a staff.name and returns their StarBars or an error
func (s *Store) GetStarBars(teacher string) ([]*comment.StarBar, error) {
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
