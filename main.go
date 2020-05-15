package main

import (
	"github.com/labstack/echo"
	"github.com/akamensky/argparse"
	"html/template"
	"fmt"
	"strconv"
	"os"
)

func main() {

	// Parse command line arguments
	parser := argparse.NewParser("avologo", "A fast, lightweight, and modular log aggregation tool")
	mode := parser.String("m", "mode", &argparse.Options{Required: true, Help: "server/client"})
	err := parser.Parse(os.Args)
	if (err != nil) {
		fmt.Print(parser.Usage(err))
		return
	}

	// Parse configuration file
	global_cfg = parseConfig("/etc/avologo.conf")

	// Determine appropriate mode to initialize in
	if (*mode == "server") {
		initializeServer();

	} else if (*mode == "client") {
		initializeClient();
	} else {
		fmt.Println("Invalid mode specified")
	}
}

/*
	Server mode
*/
func initializeServer() {
	e := echo.New()
	e.HideBanner = true
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "assets")

	// Initialize GET handlers
	for endpoint, handler := range getHandlers {
		e.GET(endpoint, handler)
	}
	// Initialize POST handlers
	for endpoint, handler := range postHandlers {
		e.POST(endpoint, handler)
	}

	db_con, db_err = getDBHandle()
	
	// Start listening
	e.Logger.Fatal(e.Start(global_cfg.Server.Host + ":" + strconv.Itoa(global_cfg.Server.Port)))
	db_con.Close()
}

/*
	Client mode
*/
func initializeClient() {
	
	// Spawn threads for watching files
	for _, path := range global_cfg.Client.Watch {
		go monitorFile(path)
	}

	// Main thread sleep
	for { }
}