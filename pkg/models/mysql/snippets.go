// snippets.go will contains the code specifically for working with the snippets in our MySQL database
// Here we're going to define a new SnippetModel type and implement some methods on it to access and manipulate the database
package mysql

import (
	"database/sql"
	"errors"

	// Import statement should look like this: "{your-module-path}/pkg/models"

	"roman-hds/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	//Write the SQL statement wewant to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backqoutes instead of normal double qoutes)
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"

	// Use the Exec() method on the emdedded connection pool to execute the statement.
	// The first parameter is the SQL statement, followed by the title, content and expiry values
	// fir the placeholder parameters. This method returns a sql.Result object, which contains some basic information
	// about what happened when the statement was executed
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertID() method on the result object to get the ID of our newly inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type before returning
	return int(id), nil
}

// This will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}
	// errors from DB.QueryRow() are deferred until Scan() is called.
	err := m.DB.QueryRow(`SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a sql. ErrNoRows error. Use the errors.Is() function check for that error
		// specifically, and return our own models.ErrNoRecord error instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// Write the SQL statement we want to execute
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	// Use the Query() method on the connection pool to execute our SQL statement. This returns a sql.Rows resulset containing the results of our query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure thje sql.Rows() resultset is always properly closed before the Latest() method returns.
	// This defer statement should come *after* you check for an error from the Query() method. Otherwise, if Query() returns an error, you'll get a panic trying to close a nil resultset
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects
	snippets := []*models.Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This prepares the first (and then each subsequent) row to be acted on by the rows.Scan() method.
	// If iteration over all the rows completes, then the resultset automatically closes itself and frees-up the underlying database connection
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct
		s := &models.Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy data into, and the number of arguments must be exactly the same as the number of columns returned by your statement
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any error that we encountered during the iteration.
	// It's important to call this - don't assume that a successfull iteration was completed over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice
	return snippets, nil
}

// Important: Closing the resultset with defer rows.Close() is critical here. As long as a resultset is open it will keep the underlying database connection open... so if something goes wrong
// in this method and the resultset isn't closede, it can rapidly lead to all the connection in your pool being used up
