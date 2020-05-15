package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"path/filepath"
	"strings"
	"github.com/hpcloud/tail"
)

/*
	Struct defining a single log line on the client side
*/
type (
	ClientLog struct {
		Message 		string `json:"message"`
		Source			string `json:"source"`
		Host			string `json:"host"`
	}
)

/*
	Send log entry to remote server, as configured in avologo.conf
*/
func forwardToServer(entry ClientLog) {
	encoded, _ := json.Marshal(entry)
	req, _ := http.NewRequest("POST", "http://" + global_cfg.Client.Destination + "/log", bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
    client.Do(req)
}

/*
	Tail a file, sending new lines to the remote server
*/
func monitorFile(path string) {
	path, _ = filepath.Abs(path)
	seek := new(tail.SeekInfo)
	seek.Whence = 2
	seek.Offset = 0

	t, _ := tail.TailFile(path, tail.Config{Follow: true, Logger: tail.DiscardingLogger, Location: seek, Poll: true})
	for line := range t.Lines {
		var entry ClientLog

		if (strings.TrimSpace(line.Text) != ""){
			entry.Message = line.Text
			entry.Source = path

			// Set host to friendly name if requested
			if (global_cfg.Client.FriendlyName != "") {
				entry.Host = global_cfg.Client.FriendlyName
			}
	
			forwardToServer(entry)	
		}
	}
}