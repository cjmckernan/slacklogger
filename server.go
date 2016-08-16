package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type WebResponse struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.PostFormValue("token") == tokenID {
			fmt.Println(r.PostFormValue("user_name"))
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

func Serve(port int) {
	log.Printf("Starting server on %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Serve error", err)
	}
}
