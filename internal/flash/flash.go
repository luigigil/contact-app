package flash

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

const cookieName = "messages"

func SetFlash(w http.ResponseWriter, r *http.Request, value []byte) {
	messages := getMessages(r)
	messages = append(messages, value)

	data, err := json.Marshal(messages)
	if err != nil {
		return
	}

	c := &http.Cookie{
		Name:    cookieName,
		Value:   encode(data),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour), // expires in 24 hours if not consumed
	}
	http.SetCookie(w, c)
}

func GetFlash(w http.ResponseWriter, r *http.Request) ([]string, error) {
	messages := getMessages(r)
	if len(messages) == 0 {
		return nil, nil
	}

	// Clear the cookie
	dc := &http.Cookie{
		Name:    cookieName,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Path:    "/",
	}
	http.SetCookie(w, dc)

	// Convert [][]byte to []string
	strMessages := make([]string, len(messages))
	for i, msg := range messages {
		strMessages[i] = string(msg)
	}

	return strMessages, nil
}

// getMessages retrieves the current messages from the cookie
func getMessages(r *http.Request) [][]byte {
	if r == nil {
		return [][]byte{}
	}

	c, err := r.Cookie(cookieName)
	if err != nil {
		return [][]byte{}
	}

	value, err := decode(c.Value)
	if err != nil {
		return [][]byte{}
	}

	var messages [][]byte
	if err := json.Unmarshal(value, &messages); err != nil {
		return [][]byte{}
	}

	return messages
}

// -------------------------

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
