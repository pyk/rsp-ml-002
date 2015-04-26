package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/robfig/cron"
)

var (
	DATABASE_URL = os.Getenv("DATABASE_URL")
)

func waitNotifications(l *pq.Listener) {
	for {
		log.Println("Wait a notification...")
		select {
		case n := <-l.Notify:
			log.Printf("Notification received, schedule a new cron jobs with url %q\n", n.Extra)
			// TODO: run each function on go routine
			c := cron.New()
			c.AddFunc("@every 5s", func() {
				fmt.Printf("cron: fetch %v\n", n.Extra)
			})
			c.Start()
			return
		case <-time.After(90 * time.Second):
			go func() {
				l.Ping()
			}()
			// Check if there's more work available, just in case it takes
			// a while for the Listener to notice connection loss and
			// reconnect.
			fmt.Println("received no work for 90 seconds, checking for new work")
			return
		}
	}
}

func main() {
	_, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		os.Exit(0)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(DATABASE_URL, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("sources")
	if err != nil {
		log.Fatal(err)
	}

	for {
		waitNotifications(listener)
	}
}
