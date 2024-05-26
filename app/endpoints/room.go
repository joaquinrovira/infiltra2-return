package endpoints

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/joaquinrovira/infiltra2-returns/app/components"
	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/routes"
	"github.com/joaquinrovira/infiltra2-returns/app/services"
	"github.com/joaquinrovira/infiltra2-returns/app/util"
	"github.com/samber/do/v2"
)

const USER_ID_COOKIE = "user-id"

func Room(svc *do.RootScope) http.HandlerFunc {
	controller := do.MustInvoke[*services.RoomsManager](svc)

	return func(w http.ResponseWriter, r *http.Request) {
		room_id_raw := r.PathValue(constants.PATH_PARAM_ROOM_ID)
		room_id := util.ToSlug(room_id_raw)
		
		if(room_id_raw != room_id) { // ensure room_id are slugs or redirect to its slug
			http.Redirect(w, r, routes.RoomSpecific(room_id), http.StatusFound)
		}
		
		user_id := RequestingUserId(w, r)
		room := controller.GetOrCreateRoom(room_id)

		if IsSSETriggeredRequest(room_id, r) {
			// Request sent by HTMX SSE trigger
			templ.Handler(components.Room(user_id, room)).ServeHTTP(w, r)
		} else {
			// Normal user request
			templ.Handler(components.RoomFull(user_id, room)).ServeHTTP(w, r)
		}
	}
}

func RequestingUserId(w http.ResponseWriter, r *http.Request) (user_id string) {
	if user_cookie, err := r.Cookie(USER_ID_COOKIE); err == http.ErrNoCookie {
		user_id = uuid.NewString()
		http.SetCookie(w, &http.Cookie{
			Name:  USER_ID_COOKIE,
			Value: user_id,
		})
	} else {
		user_id = user_cookie.Value
	}

	return user_id
}

func IsSSETriggeredRequest(room_id string, r *http.Request) bool {

	values, exists := r.Header["Hx-Request"]

	if !exists {
		return false
	}
	if values[0] != "true" {
		return false
	}

	values, exists = r.Header["Hx-Current-Url"]
	if !exists {
		return false
	}

	current := values[0]
	expectedSuffix := routes.RoomSpecific(room_id)
	return strings.HasSuffix(current, expectedSuffix)
}
