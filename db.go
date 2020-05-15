package main

import (
	"database/sql"
	"fmt"
)

/*
	Global database handle
*/
var (
	db_con *sql.DB
	db_err error
)

/*
	Struct defining a connection test json object and response
*/
type (
	TestConnection struct {
		Host 		string `json:"host"`
		Port 		int `json:"port"`
		User 		string `json:"user"`
		Password 	string `json:"password"`
		DBName 		string `json:"dbname"`
	}
	TestConnectionResponse struct {
		Valid 		bool `json:"valid"`
	}
)

/*
	Test a postgres connection and return a boolean
*/
func testDBConnection(conn *TestConnection) TestConnectionResponse {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
		conn.Host,
		conn.Port,
		conn.User,
		conn.Password,
		conn.DBName)

	// Response object
	var response TestConnectionResponse
	response.Valid = true

	context, err := sql.Open("postgres", connectionString)
	if (err != nil) {
		response.Valid = false
		return response
	}
		
	err = context.Ping()
	if (err != nil) {
		response.Valid = false
		return response
	}

	return response
}

/* 
	Returns a database handle
*/
func getDBHandle() (*sql.DB, error) {
	
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		global_cfg.Database.Host,
		global_cfg.Database.Port,
		global_cfg.Database.User,
		global_cfg.Database.Password,
		global_cfg.Database.DBName)

	context, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	
	err = context.Ping()
	if err != nil {
		panic(err)
	}

	return context, err
}