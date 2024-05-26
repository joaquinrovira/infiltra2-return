package model


type UserRoomState struct {
	Ready bool
}

func NewUserRoomState() UserRoomState {
	return UserRoomState{
		Ready: false,
	}
}