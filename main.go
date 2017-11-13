package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
)

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "127.0.0.1:27017"
	}
	session, err := mgo.Dial(mongodbURI)
	if err != nil {
		panic(err)
	}
	fmt.Println("Using DB: 127.0.0.1:27017")
	session.DB("heroku_9867bqd6").DropDatabase()
	dbm = session.DB("heroku_9867bqd6").C("messages")
	// dbc = session.DB("chat").C("channels")
	defer session.Close()

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleConn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7331"
	}
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
