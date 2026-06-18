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

	// Blog retrieval queries
	SelectAllBlogsQuery = `
		SELECT id, title, description, blogcontent, image, category, author, created_at
		FROM blogs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;`

	SelectBlogsByCategoryQuery = `
		SELECT id, title, description, blogcontent, image, category, author, created_at
		FROM blogs
		WHERE category = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;`

	SelectBlogsBySearchQuery = `
		SELECT id, title, description, blogcontent, image, category, author, created_at
		FROM blogs
		WHERE (title ILIKE $1 OR description ILIKE $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;`

	SelectBlogsBySearchAndCategoryQuery = `
		SELECT id, title, description, blogcontent, image, category, author, created_at
		FROM blogs
		WHERE (title ILIKE $1 OR description ILIKE $1) AND category = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4;`

	SelectSavedBlogsDetailedByUserIDQuery = `
		SELECT b.id, b.title, b.description, b.blogcontent, b.image, b.category, b.author, b.created_at
		FROM savedblogs sb
		JOIN blogs b ON sb.blogid = CAST(b.id AS VARCHAR)
		WHERE sb.userid = $1
		ORDER BY sb.created_at DESC;`

	SelectBlogByIDQuery = `
		SELECT id, title, description, blogcontent, image, category, author, created_at
		FROM blogs
		WHERE id = $1;`

	// Comment queries
	InsertCommentQuery = `
		INSERT INTO comments (comment, blogid, userid, username)
		VALUES ($1, $2, $3, $4)
		RETURNING id, comment, userid, username, blogid, created_at;`

	SelectCommentsByBlogIDQuery = `
		SELECT id, comment, userid, username, blogid, created_at
		FROM comments
		WHERE blogid = $1
		ORDER BY created_at DESC;`

	SelectCommentByIDQuery = `
		SELECT id, comment, userid, username, blogid, created_at
		FROM comments
		WHERE id = $1;`

	DeleteCommentByIDQuery = `
		DELETE FROM comments
		WHERE id = $1;`

	// Saved blog queries
	SelectSavedBlogQuery = `
		SELECT id, userid, blogid, created_at
		FROM savedblogs
		WHERE userid = $1 AND blogid = $2;`

	InsertSavedBlogQuery = `
		INSERT INTO savedblogs (userid, blogid)
		VALUES ($1, $2)
		RETURNING id, userid, blogid, created_at;`

	DeleteSavedBlogQuery = `
		DELETE FROM savedblogs
		WHERE userid = $1 AND blogid = $2;`

	SelectSavedBlogsByUserIDQuery = `
		SELECT id, userid, blogid, created_at
		FROM savedblogs
		WHERE userid = $1
		ORDER BY created_at DESC;`
)
