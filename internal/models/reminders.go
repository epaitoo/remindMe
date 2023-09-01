package models

import (
	"database/sql"
	"errors"
	"time"
)

type Reminder struct {
	ID          int
	Item        string
	Description string
	Created     time.Time
	Due         time.Time
}

// Define a ReminderModel type which wraps a sql.DB connection pool
type ReminderModel struct {
	DB *sql.DB
}


// This will insert a new reminder into the database.
func (m *ReminderModel) Insert(item string, description string, due int) (int, error) {
	stmt := `INSERT INTO reminders (item, description, created, due)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, item, description, due)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return a reminder based on the ID
func (m *ReminderModel) Get(id int) (*Reminder, error) {
	stmt := `SELECT id, item, description, created, due FROM reminders WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	r := &Reminder{}

	err := row.Scan(&r.ID, &r.Item, &r.Description, &r.Created, &r.Due)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	

	return r, nil
}

// Return list of reminders
func (m *ReminderModel) Latest() ([]*Reminder, error) {
	stmt := `SELECT id, item, description, created, due FROM reminders where due > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reminders := []*Reminder{}

	for rows.Next() {
		r := &Reminder{}
		err := rows.Scan(&r.ID, &r.Item, &r.Description, &r.Created, &r.Due)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}


	return reminders, nil
}


