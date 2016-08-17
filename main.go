package main

import (
	"flag"
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
	flag.IntVar(&httpPort, "port", 9001, "HTTP Port to use")
	flag.StringVar(&logFileName, "logname", "chaterino.log", "Name of the log file")
	flag.StringVar(&botUsername, "botName", "grepbot", "Name the bot for grep")
	flag.StringVar(&tokenID, "token", "", "Token for the outgoing webhook")
	flag.Parse()
	Serve(httpPort)
}

func SplitText(message string) (grep string) {
	split := strings.Split(message, " ")
	_, grep = split[0], split[1]
	return
}
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
func Writelog(message string) {
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(message); err != nil {
		panic(err)
	}
}
func rungrep(phrase string) (output string) { //(answer string) {
	cmd := exec.Command("grep", "-i", phrase, logFileName)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		println(err.Error())
	}
	output = string(resp)
	return
}
