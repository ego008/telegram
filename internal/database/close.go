package database

func (db *DataBase) Close() error {
	return db.DB.Close()
}
