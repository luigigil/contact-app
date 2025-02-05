package contact

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Contact struct {
	ID     int
	First  string
	Last   string
	Phone  string
	Email  string
	Errors map[string]string
}

var contacts []Contact

func init() {
	Load()
}

func Load() {
	data, err := os.ReadFile("contacts.json")
	if err != nil {
		log.Fatalf("couldn't read contacts file")
	}

	err = json.Unmarshal(data, &contacts)
	if err != nil {
		log.Fatalf("couldn't unmarshal contacts")
	}
}

func All() []Contact {
	return contacts
}

func Search(query string) []Contact {
	var found []Contact
	for _, c := range contacts {
		query = strings.ToLower(query)
		if strings.Contains(strings.ToLower(c.First), query) ||
			strings.Contains(strings.ToLower(c.Last), query) ||
			strings.Contains(strings.ToLower(c.Email), query) {
			found = append(found, c)
		}
	}
	return found
}

func Save(contact Contact) bool {
	// Initialize errors map if it doesn't exist
	if contact.Errors == nil {
		contact.Errors = make(map[string]string)
	}

	// Basic validation
	if strings.TrimSpace(contact.First) == "" {
		contact.Errors["first"] = "First name is required"
	}
	if strings.TrimSpace(contact.Email) == "" {
		contact.Errors["email"] = "Email is required"
	} else if !strings.Contains(contact.Email, "@") {
		contact.Errors["email"] = "Invalid email format"
	}
	if strings.TrimSpace(contact.Phone) != "" && len(contact.Phone) < 10 {
		contact.Errors["phone"] = "Phone number must be at least 10 digits"
	}

	// If there are any errors, return false
	if len(contact.Errors) > 0 {
		return false
	}

	// Generate new ID if it's a new contact
	if contact.ID == 0 {
		maxID := 0
		for _, c := range contacts {
			if c.ID > maxID {
				maxID = c.ID
			}
		}
		contact.ID = maxID + 1
	}

	// Add or update the contact
	found := false
	for i, c := range contacts {
		if c.ID == contact.ID {
			contacts[i] = contact
			found = true
			break
		}
	}
	if !found {
		contacts = append(contacts, contact)
	}

	return true
}

func Find(contactID int) (Contact, error) {
	for _, c := range contacts {
		if c.ID == contactID {
			return c, nil
		}
	}
	return Contact{}, fmt.Errorf("didn't find contact")
}

func Delete(contactID int) error {
	for i, c := range contacts {
		if c.ID == contactID {
			contacts = slices.Delete(contacts, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("didn't find contact")
}
