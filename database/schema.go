package database

import "database/sql"

func RunSchema(db *sql.DB) {

	const create string = `
	CREATE TABLE IF NOT EXISTS user (
		id TEXT NOT NULL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS user_session (
		id TEXT NOT NULL PRIMARY KEY,
		user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
		active_expires INTEGER NOT NULL,
		idle_expires INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS projects(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		creator_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
		name TEXT NOT NULL 
	)
	`

	db.Exec(create)

}
