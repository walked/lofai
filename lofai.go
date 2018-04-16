package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	tail "github.com/hpcloud/tail"
)

//DefaultOffset  -  File offset for reading from end (in bytes)
var DefaultOffset *int64

//LogFile  -  Location of logfile to monitor
var LogFile *string

func main() {
	DefaultOffset = flag.Int64("offset", 2048, "Offset for how far back to read a log file")
	LogFile = flag.String("logfile", os.Getenv("LOFAI_LOG"), "Absolute path to monitored logfile; will search for EnvVar LOFAI_LOG")
	Port := flag.String("port", "8000", "Default port to listen on")
	flag.Parse()
	if *LogFile == "" {
		fmt.Println("Please Specify Logfile")
		return
	}
	runWeb(*Port)
}

// Kickoff web service on port p
func runWeb(p string) {
	router := mux.NewRouter()
	router.HandleFunc("/streamer", streamer).Methods("GET")
	router.HandleFunc("/search/{searchTerm}", search).Methods("GET")
	router.HandleFunc("/get/{numchar}", getData).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+p, router)) //opens on specified port
}

//Streams actively the most recent lines; starting DefaultOffset bytes inward
func streamer(w http.ResponseWriter, r *http.Request) {
	// headers necessary for keepalive
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	seek := tail.SeekInfo{Offset: -*DefaultOffset, Whence: 2} // Will seek to the last DefaultOffset bytes of file for initial read
	t, err := tail.TailFile(*LogFile, tail.Config{Follow: true, ReOpen: true, Location: &seek})
	if err != nil {
		panic(err)
	}
	afterFirst := false //flag to skip first line due to truncation formatting
	for line := range t.Lines {

		if afterFirst {
			w.Write([]byte(line.Text + "\n"))
			w.(http.Flusher).Flush()
		}
		afterFirst = true
	}

}

// request the number of bytes to read from the bottom of the given log file
func getData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	i, err := strconv.ParseInt(vars["numchar"], 10, 64)
	seek := tail.SeekInfo{Offset: -i * 64, Whence: 2} // Will seek to the last Offset characters of file for initial read

	t, _ := tail.TailFile(*LogFile, tail.Config{Follow: false, Location: &seek})

	if err != nil {
		panic(err)
	}

	afterFirst := false // The first line is truncated to var["numchar] bytes; this skips the first line using this flag
	for line := range t.Lines {
		if afterFirst {
			w.Write([]byte(line.Text + "\n"))
		}
		afterFirst = true
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	inFile, _ := os.Open(*LogFile)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)

	scanner.Split(bufio.ScanLines)

	re := regexp.MustCompile("^.*" + regexp.QuoteMeta(vars["searchTerm"]) + ".*$")

	for scanner.Scan() {
		if re.Match([]byte(scanner.Text())) {
			w.Write([]byte(scanner.Text() + "\n"))
		}
	}
}
