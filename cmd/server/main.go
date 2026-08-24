package main

import (
	"fmt"
	"gymrecommend/internal/config"
	"gymrecommend/internal/httpapi"
	"gymrecommend/internal/importcsv"
	"gymrecommend/internal/storage"
	"gymrecommend/internal/workflow"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := config.Load()
	store, e := storage.Open(cfg.DBPath)
	if e != nil {
		log.Fatal(e)
	}
	defer store.Close()
	svc := workflow.NewService(store)
	if f, e := os.Open(cfg.DataDir + "/members.csv"); e == nil {
		members, _ := importcsv.ReadMembers(f)
		f.Close()
		for _, m := range members {
			_ = svc.RegisterMember(m)
		}
	}
	if f, e := os.Open(cfg.DataDir + "/classes.csv"); e == nil {
		classes, _ := importcsv.ReadClasses(f)
		f.Close()
		for _, c := range classes {
			_ = svc.RegisterClass(c)
		}
	}
	fmt.Printf("gym recommendation server listening on %s\n", cfg.Address())
	log.Fatal(http.ListenAndServe(cfg.Address(), httpapi.New(svc).Handler()))
}
