package main

import (
    "errors"
	"github.com/satori/go.uuid"
)

type Room struct {
    members map[string]bool
    uuid    uuid.UUID
}

func newRoom() *Room {
    r := &Room{}
    r.members = make(map[string]bool)
    r.uuid = uuid.Must(uuid.NewV4())
    return r
}

func (this *Room) addMember(remoteAddr string) error {
    if _, ok := this.members[remoteAddr]; ok {
        return errors.New("User already member of room")
    }
    this.members[remoteAddr] = true
    return nil
}
