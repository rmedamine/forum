package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forim/bcryptp"
	"forim/database"
	sessions "forim/session"
)

var limit = 0

func GetHome(w http.ResponseWriter, r *http.Request) {
	catigorie := r.FormValue("category")
	action := r.FormValue("Next")
	if action != "" && database.CountPost(limit+1) {
		limit += 5
	}
	fmt.Print(database.CountPost(limit + 1))
	action = r.FormValue("Back")
	if action != "" && limit != 0 {
		limit -= 5
	}
	posts, err := database.GetPosts(catigorie, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, "./assets/templates/post.html", posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if title != "" && content != "" && category != "" {
			if len(title) < 5 || len(title) > 50 {
				http.Error(w, "title is too long or too short", http.StatusBadRequest)
				return
			}
			if len(content) < 10 || len(content) > 500 {
				http.Error(w, "content is too long or too short", http.StatusBadRequest)
				return
			}
			if err := database.InsertPost(title, content, cookie.Value, category); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/post", http.StatusSeeOther)
			return
		}
	}

	RenderTemplate(w, "./assets/templates/post.create.page.html", nil)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	id_post := r.FormValue("id-post")
	comments, err := database.GetComment(id_post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("comments : ", comments)
	RenderTemplate(w, "./assets/templates/comment.html", comments)
}

func NewComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id_post := r.FormValue("id-post")
	comment := r.FormValue("comment")
	if len(comment) > 200 {
		http.Error(w, "comment is too long", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if cookie == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if comment != "" {
		if err := database.Createcomment(comment, id_post, cookie.Value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/post", http.StatusSeeOther)
}

func Like_post(w http.ResponseWriter, r *http.Request) {
	like := r.FormValue("like_post")
	deslike := r.FormValue("deslike_post")
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if like != "" {
		err = database.InsertLike(like, cookie.Value, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err = database.InsertLike(deslike, cookie.Value, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/post", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")
	success, userId, err := database.GetLogin(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if success {
		duration := 24 * time.Hour
		err = sessions.CleanupExpiredSessions(database.GetDB())
		if err != nil {
			log.Fatal("Erreur lors du nettoyage des sessions:", err)
		}

		token, err := sessions.CreateSession(database.GetDB(), fmt.Sprintf("%d", userId), duration)
		if err != nil {
			log.Printf("Erreur lors de la création de la session  : %v\n", err)
			http.Error(w, "Erreur lors de la création de la session", http.StatusInternalServerError)
			return
		}
		log.Printf("Paramètres envoyés : userId=%d, duration=%s\n", userId, duration)

		_, err = sessions.IsValidSession(database.GetDB(), token)
		if err != nil {
			log.Fatal("Erreur lors de la validation de session:", err)
		}
		

		catigorie := r.FormValue("category")
		posts, err := database.GetPosts(catigorie, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:  "session",
			Value: token,
		}
		http.SetCookie(w, &cookie)
		RenderTemplate(w, "./assets/templates/post.html", posts)
	} else {
		errorMessage := ""
		if email != "" || password != "" {
			errorMessage = "Password or email not working"
		}
		RenderTemplate(w, "./assets/templates/index.html", errorMessage)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		name := r.FormValue("username")
		p, err := bcryptp.HashPassword(password)
		if email == "" || p == "" || name == "" {
			RenderTemplate(w, "./assets/templates/register.html", nil)
			return
		}
		if err != nil {
			RenderTemplate(w, "./assets/templates/register.html", nil)
			return
		}
		if err := database.CreateAcount(name, email, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		RenderTemplate(w, "./assets/templates/register.html", nil)
	}
}
