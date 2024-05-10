package queries

const (
	CreateUserQuery = `
        INSERT INTO users (firstname,lastname,email,age)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	UpdateUserQuery = `
        UPDATE users
        SET firstname = $1, lastname = $2, email = $3, age = $4
        WHERE id = $5`

	SelectUserByIDQuery = `
        SELECT firstname, lastname, email, age
        FROM users
        WHERE id = $1`

	CreateUserTableQuery = `
        CREATE TABLE IF NOT EXISTS users (
			id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
)
