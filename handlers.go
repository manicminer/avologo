package main

import (
	"net/http"
	"strconv"
	"time"
	"strings"
	"io/ioutil"
	"github.com/labstack/echo"
)

/*
	Map of GET handlers
*/
var getHandlers = map[string]echo.HandlerFunc {
	"/": rootHandlerG,
	"/clients" : clientsHandlerG,
	"/view/:id": viewHandlerG,
	"/log/:id": logHandlerG,
	"/query": queryHandlerG,
	"/queryClients" : queryClientsHandlerG,
}

/*
	Map of POST handlers
*/
var postHandlers = map[string]echo.HandlerFunc {
	"/log": logHandlerP,
	"/logRaw": lawRawHandlerP,
}

/*
	Dashboard page
*/
func rootHandlerG(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

/*
	View single entry page
*/
func viewHandlerG(c echo.Context) error {
	return c.Render(http.StatusOK, "view.html", map[string]interface{}{})
}

/*
	View clients page
*/
func clientsHandlerG(c echo.Context) error {
	return c.Render(http.StatusOK, "clients.html", map[string]interface{}{})
}

/*
	POST /log
*/
func logHandlerP(c echo.Context) (err error) {
	log := new(LogEntry)
	if err = c.Bind(log); err != nil {
		return
	}

	// Set timestamp, if not included
	if (log.Timestamp == 0) {
		log.Timestamp = time.Now().Unix()
	}

	// Set source to raw, if not included
	if (log.Source == "") {
		log.Source = "raw"
	}

	// Set host to request ip, if not included
	if (log.Host == "") {
		log.Host = strings.Split(c.Request().RemoteAddr, ":")[0]
	}

	log.Severity = calculateSeverity(log.Message)

	writeLogEntry(log)
	return
}

/*
	POST /logRaw
*/
func lawRawHandlerP(c echo.Context) (err error) {
	log := new(LogEntry)
	
	// Set timestamp
	log.Timestamp = time.Now().Unix()

	// Set message body, source, and host
	var body []byte
	if c.Request().Body != nil {
		body, _ = ioutil.ReadAll(c.Request().Body)
	}
	log.Message = string(body)
	log.Source = "raw"
	log.Host = strings.Split(c.Request().RemoteAddr, ":")[0]

	// Write to database
	writeLogEntry(log)
	return
}

/*
	GET /log/:id
*/
func logHandlerG(c echo.Context) error {
	search_id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	// Retreive entry from DB
	entry := getLogEntry(search_id)

	// Convert timestamp
	t := time.Unix(entry.Timestamp, 0)
	entry.TimeFormatted = t.Format(time.RFC3339)
	
	return c.JSON(http.StatusOK, entry)
}

/*
	GET /query
	Query parameters:
	* search = string
	* limit = int
	* host = string
	* source = string
	* severity = int
	* upper = int
	* lower = int
	* order = int (0 = desc, 1 = asc)
*/
func queryHandlerG(c echo.Context) (err error) {

	// Build query struct from params
	var query LogQuery
	query.SearchFilter = c.QueryParam("search")
	query.HostFilter = c.QueryParam("host")
	query.SourceFilter = c.QueryParam("source")
	query.SeverityFilter, _ = strconv.ParseInt(c.QueryParam("severity"), 10, 64)
	query.Limit, _ = strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	query.UpperBound, _ = strconv.ParseInt(c.QueryParam("upper"), 10, 64)
	query.LowerBound, _ = strconv.ParseInt(c.QueryParam("lower"), 10, 64)
	query.Order, _ = strconv.ParseInt(c.QueryParam("order"), 10, 64)
	query.Offset, _ = strconv.ParseInt(c.QueryParam("offset"), 10, 64)

	return c.JSON(http.StatusOK, performQuery(query))
}

/*
	GET /queryClients
*/
func queryClientsHandlerG(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, performClientQuery()) 
}