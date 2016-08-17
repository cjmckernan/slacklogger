package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type WebResponse struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

//Handles a request at the root and processes if the token is correct
func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//Checks token is same as one on launch
		if r.PostFormValue("token") == tokenID {
			//Check that values its not a bot trying to grep itself >_>
			if strings.Contains(r.PostFormValue("text"), "grepbot") && r.PostFormValue("user_name") != "slackbot" {
				grepMessage := SplitText(r.PostFormValue("text"))
				cmd := rungrep(grepMessage)
				var resp WebResponse
				resp.Username = "grepbot"
				resp.Text = cmd
				b, err := json.Marshal(resp)
				if err != nil {
					log.Fatal(err)

				}
				w.Write(b)
			}
			fullMessage := FormatMessage(r)
			Writelog(fullMessage)
		} else {
			var rejectresp WebResponse
			rejectresp.Username = "grepbot"
			rejectresp.Text = "You do not have permission to access this bot, eat shit"
			b, err := json.Marshal(rejectresp)
			if err != nil {
				log.Fatal(err)
			}
			w.Write(b)

		}
	})

}

//Starts server on the port that was passed through on start
func Serve(port int) {
	log.Printf("Starting server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Serve error", err)
	}
}
