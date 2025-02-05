package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/luigigil/contact-app/domain/contact"
	"github.com/luigigil/contact-app/internal/flash"
	views "github.com/luigigil/contact-app/templates"
)

func Render(w http.ResponseWriter, templates *template.Template, name string, data interface{}) {
	tmpl := template.Must(templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(views.TmplFS, name))

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := chi.NewRouter()

	funcMap := template.FuncMap{
		"sum": func(i, j int) int {
			return i + j
		},
		"subtract": func(i, j int) int {
			return i - j
		},
	}

	templates := template.Must(
		template.New("").Funcs(funcMap).ParseFS(
			views.TmplFS,
			"0-layout.html",
			"index.html",
			"new.html",
		))

	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
	r.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
		page := 1
		if r.URL.Query().Has("page") {
			p, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				p = 1
			}
			page = p
		}

		var contacts []contact.Contact
		hasNext := false
		if r.URL.Query().Has("q") {
			contacts = contact.Search(r.URL.Query().Get("q"))
		} else {
			contacts, hasNext = contact.All(page)
		}

		messages, _ := flash.GetFlash(w, r)

		Render(w, templates, "index.html", map[string]interface{}{
			"Contacts": contacts,
			"Messages": messages,
			"Page":     page,
			"HasNext":  hasNext,
		})
	})
	r.Get("/contacts/new", func(w http.ResponseWriter, r *http.Request) {
		Render(w, templates, "new.html", map[string]interface{}{
			"Contact": contact.Contact{},
		})
	})
	r.Post("/contacts/new", func(w http.ResponseWriter, r *http.Request) {
		c := contact.Contact{
			ID:    0,
			First: r.FormValue("first_name"),
			Last:  r.FormValue("last_name"),
			Phone: r.FormValue("phone"),
			Email: r.FormValue("email"),
		}

		if contact.Save(c) {
			flash.SetFlash(w, r, []byte("Created New Contact!"))
			http.Redirect(w, r, "/contacts", http.StatusSeeOther)
			return
		}

		Render(w, templates, "new.html", map[string]interface{}{
			"Contact": c,
		})
	})
	r.Get("/contacts/{contactID}", func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		if contactID == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(contactID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		c, err := contact.Find(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		Render(w, templates, "show.html", map[string]interface{}{
			"Contact": c,
		})
	})
	r.Get("/contacts/{contactID}/edit", func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		if contactID == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(contactID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		c, err := contact.Find(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		Render(w, templates, "edit.html", map[string]interface{}{
			"Contact": c,
		})
	})
	r.Post("/contacts/{contactID}/edit", func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		if contactID == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(contactID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		c, err := contact.Find(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		c.First = r.FormValue("first_name")
		c.Last = r.FormValue("last_name")
		c.Phone = r.FormValue("phone")
		c.Email = r.FormValue("email")

		if contact.Save(c) {
			flash.SetFlash(w, r, []byte("Updated Contact!"))
			http.Redirect(w, r, "/contacts", http.StatusSeeOther)
			return
		}

		Render(w, templates, "edit.html", map[string]interface{}{
			"Contact": c,
		})
	})
	r.Delete("/contacts/{contactID}", func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		if contactID == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(contactID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		err = contact.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		flash.SetFlash(w, r, []byte("Deleted Contact!"))
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
	r.Get("/contacts/{contactID}/email", func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		if contactID == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(contactID)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}
		c, err := contact.Find(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		c.Email = r.FormValue("email")
		fmt.Println(c.Email)
		contact.Validate(c)

		fmt.Println(c.Errors["email"])
		w.Write([]byte(c.Errors["email"]))
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalf("failed to initialize server: %s", err)
	}
}
