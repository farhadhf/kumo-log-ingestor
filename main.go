package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/farhadhf/kumo-log-injestor/store"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("PostgreSQL DATABASE_URL env variable is required.")
	}
	db, err := store.Connect(dbURL)
	if err != nil {
		panic(err)
	}
	if err := db.InitDatabase(); err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	http.HandleFunc("POST /", eventHandler(db, logger))

	listenAddr := "127.0.0.1:3000"
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		listenAddr = addr
	}
	logger.Info(fmt.Sprintf("Listening on %s", listenAddr))

	// @TODO You might want to launch a go routine to periodically delete old records to keep the table size manageable.

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		panic(err)
	}
}

func eventHandler(db *store.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We don't want to block the request (that is, KumoMTA) while decoding the body and inserting it into the database.
		// Also, we don't care about errors. Returning non-2xx response causes KumoMTA to reschedule the event for delivery,
		// which can cause the KumoMTA queues to quickly grow to unmanageable sizes in case there's a problem here.
		// So, we accept that there will be a small number missing events from the database due to errors (which do get logged),
		// and just respond with 200 OK.
		io.WriteString(w, "OK")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("error while reading the request body", "error", err.Error())
			return
		}

		// And then go handle the decoding and database insert in a separate goroutine.
		go func(body []byte) {
			event := store.Event{}
			if err := json.Unmarshal(body, &event); err != nil {
				logger.Error("error while unmarshalling the KumoMTA event", "error", err.Error())
				return
			}

			if err := db.InsertEvent(&event); err != nil {
				logger.Error("error while inserting the event into the database", "error", err.Error())
			}
		}(body)
	}
}
