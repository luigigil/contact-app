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

const PAGE_SIZE = 5

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

func All(page int) ([]Contact, bool) {
	index := (page - 1) * PAGE_SIZE
	hasNext := false

	if index >= len(contacts) {
		return []Contact{}, hasNext
	}

	hasNext = index+PAGE_SIZE < len(contacts)

	return contacts[index : index+PAGE_SIZE], hasNext
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
	contact = Validate(contact)

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
			c.Errors = make(map[string]string)
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

func Validate(contact Contact) Contact {
	// Initialize errors map if it doesn't exist
	if contact.Errors == nil {
		contact.Errors = make(map[string]string)
	}

	// First name validation
	if strings.TrimSpace(contact.First) == "" {
		contact.Errors["first"] = "First name is required"
	} else if len(contact.First) < 2 {
		contact.Errors["first"] = "First name must be at least 2 characters"
	}

	// Last name validation
	if strings.TrimSpace(contact.Last) == "" {
		contact.Errors["last"] = "Last name is required"
	} else if len(contact.Last) < 2 {
		contact.Errors["last"] = "Last name must be at least 2 characters"
	}

	// Email validation
	if strings.TrimSpace(contact.Email) == "" {
		contact.Errors["email"] = "Email is required"
	} else if !strings.Contains(contact.Email, "@") || !strings.Contains(contact.Email, ".") {
		contact.Errors["email"] = "Invalid email format"
	}
	emailExists := false
	for _, c := range contacts {
		if c.Email == contact.Email {
			emailExists = true
		}
	}
	if emailExists {
		contact.Errors["email"] = "Email already exists"
	}

	// Phone validation (optional field)
	if phone := strings.TrimSpace(contact.Phone); phone != "" {
		// Remove any non-digit characters for validation
		digits := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, phone)

		if len(digits) < 10 {
			contact.Errors["phone"] = "Phone number must have at least 10 digits"
		} else if len(digits) > 15 {
			contact.Errors["phone"] = "Phone number cannot exceed 15 digits"
		}
	}

	return contact
}
