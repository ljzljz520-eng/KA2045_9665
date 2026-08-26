package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"workpay/api"
	"workpay/store"
	"workpay/workflow"
)

func main() {
	path := flag.String("db", "payroll.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	db, e := store.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	svc := workflow.New(db)
	log.Printf("listening on %s", *addr)
	if e = http.ListenAndServe(*addr, api.Routes(api.New(svc))); e != nil && e != http.ErrServerClosed {
		os.Exit(1)
	}
}
