package session

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID string, duration time.Duration) (string, error) {
	token := uuid.New().String()
	createdAt := time.Now()
	expiresAt := createdAt.Add(duration)
	createdAtStr := createdAt.Format("2006-01-02 15:04:05")
	expiresAtStr := expiresAt.Format("2006-01-02 15:04:05")
	log.Printf("Paramètres de la session :  userID=%s, token=%s, createdAt=%s, expiresAt=%s\n",
		 userID, token, createdAtStr, expiresAtStr)

	_, err := db.Exec(`
		INSERT INTO sessions ( user_id, token, created_at, expires_at)
		VALUES ( ?, ?, ?, ?);
	`,  userID, token, createdAt, expiresAt)
	if err != nil {
		return "", fmt.Errorf("erreur d'insertion dans la base de données : %w", err)
	}

	return token, nil
}

func IsValidSession(db *sql.DB, token string) (bool, error) {
	var expiresAt time.Time

	err := db.QueryRow(`
		SELECT expires_at FROM sessions WHERE token = ?;
	`, token).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return time.Now().Before(expiresAt), nil
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(`
		DELETE FROM sessions WHERE token = ?;
	`, token)
	return err
}
func CleanupExpiredSessions(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM sessions WHERE expires_at <= ?;
	`, time.Now())
	return err
}
