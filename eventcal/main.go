package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/tenntenn/connpass"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	cli := connpass.NewClient()
	params, err := connpass.SearchParam(connpass.Keyword("golang"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	result, err := cli.Search(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	for _, e := range result.Events {
		id := fmt.Sprintf("connpass-%d", e.ID)
		event := cal.AddEvent(id)
		event.SetCreatedTime(e.UpdatedAt)
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(e.UpdatedAt)
		event.SetStartAt(e.StartedAt)
		event.SetEndAt(e.EndedAt)
		event.SetSummary(e.Title)
		event.SetLocation(e.Place)
		event.SetDescription(e.Description)
		event.SetURL(e.URL)
	}

	s := cal.Serialize()
	fmt.Fprint(w, s)
}
