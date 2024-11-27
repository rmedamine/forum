package database

import (
	"fmt"
	"strconv"
	"time"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Date     time.Time
	Like     int
	Deslike  int
	User     string
	Category string
	Comm     []Comment
}
type Comment struct {
	Comment string
	User    string
}

func InsertPost(title, content, email, categories string) error {
	/* if(len(content)> 500){

		return <max len error>

	}
	*/

	id := 0
	category_id := 0
	erre := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id)
	_ = erre
	erre = db.QueryRow("SELECT category_id FROM categories WHERE name = ?", categories).Scan(&category_id)
	// fmt.Println(categories)
	_ = erre
	_, err := db.Exec("INSERT INTO posts (title, content, createdAt,user_id, category_id) VALUES (?, ?, datetime('now'),?,?)", title, content, id, category_id)
	return err
}

/*
	func InsertComment(comment,email string) error{
		id := 0
		erre := db.QueryRow("SELECT FROM users WHERE email = ?",email).Scan(&id)
		return err
	}
*/

func CountPost(limit int) bool {
	var l string
	err := db.QueryRow("SELECT COUNT(post_id) FROM posts").Scan(&l)
	if err != nil {
		fmt.Println("Error querying row:", err)
		return false
	}
	// fmt.Println(l)

	max, err := strconv.Atoi(l)
	if err != nil {
		fmt.Println("Error converting count to integer:", err)
		return false
	}

	return max > limit
}

func GetPosts(catigorie string, limit int) ([]Post, error) {
	if limit < 0 {
		limit = 0
	}
	var query string
	query = `SELECT 
    posts.post_id,
    posts.title,
    posts.content,
    posts.createdAt, 
    users.name,
    categories.name AS category,
    COALESCE(SUM(CASE WHEN likes.is_like = 1 THEN 1 ELSE 0 END), 0) AS like_count,
    COALESCE(SUM(CASE WHEN likes.is_like = 2 THEN 1 ELSE 0 END), 0) AS dislike_count
FROM 
    posts
LEFT JOIN users ON users.user_id = posts.user_id
LEFT JOIN categories ON categories.category_id = posts.category_id
LEFT JOIN likes ON likes.post_id = posts.post_id
GROUP BY 
    posts.post_id, posts.title, posts.content, posts.createdAt, users.name, categories.name
`

	if catigorie != "" {
		query += " WHERE category = ?"
	}
	query += " ORDER BY posts.createdAt DESC LIMIT 5 OFFSET " + strconv.Itoa(limit)
	rows, err := db.Query(query, catigorie)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Date, &post.User, &post.Category, &post.Like, &post.Deslike); err != nil {
			return nil, err
		}
		/*rows1, err := db.Query(`SELECT COUNT(likes.post_id) FROM likes WHERE likes.post_id = ? AND is_like = 1`, post.ID)
		if err != nil {
			return nil, err
		}
		defer rows1.Close()
		if rows1.Next() {
			rows1.Scan(&post.Like)
		}
			rows2, err := db.Query(`SELECT COUNT(likes.post_id) FROM likes WHERE likes.post_id = ? AND is_like = 2`, post.ID)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		if rows2.Next() {
			rows2.Scan(&post.Deslike)
		}*/

		posts = append(posts, post)
	}
	return posts, nil
}

func GetComment(id string) (any, error) {
	rows, err := db.Query(`SELECT content, users.name FROM comments LEFT Join users ON users.user_id=comments.user_id WHERE  post_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows1, err := db.Query(`SELECT posts.title, posts.content, posts.createdAt, users.name ,categories.name AS category_name
                    FROM posts 
                    LEFT JOIN users ON posts.user_id = users.user_id 
					LEFT JOIN categories ON categories.category_id = posts.category_id
                    WHERE posts.post_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()
	var post Post
	if rows1.Next() {
		if err := rows1.Scan(&post.Title, &post.Content, &post.Date, &post.User, &post.Category); err != nil {
			return nil, err
		}
	}
	// fmt.Println(post)
	var Comments []Comment
	for rows.Next() {
		var Comment Comment
		if err := rows.Scan(&Comment.Comment, &Comment.User); err != nil {
			return nil, err
		}
		Comments = append(Comments, Comment)
	}
	post.Comm = Comments
	return post, nil
}

func InsertLike(id, email string, is_like bool) error {
	var id_user int
	var checkrow int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id_user)
	if err != nil {
		return err
	}
	pre, err := db.Prepare("SELECT COUNT(like_id) FROM likes WHERE post_id = ? AND user_id = ? ")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer pre.Close()
	err = pre.QueryRow(id, id_user).Scan(&checkrow)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if checkrow == 0 && is_like {
		if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 1); err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else if checkrow == 0 && !is_like {
		if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 2); err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else if checkrow != 0 && is_like {
		var reaction int
		pre, err := db.Prepare("SELECT is_like FROM likes WHERE post_id = ? AND user_id = ? ")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer pre.Close()
		err = pre.QueryRow(id, id_user).Scan(&reaction)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if _, err := db.Exec(`DELETE FROM Likes WHERE post_id = ? AND user_id = ? `, id, id_user); err != nil {
			fmt.Println(err.Error())
			return err
		}
		if reaction == 2 {
			if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 1); err != nil {
				fmt.Println(err.Error())
				return err
			}
		}

	} else if checkrow != 0 && !is_like {

		var reaction int
		pre, err := db.Prepare("SELECT is_like FROM likes WHERE post_id = ? AND user_id = ? ")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer pre.Close()
		err = pre.QueryRow(id, id_user).Scan(&reaction)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if _, err := db.Exec(`DELETE FROM Likes WHERE post_id = ? AND user_id = ? `, id, id_user); err != nil {
			fmt.Println(err.Error())
			return err
		}
		if reaction == 1 {
			if _, err = db.Exec(`INSERT INTO likes (post_id,user_id,is_like,type) VALUES (?,?,?,'post')`, id, id_user, 2); err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}
	return nil
}
