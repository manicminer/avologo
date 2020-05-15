package main

import (
	"github.com/labstack/echo"
	"github.com/akamensky/argparse"
	"html/template"
	"fmt"
	"strconv"
	"os"
	_ "github.com/lib/pq"
)

func main() {

	// Parse command line arguments
	parser := argparse.NewParser("avologo", "A fast, lightweight, and modular log aggregation tool")
	mode := parser.String("m", "mode", &argparse.Options{Required: true, Help: "server/client"})
	config_path := parser.String("c", "config", &argparse.Options{Required: false, Help: "path to alternate config file"})

	err := parser.Parse(os.Args)
	if (err != nil) {
		fmt.Print(parser.Usage(err))
		return
	}

	// Config path
	global_cfg_path = "/etc/avologo.conf"
	if (*config_path != "") {
		if (!fileExists(*config_path)) {
			fmt.Println("Error: specified config file not found")
			return
		}
		global_cfg_path = *config_path
	}

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
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "./assets")

	// Initialize GET handlers
	for endpoint, handler := range getHandlers {
		e.GET(endpoint, handler)
	}
	// Initialize POST handlers
	for endpoint, handler := range postHandlers {
		e.POST(endpoint, handler)
	}

	// Parse configuration file and get db handle
	if (fileExists(global_cfg_path)) {
		global_cfg = parseConfig(global_cfg_path)
		db_con, db_err = getDBHandle()

		// Start listening
		e.Logger.Fatal(e.Start(global_cfg.Server.Host + ":" + strconv.Itoa(global_cfg.Server.Port)))
		db_con.Close()
	} else {
		// Listen on default settings
		e.Logger.Fatal(e.Start("0.0.0.0:4747"))
	}
	
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