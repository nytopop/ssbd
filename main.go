package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"golang.org/x/crypto/ssh"

	"github.com/nytopop/ssbd/config"
	"github.com/nytopop/ssbd/data"
	"github.com/nytopop/ssbd/logs"
	"github.com/nytopop/ssbd/models"
	"github.com/nytopop/ssbd/web"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// CMD flags
	confFile := flag.String("conf", "etc/test.conf", "path to configuration file.")
	flag.Parse()

	/*
	 We don't use file loggers for initialization, switch to
	 them after server is running. This way, errors are logged
	 to syslog even if file loggers do not initialize correctly.
	*/

	// Load config
	err := config.LoadConfig(*confFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Start logging
	err = logs.InitLoggers()
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize database
	db, err := models.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Start job scheduler
	s := data.NewScheduler(db)
	go s.Run()
	defer s.Close()

	/****************************/
	// testingBacks()
	/****************************/

	// Start HTTP handlers : THIS WILL BLOCK UNTIL EXIT
	err = web.StartServer(db)
	if err != nil {
		log.Fatalln(err)
	}
}

// temporary while infra is being written
func testingBacks() {
	var user string
	var pw string
	fmt.Println("[?user] [?pw]")
	fmt.Scanln(&user, &pw)
	fmt.Println("got it!")

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pw),
		},
	}

	client, err := ssh.Dial("tcp", "172.18.9.241:22", cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	err = data.GetFullTar(client,
		"/home/eric/doc/notes",
		"run/bak/f0")
	if err != nil {
		log.Fatalln(err)
	}

	err = data.GetDiffTar(client,
		"run/bak/f0",
		"/home/eric/doc/notes",
		"run/bak/d0")
	if err != nil {
		log.Fatalln(err)
	}
}
