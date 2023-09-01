package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/epaitoo/remindme/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	reminders, err := app.reminders.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Reminders = reminders

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}


func (app *application) reminderView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	reminder, err := app.reminders.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Reminder = reminder

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}


func (app *application) reminderCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	item :=  "O snail"
	description := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	due := 7


	id, err := app.reminders.Insert(item, description, due)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the newly created Reminder
	http.Redirect(w, r, fmt.Sprintf("/reminder/view?id=%d", id), http.StatusSeeOther)
}


func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData)  {
	// retrieve the page
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write out the provided HTTP status code
	w.WriteHeader(status)

	buf.WriteTo(w)

	
}