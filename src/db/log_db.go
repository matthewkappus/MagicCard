package db

const (
	createLog = ` CREATE TABLE IF NOT EXISTS ? (logType TEXT, when DATETIME, message TEXT)`

	insertLog = `INSERT INTO ? (logType, type, when, message) VALUES (?, ?, ?, ?)`
)

func (s *Store) NewLogTable(name string) error {
	_, err := s.db.Exec(createLog, name)
	return err
}

func (s *Store) InsertLog(table string, logType string, when string, message string) error {
	_, err := s.db.Exec(insertLog, table, logType, when, message)
	return err
}
