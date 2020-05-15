package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

/*
	Struct describing a single log line
*/
type (
	LogEntry struct {
		Id				int `json:"id"`
		Host			string `json:"host"`
		Message 		string `json:"message"`
		Source			string `json:"source"`
		Severity		int64 `json:"severity"`
		Timestamp		int64 `json:"timestamp"`
		TimeFormatted	string `json:"timeformatted"`
	}
)

/*
	Struct describing a query for log data
*/
type (
	LogQuery struct {
		SearchFilter	string `json:"search"`
		SourceFilter	string `json:"source"`
		SeverityFilter	int64 `json:"severity"`
		Limit			int64 `json:"limit"`
		HostFilter	 	string `json:"host"`
		UpperBound	 	int64 `json:"upper"`
		LowerBound	 	int64 `json:"lower"`
		Order	 		int64 `json:"order"`
		Offset			int64 `json:"offset"`
	}
)

/*
	Struct describing a client and its metrics
*/
type (
	ClientData struct {
		Host			string `json:"host"`
		Count 			int64 `json:"count"`
		Rate			float64 `json:"rate"`
		LastMessage		string `json:"lastmessage"`
	}
)

/*
	Insert log entry to DB
*/
func writeLogEntry(log *LogEntry) {
	sqlStatement := `
		INSERT INTO log (host, message, source, severity, timestamp, document)
		VALUES ($1, $2, $3, $4, $5, to_tsvector($2))`
	_, db_err = db_con.Exec(sqlStatement, log.Host, log.Message, log.Source, log.Severity, log.Timestamp)
	if db_err != nil {
		panic(db_err)
	}
}

/*
	Get log entry from DB by id
*/
func getLogEntry(id int64) LogEntry {
	row := db_con.QueryRow("SELECT id, host, message, source, severity, timestamp FROM log WHERE id=$1", id)

	var entry LogEntry
	switch err := row.Scan(&entry.Id, &entry.Host, &entry.Message, &entry.Source, &entry.Severity, &entry.Timestamp); err {
		case sql.ErrNoRows:
			return LogEntry {Id: -1, Message: "An entry with the specified index was not found"}
		case nil:
			return entry
		default:
			panic(err)
	}
}

/*
	Perform log entry query with the given LoqQuery struct
*/
func performQuery(query LogQuery) []LogEntry {

	results := []LogEntry{}

	// Set reasonable default limit if none is given
	if (query.Limit <= 0) {
		query.Limit = 50 
	}

	// Calculate constraints, if necessary
	whereComponents := []string{}
	if (query.SearchFilter != "") {
		searchClause := fmt.Sprintf("document @@ to_tsquery('%s')", query.SearchFilter)
		whereComponents = append(whereComponents, searchClause)
	}
	if (query.HostFilter != "") {
		hostClause := fmt.Sprintf("host='%s'", query.HostFilter)
		whereComponents = append(whereComponents, hostClause)
	}
	if (query.SourceFilter != "") {
		sourceClause := fmt.Sprintf("source='%s'", query.SourceFilter)
		whereComponents = append(whereComponents, sourceClause)
	}
	if (query.SeverityFilter > 0) {
		severityClause := fmt.Sprintf("severity=%d", query.SeverityFilter)
		whereComponents = append(whereComponents, severityClause)
	}

	// Allow for negative bounds (n seconds ago)
	now := time.Now().Unix()
	if (query.LowerBound < 0) {
		query.LowerBound += now
	}
	
	// Timestamp constraints (lower can always be 0)
	lowerClause := fmt.Sprintf("timestamp >= %d", query.LowerBound)
	whereComponents = append(whereComponents, lowerClause)
	if (query.UpperBound != 0) {
		if (query.UpperBound < 0) {
			query.UpperBound += now
		}

		upperClause := fmt.Sprintf("timestamp <= %d", query.UpperBound)
		whereComponents = append(whereComponents, upperClause)
	}

	// Build WHERE clause
	var whereClauseBuilder strings.Builder
	if (len(whereComponents) > 0) {
		whereClauseBuilder.WriteString("WHERE ")
	}
	for _, component := range whereComponents {
		whereClauseBuilder.WriteString(component)
		whereClauseBuilder.WriteString(" AND ")
	}
	whereClause := strings.TrimSuffix(whereClauseBuilder.String(), " AND ")

	// Build ORDER BY clause, if necessary
	var orderClause string
	if (query.Order == 0) {
		orderClause = "ORDER BY timestamp DESC"
	} else {
		orderClause = "ORDER BY timestamp ASC"
	}

	// Build and execute query
	queryString := fmt.Sprintf("SELECT id, host, message, source, severity, timestamp FROM log %s %s LIMIT %d OFFSET %d", whereClause, orderClause, query.Limit, query.Offset)
	fmt.Println(queryString)
	rows, _ := db_con.Query(queryString)

	// Loop through rows
	defer rows.Close()
	for rows.Next() {
		var entry LogEntry
		rows.Scan(&entry.Id, &entry.Host, &entry.Message, &entry.Source, &entry.Severity, &entry.Timestamp)
		
		// Convert timestamp
		t := time.Unix(entry.Timestamp, 0)
		entry.TimeFormatted = t.Format(time.RFC3339)

		results = append(results, entry)
	}

	return results
}

/*
	Returns an array of clients and their respective metrics
*/
func performClientQuery() []ClientData {
	results := []ClientData{}
	clientRows, err := db_con.Query("SELECT DISTINCT host FROM log")

	// Loop through rows
	defer clientRows.Close()
	for clientRows.Next() {
		var client ClientData
		err = clientRows.Scan(&client.Host)
		if (err != nil) { break }

		// Get row count
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE host='%s'", client.Host)
		countRows, err := db_con.Query(countQuery)
		defer countRows.Close()
		countRows.Next()
		err = countRows.Scan(&client.Count)
		if (err != nil) { break }

		// Get last message
		lastQuery := fmt.Sprintf("SELECT timestamp FROM log WHERE host='%s' ORDER BY timestamp DESC LIMIT 1", client.Host)
		lastRows, err := db_con.Query(lastQuery)
		defer lastRows.Close()
		lastRows.Next()
		var last int64
		err = lastRows.Scan(&last)
		t := time.Unix(last, 0)
		client.LastMessage = t.Format(time.RFC3339)
		if (err != nil) { break }

		// Get rate in last 24 hours
		dayAgo := time.Now().Unix() - 86400
		dayQuery := fmt.Sprintf("SELECT COUNT(*) FROM log WHERE host='%s' AND timestamp >= %d", client.Host, dayAgo)
		dayRows, err := db_con.Query(dayQuery)
		defer dayRows.Close()
		dayRows.Next()
		var inLastDay float64
		err = countRows.Scan(&inLastDay)
		client.Rate = inLastDay / 24
		if (err != nil) { break }

		results = append(results, client)
	}
	return results
}
