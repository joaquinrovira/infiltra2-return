package endpoints

import (
	"fmt"
	"net/http"

	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/model"
	"github.com/joaquinrovira/infiltra2-returns/app/routes"
	"github.com/joaquinrovira/infiltra2-returns/app/services"
	"github.com/samber/do/v2"
)

func Ready(svc *do.RootScope) http.HandlerFunc {
	controller := do.MustInvoke[*services.RoomsManager](svc)

	return func(w http.ResponseWriter, r *http.Request) {
		room_id := r.PathValue(constants.PATH_PARAM_ROOM_ID)
		user_id := RequestingUserId(w, r)

		var exists bool
		var room *model.Room

		if room, exists = controller.Room(room_id); !exists {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf("unknown room '%v'", room_id)))
			return
		}

		if _, exists = room.User(user_id); !exists {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(fmt.Sprintf("you are not a member of room '%v'", user_id)))
			return
		}

		if err := controller.ToggleReady(room_id, user_id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		http.Redirect(w, r, routes.RoomSpecific(room_id), http.StatusFound)
	}
}
