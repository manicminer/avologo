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