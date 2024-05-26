package endpoints

import (
	"fmt"
	"net/http"

	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/services"
	"github.com/samber/do/v2"
)

func RoomEvents(svc *do.RootScope) http.HandlerFunc {
	controller := do.MustInvoke[*services.RoomsManager](svc)

	return func(w http.ResponseWriter, r *http.Request) {
		room_id := r.PathValue(constants.PATH_PARAM_ROOM_ID)
		user_id := RequestingUserId(w, r)

		defer controller.LeaveIfNoListeners(room_id, user_id) // NOTE: Very important LeaveIfNoListeners() *after* RemoveListener()

		// log.Printf("received connection from user '%v'", user_id)
		// defer log.Printf("closed connection from user '%v'", user_id)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		events := make(chan string, 1)
		defer close(events)

		controller.AddListener(room_id, user_id, events)
		defer controller.RemoveListener(room_id, user_id, events)
		controller.Join(room_id, user_id)

		flusher := w.(http.Flusher)

		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		for {
			select {
			case event := <-events:
				fmt.Fprintf(w, "event:RoomUpdate\ndata: %s\n\n", event)
				flusher.Flush()

			case <-r.Context().Done():
				return
			}
		}
	}
}
