package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func SendMail(toEmail string, subject string, body string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO API KEY isn't set yet.")
	}

	payload, _ := json.Marshal(map[string]any{
		"sender":map[string]string{
			"email": os.Getenv("BREVO_SENDER_EMAIL"),
			"name": os.Getenv("BREVO_SENDER_NAME"),
		},
		"to": []map[string]string{{"email":toEmail}},
		"subject": subject,
		"htmlContent": body,
	})

	req, _ := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(payload))
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Brevo: status %d - %s", res.StatusCode, b)
	}

	return nil
}