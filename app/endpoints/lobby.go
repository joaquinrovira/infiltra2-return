package endpoints

import (
	"log"
	"net/http"

	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/routes"
	"github.com/joaquinrovira/infiltra2-returns/app/util"
)

func Lobby() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			// TODO: Error handling
			log.Printf("[ERR] %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	
		lobby := r.FormValue(constants.FORM_LOBBY_NAME)
		if lobby == "" {
			http.Redirect(w, r, routes.Home(), http.StatusSeeOther)
			return
		}
	
		lobbySlug := util.ToSlug(lobby)
		http.Redirect(w, r, routes.RoomSpecific(lobbySlug), http.StatusFound)
	}
}
