package postgresql

const (
	newUser       = `INSERT INTO users (email, enc_password) VALUES ($1, $2) RETURNING id;`
	authorizeUser = `SELECT id FROM users WHERE email = $1 AND enc_password = $2 LIMIT 1;`
	checkUser     = `SELECT exists(SELECT 1 FROM users WHERE email = $1);`
	signIn        = `UPDATE users SET is_authorized = true WHERE email = $1;`
	signOut       = `UPDATE users SET is_authorized = false WHERE email = $1;`
	addNote       = `INSERT INTO notes (email, title, text, start_date) VALUES ($1, $2, $3, $4) RETURNING id;`
	deleteNotes   = `DELETE FROM notes WHERE email = $1 AND id = ANY ($2);`
	updateNote    = `UPDATE notes SET title = $1, text = $2 WHERE id = $3 AND email = $4;`
	markDone      = `UPDATE notes SET is_done = true, end_date = $1 WHERE id = $2 AND email = $3;`
	markUndone    = `UPDATE notes SET is_done = false, end_date = null WHERE id = $1 AND email = $2;`
	getById       = `SELECT * FROM notes WHERE id = $1 AND email = $2;`
	getByEmail    = `SELECT * FROM notes WHERE email = $1;`
	getByText     = `SELECT * FROM notes WHERE email = $1 AND title LIKE $2 OR text LIKE $2`
)
