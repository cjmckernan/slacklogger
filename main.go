package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	httpPort    int
	botUsername string
	logFileName string
	tokenID     string
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: ./slackbot -port=9001 -token=xxxxx -logname=abc.log\n")
		flag.PrintDefaults()

	}
	flag.IntVar(&httpPort, "port", 9001, "HTTP Port to use")
	flag.StringVar(&logFileName, "logname", "chaterino.log", "Name of the log file")
	flag.StringVar(&botUsername, "botName", "grepbot", "Name the bot for grep")
	flag.StringVar(&tokenID, "token", "", "Token for the outgoing webhook")
	flag.Parse()
	Serve(httpPort)
}

//Function to split text in grepbot command
func SplitText(message string) (grep string) {
	split := strings.Split(message, " ")
	_, grep = split[0], split[1]
	return
}

//Format message to be written to the log file
func FormatMessage(r *http.Request) (message string) {
	incomingVals := r.PostFormValue("text")
	channelVal := r.PostFormValue("channel_name")
	username := r.PostFormValue("user_name")
	timestamp := r.PostFormValue("timestamp")
	x, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		panic(err)
	}
	var y int64 = int64(x)
	realTime := time.Unix(y, 0)
	message = realTime.String() + "   " + channelVal + "    " + username + "    " + incomingVals + "\n"
	return
}

//Function to write messages to the log file.
func Writelog(message string) {
	var _, err = os.Stat(logFileName)
	//Check if file exists if not create one
	if os.IsNotExist(err) {
		//Create file and add the message
		var file, err = os.Create(logFileName)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = file.WriteString(message); err != nil {
			panic(err)
		}

		defer file.Close()
		//Else just add the message to the log file
	} else {
		f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err = f.WriteString(message); err != nil {
			panic(err)
		}
	}
}

//Run grep command from input.
func rungrep(phrase string) (output string) {
	cmd := exec.Command("grep", "-i", phrase, logFileName)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		println(err.Error())
	}
	output = string(resp)
	return
}
