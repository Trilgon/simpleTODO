package postgresql

const (
	newUser       = `INSERT INTO users (email, enc_password) VALUES ($1, $2) RETURNING id;`
	authorizeUser = `SELECT id FROM users WHERE email = $1 AND enc_password = $2 LIMIT 1;`
	checkUser     = `SELECT exists(SELECT 1 FROM users WHERE email = $1);;`
	signIn        = `UPDATE users SET is_authorized = true WHERE email = $1;`
	signOut       = `UPDATE users SET is_authorized = false WHERE email = $1;`
	addNote       = `INSERT INTO notes (email, title, text, start_date) VALUES ($1, $2, $3, $4) RETURNING id;`
	deleteNote    = `DELETE FROM notes WHERE id = $1`
	updateNote    = `UPDATE notes SET title = $1, text = $2 WHERE id = $3;`
	markNote      = `UPDATE notes SET is_done = $1 WHERE id = $2;`
	getById       = `SELECT * FROM notes WHERE id = $1;`
	getByEmail    = `SELECT * FROM notes WHERE email = $1;`
)