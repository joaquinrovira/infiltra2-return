package routes

import (
	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/util"
)

func Home() string {
	return "/"
}

func CatchAll() string {
	return "/*"
}

func Lobby() string {
	return "/lobby"
}

func RoomTemplate() string {
	return "/room/" + util.WrapPathKey(constants.PATH_PARAM_ROOM_ID)
}

func RoomSpecific(id string) string {
	return "/room/" + id
}

func RoomSSETemplate() string {
	return "/room/" + util.WrapPathKey(constants.PATH_PARAM_ROOM_ID) + "/events"
}

func RoomSSESpecific(id string) string {
	return "/room/" + id + "/events"
}


func ReadyTemplate() string {
	return "/room/" + util.WrapPathKey(constants.PATH_PARAM_ROOM_ID) + "/ready"
}

func ReadySpecific(id string) string {
	return "/room/" + id + "/ready"
}