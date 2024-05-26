package model

import (
	"fmt"
	"math/rand"
	"sync"
)

type Room struct {
	name            string
	mutex           sync.RWMutex
	users           map[string]UserRoomState
	word            *Word
	targetUser      *string
	CountdownActive bool
}

func NewRoom(name string) *Room {
	room := &Room{
		name:  name,
		mutex: sync.RWMutex{},
		users: make(map[string]UserRoomState),
	}

	return room
}

func (room *Room) Name() string {
	return room.name
}

func (room *Room) Join(user string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	room.users[user] = NewUserRoomState()
}

func (room *Room) Leave(user string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	delete(room.users, user)
}

func (room *Room) Ready(user string) error {
	var state UserRoomState
	var ok bool

	room.mutex.Lock()
	defer room.mutex.Unlock()

	if state, ok = room.users[user]; !ok {
		return fmt.Errorf("user '%v' is not a memeber of this room", user)
	}

	state.Ready = true
	room.users[user] = state
	return nil
}

func (room *Room) NotReady(user string) error {
	var state UserRoomState
	var ok bool

	room.mutex.Lock()
	defer room.mutex.Unlock()

	if state, ok = room.users[user]; !ok {
		return fmt.Errorf("user '%v' is not a memeber of this room", user)
	}

	state.Ready = false
	room.users[user] = state
	return nil
}

func (room *Room) ToggleReady(user string) (isReady bool, err error) {
	var state UserRoomState
	var ok bool

	room.mutex.Lock()
	defer room.mutex.Unlock()

	if state, ok = room.users[user]; !ok {
		return false, fmt.Errorf("user '%v' is not a member of this room", user)
	}

	state.Ready = !state.Ready
	room.users[user] = state
	return state.Ready, nil
}

func (room *Room) UserCount() int {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	return len(room.users)
}

func (room *Room) Users() map[string]UserRoomState {
	return room.users
}

func (room *Room) User(user_id string) (state UserRoomState, exists bool) {
	state, exists = room.users[user_id]
	return
}

func (room *Room) AllReady() bool {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	for _, user := range room.users {
		if !user.Ready {
			return false
		}
	}
	return true
}

func (room *Room) SelectRandomUser() string {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	targetIndex := rand.Intn(len(room.users))
	i := 0

	for user_id := range room.users {
		if i == targetIndex {
			room.targetUser = &user_id
			break
		}
		i = i + 1
	}

	return *room.targetUser
}

func (room *Room) SelectedUser() (user_id string, exists bool) {
	if room.targetUser == nil {
		return "", false
	} else {
		return *room.targetUser, true
	}
}

func (room *Room) SelectedWord() (word Word, exists bool) {
	if room.word == nil {
		return Word{}, false
	} else {
		return *room.word, true
	}
}

func (room *Room) SelectedWordOrEmpty() string {
	if room.word == nil {
		return ""
	} else {
		return room.word.Word
	}
}

func (room *Room) SetWord(word Word) {
	room.word = &word
}

func (room *Room) SetAllNotReady() {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	for user_id, user := range room.users {
		user.Ready = false
		room.users[user_id] = user
	}
}
