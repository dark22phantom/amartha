package notification

import "log"

func (r *Repository) SendEmail(email string, subject string, message string) error {
	log.Printf("Send email to %s with subject %s and message %s", email, subject, message)
	return nil
}
