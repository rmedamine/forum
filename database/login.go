package database

import (
	"database/sql"
	"forim/bcryptp"
	sess "forim/session"
	"net/http"
)

func GetLogin(email, password string) (bool, int, error) {
	var userID int
	var passwords string

	err := db.QueryRow("SELECT user_id ,	password FROM users WHERE email =  $1", email).Scan(&userID, &passwords)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	if bcryptp.CheckPasswordHash(password, passwords) {
		return true, userID, nil
	}

	return false, 0, nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = sess.DeleteSession(GetDB(), cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
