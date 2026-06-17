package db

const (
	// Table creation schema queries
	CreateBlogsTableQuery = `
		CREATE TABLE IF NOT EXISTS blogs(
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL,
			blogcontent TEXT NOT NULL,
			image VARCHAR(255) NOT NULL,
			category VARCHAR(255) NOT NULL,
			author VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	CreateCommentsTableQuery = `
		CREATE TABLE IF NOT EXISTS comments(
			id SERIAL PRIMARY KEY,
			comment VARCHAR(255) NOT NULL,
			userid VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			blogid VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	CreateSavedBlogsTableQuery = `
		CREATE TABLE IF NOT EXISTS savedblogs(
			id SERIAL PRIMARY KEY,
			userid VARCHAR(255) NOT NULL,
			blogid VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	// CRUD Blog operations
	InsertBlogQuery = `
		INSERT INTO blogs (title, description, image, blogcontent, category, author) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, title, description, image, blogcontent, category, author, created_at;`

	SelectBlogByIDQuery = `
		SELECT id, title, description, image, blogcontent, category, author, created_at 
		FROM blogs 
		WHERE id = $1;`

	UpdateBlogQuery = `
		UPDATE blogs 
		SET title = $1, description = $2, image = $3, blogcontent = $4, category = $5 
		WHERE id = $6 
		RETURNING id, title, description, image, blogcontent, category, author, created_at;`

	DeleteBlogByIDQuery = `
		DELETE FROM blogs 
		WHERE id = $1;`

	// Dependent table cascades
	DeleteCommentsByBlogIDQuery = `
		DELETE FROM comments 
		WHERE blogid = $1;`

	DeleteSavedBlogsByBlogIDQuery = `
		DELETE FROM savedblogs 
		WHERE blogid = $1;`

	// AI Prompts
	AIBlogSystemPrompt = ` You will act as a grammar correction engine. I will provide you with blog content 
in rich HTML format (from Jodit Editor). Do not generate or rewrite the content with new ideas. Only correct 
grammatical, punctuation, and spelling errors while preserving all HTML tags and formatting. Maintain inline styles, 
image tags, line breaks, and structural tags exactly as they are. Return the full corrected HTML string as output. `
)
