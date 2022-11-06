package postgresql

const (
	isUserExists = `SELECT exists(SELECT 1 FROM users WHERE email = $1)`
	signUp       = `INSERT INTO users (email, enc_password) VALUES ($1, $2)`
	signIn       = `SELECT exists(SELECT 1 FROM users WHERE email = $1 AND enc_password = $2)`
)
