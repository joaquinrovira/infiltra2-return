package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joaquinrovira/infiltra2-returns/app/constants"
	"github.com/joaquinrovira/infiltra2-returns/app/model"
	"github.com/samber/do/v2"
)

type RoomsManager struct {
	rooms      map[string]*RoomState
	ctx        context.Context
	randomWord *RandomWordService
}

type RoomState struct {
	room            *model.Room
	listeners       map[string][]chan string
	listenersMutex  sync.RWMutex
	cancelCountdown func()
	countdownMutex  sync.Mutex
}

func newRoomState(roomName string) *RoomState {
	return &RoomState{
		room:            model.NewRoom(roomName),
		listeners:       make(map[string][]chan string),
		cancelCountdown: func() {},
	}
}

func NewRoomsManager(svc do.Injector) (*RoomsManager, error) {
	// NOTE: To retrieve dependencies thtough Dependency Injection:
	// dependency := do.MustInvoke[DependecyType](svc)
	lobby := &RoomsManager{
		rooms:      make(map[string]*RoomState),
		ctx:        do.MustInvoke[context.Context](svc),
		randomWord: do.MustInvoke[*RandomWordService](svc),
	}
	return lobby, nil
}

func (manager *RoomsManager) AddListener(roomName string, user string, c chan string) {
	var userListeners []chan string
	var ok bool

	state := manager.getOrCreateRoom(roomName)

	state.listenersMutex.Lock()
	defer state.listenersMutex.Unlock()

	if userListeners, ok = state.listeners[user]; !ok {
		userListeners = []chan string{c}
	} else {
		userListeners = append(userListeners, c)
	}
	state.listeners[user] = userListeners

	log.Printf("[%v] user '%v' is listening", roomName, user)
}

func (manager *RoomsManager) RemoveListener(roomName string, user string, c chan string) error {
	var userListeners []chan string

	room, ok := manager.rooms[roomName]
	if !ok {
		return fmt.Errorf("undefined room '%v'", roomName)
	}

	if userListeners, ok = room.listeners[user]; !ok {
		return nil
	}

	room.listenersMutex.Lock()
	defer room.listenersMutex.Unlock()

	for i, listener := range userListeners {
		if listener == c {
			newListeners := append(userListeners[:i], userListeners[i+1:]...)
			if len(newListeners) != 0 {
				room.listeners[user] = newListeners
			} else {
				delete(room.listeners, user)
			}
			log.Printf("[%v] user '%v' stopped (%v listeners left)", roomName, user, len(newListeners))
			break
		}
	}

	return nil
}

func (manager *RoomsManager) RoomUpdate(roomName string, update string) {

	if room, ok := manager.rooms[roomName]; ok {

		room.listenersMutex.RLock()
		defer room.listenersMutex.RUnlock()

		updateCtx, done := context.WithCancel(manager.ctx)

		// log.Printf("[%v] sending update", roomName)
		go func() {
			// for user, listeners := range room.listeners {
			for _, listeners := range room.listeners {
				// log.Printf("[%v] sending update to user '%v'", roomName, user)
				for _, listener := range listeners {
					listener <- update
				}
				// log.Printf("[%v] sent update to user '%v'", roomName, user)
			}
			done()
			// log.Printf("[%v] sending update done", roomName)
		}()

		// log.Printf("[%v] sending update waiting", roomName)
		<-updateCtx.Done()
		// log.Printf("[%v] sending update finished", roomName)
	}
}

func (manager *RoomsManager) GetOrCreateRoom(roomName string) *model.Room {
	return manager.getOrCreateRoom(roomName).room
}

func (manager *RoomsManager) getOrCreateRoom(roomName string) *RoomState {
	var state *RoomState
	var ok bool

	// Get or initialize a room
	if state, ok = manager.rooms[roomName]; !ok {
		state = newRoomState(roomName)
		manager.rooms[roomName] = state
		log.Printf("[%v] room created", roomName)
	}

	return state
}

func (manager *RoomsManager) closeRoom(roomName string) {
	var state *RoomState
	var exists bool
	defer log.Printf("[%v] room closed", roomName)
	if state, exists = manager.rooms[roomName]; !exists {
		return
	}

	state.cancelCountdown()
	delete(manager.rooms, roomName)
}

func (manager *RoomsManager) Join(roomName string, user string) {
	state := manager.getOrCreateRoom(roomName)

	if _, exists := state.room.User(user); exists {
		return
	}

	state.room.Join(user)
	log.Printf("[%v] user '%v' joins", roomName, user)
	manager.RoomUpdate(roomName, fmt.Sprintf("user '%v' joined", user))
	state.cancelCountdown()
}

func (manager *RoomsManager) LeaveIfNoListeners(roomName string, user string) {
	var state *RoomState
	var exists bool

	time.Sleep(300 * time.Millisecond)

	if state, exists = manager.rooms[roomName]; !exists {
		return
	}

	state.listenersMutex.RLock()
	defer state.listenersMutex.RUnlock()

	listeners, exists := state.listeners[user]
	if exists && len(listeners) > 0 {
		return
	}
	manager.Leave(roomName, user)
}

func (manager *RoomsManager) Leave(roomName string, user string) {
	state := manager.getOrCreateRoom(roomName)

	state.room.Leave(user)
	log.Printf("[%v] user '%v' leaves", roomName, user)
	manager.RoomUpdate(roomName, fmt.Sprintf("user '%v' left", user))

	if state.room.UserCount() == 0 {
		log.Printf("room '%v' is empty - cleaning up", roomName)
		manager.closeRoom(roomName)
		return
	}

	// Start countdown process
	if state.room.AllReady() {
		log.Printf("[%s] all users are ready - starting countdown", roomName)
		go manager.StartRound(roomName)
	}
}

func (manager *RoomsManager) ToggleReady(roomName string, user string) error {
	state := manager.getOrCreateRoom(roomName)

	state.cancelCountdown()

	if status, err := state.room.ToggleReady(user); err != nil {
		log.Printf("[ERROR] (%v,%v) %v", roomName, user, err)
		return err
	} else {
		readyStr := ReadyString(status)
		log.Printf("[%s] user '%s' is %s", roomName, user, readyStr)
		manager.RoomUpdate(roomName, fmt.Sprintf("user '%s' is %s", user, readyStr))
	}

	// Start countdown process
	if state.room.AllReady() {
		log.Printf("[%s] all users are ready - starting countdown", roomName)
		go manager.StartRound(roomName)
	}

	return nil
}

func (manager *RoomsManager) StartRound(roomName string) {
	state := manager.getOrCreateRoom(roomName)

	if !state.countdownMutex.TryLock() {
		state.cancelCountdown()
		state.countdownMutex.Lock()
	}
	defer state.countdownMutex.Unlock()

	ctx, cancel := context.WithTimeout(manager.ctx, constants.COUNTDOWN_TIME)
	state.cancelCountdown() // Cancel previous countdown just in case
	state.cancelCountdown = cancel
	state.room.CountdownActive = true
	manager.RoomUpdate(roomName, "countdown start")
	<-ctx.Done()
	state.room.CountdownActive = false

	if ctx.Err() == context.Canceled {
		log.Printf("[%s] countdown canceled", roomName)
		manager.RoomUpdate(roomName, "countdown canceled")
		return
	}

	state.room.SetWord(manager.randomWord.Next())
	selectedUser := state.room.SelectRandomUser()
	state.room.SetAllNotReady()
	log.Printf("[%s] countdown finished - user '%v' was selected", roomName, selectedUser)
	// TODO: Choose new word (random word getter can fetch next word async)
	manager.RoomUpdate(roomName, "new round")
}

func (manager *RoomsManager) Room(roomName string) (room *model.Room, exists bool) {
	state, exists := manager.rooms[roomName]
	return state.room, exists
}

func ReadyString(ready bool) string {
	if ready {
		return "ready"
	} else {
		return "not ready"
	}
}
