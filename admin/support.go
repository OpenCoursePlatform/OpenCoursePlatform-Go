package admin

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
	"github.com/gorilla/mux"
)

/*
TicketData ...
*/
type TicketData struct {
	Title     string
	User      string
	Solved    bool
	Slug      string
	Responses []string
}

// TicketsPageData ...
type TicketsPageData struct {
	Title   string
	Paths   []helpers.Path
	Tickets []TicketData
}

/*
GetTickets ...
*/
func GetTickets(db *sql.DB) ([]TicketData, error) {
	var tickets []TicketData

	rows, err := db.Query(`
	SELECT tickets.topic, users.username, tickets.solved, tickets.slug
	FROM tickets
	INNER JOIN users
		ON tickets.user_id = users.id`)
	if err != nil {
		return tickets, err
	}
	defer rows.Close()

	for rows.Next() {
		var ticket TicketData
		err = rows.Scan(&ticket.Title, &ticket.User, &ticket.Solved, &ticket.Slug)
		if err != nil {
			return tickets, err
		}
		tickets = append(tickets, ticket)
	}
	err = rows.Err()
	if err != nil {
		return tickets, err
	}
	return tickets, nil
}

/*
GetTicket ...
*/
func GetTicket(db *sql.DB, slug string) (TicketData, error) {
	var ticket TicketData

	query := `SELECT tickets.topic, users.username, tickets.solved, tickets.slug
	FROM tickets
	INNER JOIN users
		ON tickets.user_id = users.id
	WHERE slug = ?`
	err := db.QueryRow(query, slug).Scan(&ticket.Title, &ticket.User, &ticket.Solved, &ticket.Slug)
	if err != nil {
		return ticket, err
	}

	var responses []string

	rows, err := db.Query(`
	SELECT ticket_responses.text
	FROM ticket_responses
	INNER JOIN tickets
		ON ticket_responses.ticket_id = tickets.id
	WHERE tickets.slug = ?
	`, slug)
	if err != nil {
		return ticket, err
	}
	defer rows.Close()

	for rows.Next() {
		var response string
		err = rows.Scan(&response)
		if err != nil {
			return ticket, err
		}
		responses = append(responses, response)
	}
	err = rows.Err()
	if err != nil {
		return ticket, err
	}
	ticket.Responses = responses
	return ticket, nil
}

/*
Tickets ...
*/
func Tickets(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	tickets, err := GetTickets(db)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &TicketsPageData{"Tickets", []helpers.Path{{Name: "Admin", Link: "/admin"}}, tickets}
	helpers.RenderTemplate(r, w, "admin_tickets", p)
}

// TicketPageData ...
type TicketPageData struct {
	Title  string
	Paths  []helpers.Path
	Ticket TicketData
}

/*
Ticket ...
*/
func Ticket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["ticket"]

	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}
	defer db.Close()

	ticket, err := GetTicket(db, slug)
	if err != nil {
		helpers.HandleError(err)
	}

	p := &TicketPageData{ticket.Title, []helpers.Path{{Name: "Admin", Link: "/admin"}}, ticket}
	helpers.RenderTemplate(r, w, "admin_ticket", p)
}

/*
UpdateTicket ...
*/
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["ticket"]

	http.Redirect(w, r, "/admin/tickets/"+slug, http.StatusSeeOther)
}
