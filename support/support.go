package support

import (
	"database/sql"
	"net/http"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/authentication"
	"github.com/gorilla/mux"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

/*
TicketData struct is used
for storing data related to
Customer Support Tickets.
*/
type TicketData struct {
	Title  string
	Slug   string
	Solved bool
}

/*
PageData struct is used for the
data being rendered and should
not be used in other structs.
*/
type PageData struct {
	Title   string
	Paths   []helpers.Path
	Tickets []TicketData
}

/*
Page struct is used for generic
data to be rendered in HTML and should
not be used in other structs.
*/
type Page struct {
	Title string
	Paths []helpers.Path
}

// GetTicketsByUsername gets tickets associated with the username variable sent in.
func GetTicketsByUsername(db *sql.DB, username string) ([]TicketData, error) {
	query := `
	SELECT tickets.topic, tickets.slug, tickets.solved
	FROM tickets
	INNER JOIN users
		ON users.id = tickets.user_id
	WHERE users.username = ?`
	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []TicketData
	for rows.Next() {
		var ticket TicketData
		err = rows.Scan(&ticket.Title, &ticket.Slug, &ticket.Solved)
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
Support endpoint showing the user
active and finished tickets in a
table.
*/
func Support(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	tickets, err := GetTicketsByUsername(db, username)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	supportPage := PageData{"Support", []helpers.Path{{Name: "Support", Link: "/support"}}, tickets}

	helpers.RenderTemplate(r, w, "support", supportPage)
}

/*
NewTicket endpoint rendering page to create
a new customer support ticket. Has
no other dynamic data than the "standard"
data such as username.
*/
func NewTicket(w http.ResponseWriter, r *http.Request) {
	supportPage := Page{"Support", []helpers.Path{{Name: "Support", Link: "/support"}}}

	helpers.RenderTemplate(r, w, "support_new_ticket", supportPage)
}

/*
InsertNewTicketInDB only uses database
functions and inserts the new ticket into
the database as well as the first customer
message as a ticket_response.
Returns the slug and potential error.
*/
func InsertNewTicketInDB(db *sql.DB, username, title, text string) (string, error) {
	slug := helpers.GenerateSlug(title)
	userID, err := authentication.GetUserIDFromUsername(db, username)
	if err != nil {
		return "", err
	}
	insForm, err := db.Prepare("INSERT INTO tickets (user_id, topic, slug) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	res, err := insForm.Exec(userID, title, slug)
	if err != nil {
		return "", err
	}
	insForm, err = db.Prepare("INSERT INTO ticket_responses (ticket_id, user_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return "", err
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	_, err = insForm.Exec(lastInsertedID, userID, text)
	if err != nil {
		return "", err
	}
	return slug, nil
}

/*
InsertNewTicket endpoint parses the data from the
user and inserts the new ticket in the database.
*/
func InsertNewTicket(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	title := r.FormValue("title")
	text := r.FormValue("text")
	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	slug, err := InsertNewTicketInDB(db, username, title, text)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/support/"+slug, http.StatusSeeOther)
}

/*
GetTicketBySlug gets individual ticket by its slug.
Should only interact with the database.
*/
func GetTicketBySlug(db *sql.DB, slug, username string) (int, TicketData, error) {
	var ticket TicketData
	var ticketID int
	query := `
	SELECT tickets.id, tickets.topic, tickets.slug, tickets.solved
	FROM tickets
	INNER JOIN users
		ON users.id = tickets.user_id
	WHERE
		users.username = ?
	AND tickets.slug = ?`
	err := db.QueryRow(query, username, slug).Scan(&ticketID, &ticket.Title, &ticket.Slug, &ticket.Solved)
	if err != nil {
		return 0, ticket, err
	}
	return ticketID, ticket, nil
}

/*
GetTicketResponsesByTicketID gets all of the responses
for a ticket by its ticket id.
Should only interact with the database.
*/
func GetTicketResponsesByTicketID(db *sql.DB, ticketID int) ([]string, error) {
	query := `
	SELECT text
	FROM ticket_responses
	WHERE ticket_id = ?`
	rows, err := db.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []string
	for rows.Next() {
		var response string
		err = rows.Scan(&response)
		if err != nil {
			return responses, err
		}
		responses = append(responses, response)
	}
	err = rows.Err()
	if err != nil {
		return responses, err
	}
	return responses, nil
}

/*
TicketPage struct is used for rendering the
data for an individual ticket with a ticket
and all of its responses.
*/
type TicketPage struct {
	Title     string
	Paths     []helpers.Path
	Ticket    TicketData
	Responses []string
}

/*
Ticket endpoint returns the data for
the ticket to the user from the ticket
slug. I.e. what is sent after /support/{ticket}
*/
func Ticket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["ticket"]
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	ticketID, ticket, err := GetTicketBySlug(db, slug, username)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	responses, err := GetTicketResponsesByTicketID(db, ticketID)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	supportPage := TicketPage{ticket.Title, []helpers.Path{{Name: "Support", Link: "/support"}, {Name: "Ticket", Link: "/support"}}, ticket, responses}

	helpers.RenderTemplate(r, w, "support_view_ticket", supportPage)
}

/*
InsertNewTicketResponseInDB inserts
a new ticket response into the database
based on the userID and ticketID.
*/
func InsertNewTicketResponseInDB(db *sql.DB, ticketID, userID int, text string) error {
	insForm, err := db.Prepare("INSERT INTO ticket_responses (ticket_id, user_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = insForm.Exec(ticketID, userID, text)
	if err != nil {
		return err
	}
	return nil
}

/*
InsertNewTicketResponse endpoint parses the data
when responding to a ticket and inserts the new
data into the database.
*/
func InsertNewTicketResponse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["ticket"]
	text := r.FormValue("text")
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	defer db.Close()

	username, err := helpers.GetUsernameFromRequest(r)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	userID, err := authentication.GetUserIDFromUsername(db, username)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	ticketID, _, err := GetTicketBySlug(db, slug, username)
	if err != nil {
		helpers.HandleError(err)
		http.NotFound(w, r)
		return
	}
	err = InsertNewTicketResponseInDB(db, ticketID, userID, text)
	if err != nil {
		helpers.HandleError(err)
		helpers.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/support/"+slug, http.StatusSeeOther)
}
