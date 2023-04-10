package models

import "github.com/google/uuid"

type Contacts struct {
	ID                int       `db:"id"`
	ProfileID         uuid.UUID `db:"id_profile"`
	PushNotifications bool      `db:"push_notifications"`
	Email             string    `db:"email"`
	EmailSubscription bool      `db:"email_subscription"`
	HashID            string    `db:"hash_id"`
	ChatID            string    `db:"chat_id"`
}
