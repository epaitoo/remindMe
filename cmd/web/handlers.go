package main

import (
	"errors"
	"fmt"
	// "html/template"
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

	for _, reminder := range reminders {
		fmt.Fprintf(w, "%+v\n", reminder)
	}

	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }
	
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// }
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


	fmt.Fprintf(w, "%+v", reminder)

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
