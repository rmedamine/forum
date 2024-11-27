package creatdatabase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Creatdb() {
	// Open a connection to the SQLite database
	data, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer data.Close()

	// Create tables
	_, err = data.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS categories (
		category_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts (
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		createdAt DATETIME NOT NULL,
		user_id INTEGER,
		category_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER, 
		content TEXT NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE  
	);

	CREATE TABLE IF NOT EXISTS likes (
		like_id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		user_id INTEGER,
		is_like INTEGER,
		type TEXT,
		FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE  
	);
	CREATE TABLE IF NOT EXISTS sessions (
	session_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE  
	);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Insert users
	_, err = data.Exec(`
	INSERT INTO users (name, email, password) VALUES
	('Alice', 'alice@example.com', 'password123'),
	('Bob', 'bob@example.com', 'password456'),
	('Charlie', 'charlie@example.com', 'password789');
	`)
	if err != nil {
		log.Fatalf("Failed to insert users ici: %v", err)
	}

	// Insert categories
	_, err = data.Exec(`
	INSERT INTO categories (name) VALUES
	('news'),
	('tech'),
	('Lifestyle'),
	('education'),
	('health'),
	('entertainment');
	`)
	if err != nil {
		log.Fatalf("Failed to insert categories: %v", err)
	}

	// Insert posts
	_, err = data.Exec(`
	INSERT INTO posts (title, content, createdAt, user_id, category_id) VALUES
	('The Future of Tech', 'This is a post about technology trends.', '2024-11-05 10:00:00', 1, 1),
	('Healthy Living', 'Tips for a healthier life.', '2024-11-05 11:00:00', 2, 2),
	('Traveling Tips', 'Best places to visit this summer.', '2024-11-05 12:00:00', 3, 3);
	`)
	if err != nil {
		log.Fatalf("Failed to insert posts: %v", err)
	}

	// Insert comments
	_, err = data.Exec(`
	INSERT INTO comments (post_id, user_id, content) VALUES
	(1, 2, 'Great insights on technology!'),
	(1, 3, 'I love technology discussions.'),
	(2, 1, 'Very helpful tips! Thanks!'),
	(3, 1, 'I can’t wait to travel!'),
	(3, 2, 'Awesome places to consider!');
	`)
	if err != nil {
		log.Fatalf("Failed to insert comments: %v", err)
	}

	// Insert likes
	_, err = data.Exec(`
	INSERT INTO likes (post_id, is_like) VALUES
	(1, 10),
	(2, 5),
	(3, 8);
	`)
	if err != nil {
		log.Fatalf("Failed to insert likes: %v", err)
	}
	//session
	_, err = data.Exec(`
		INSERT INTO sessions (session_id, user_id, token, created_at, expires_at) VALUES
		 ('1', 'test_user', 'test_token', DATETIME('now'), DATETIME('now', '+1 hour'));
	`)

	if err != nil {
		log.Fatalf("Erreur lors de l'insertion de données fixes : %v\n", err)

	}

	fmt.Println("Data inserted successfully!")
}
