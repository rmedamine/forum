package database

func CreateAcount(name, email, passeord string) error {
	_, err := db.Exec("INSERT INTO  users (name,email,password) VALUES ($1,$2,$3)", name, email, passeord)
	return err
}

func Createcomment(comment, post_id, email string) error {
	id := 0
	erre := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&id)
	_ = erre
	_, err := db.Exec("INSERT INTO comments (post_id,content,user_id) VALUES ($1,$2,$3)", post_id, comment, id)
	return err
}


