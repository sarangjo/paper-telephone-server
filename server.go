package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	rooms map[uuid.UUID]*Room
    users map[string]*Room
}

// Create a room on this server.
func (this *Server) createRoom(w http.ResponseWriter, remoteAddr string) uuid.UUID {
    r := newRoom()
    this.rooms[r.uuid] = r
    return r.uuid
}

// Join a room on this server.
func (this *Server) joinRoom(w http.ResponseWriter, roomUuid uuid.UUID, remoteAddr string) bool {
    if _, ok := this.rooms[roomUuid]; !ok {
        fmt.Fprintf(w, "Room does not exist")
        return false
    }
    if err := this.rooms[roomUuid].addMember(remoteAddr); err != nil {
        fmt.Fprintf(w, "%s", err)
        return false
    }
    this.users[remoteAddr] = this.rooms[roomUuid]
    return true
}

// Handle room requests
func (this *Server) roomHandler(w http.ResponseWriter, r *http.Request) {
    req := strings.Split(r.URL.Path, "/")
    fmt.Fprintf(w, "Request: %s", r.RemoteAddr)

    switch req[2] {
    case "create":
        // Create a room
        u := this.createRoom(w, r.RemoteAddr)
        if this.joinRoom(w, u, r.RemoteAddr) {
            fmt.Fprintf(w, "%s", u)
        }
        break
    case "join":
        // Join a room
        roomUuid, err := uuid.FromString(r.Header[http.CanonicalHeaderKey("room")][0])
        if err != nil {
            fmt.Fprintf(w, "Non-UUID room provided")
            break
        }
        this.joinRoom(w, roomUuid, r.RemoteAddr)
        break
    case "start":
        // Start a game
        break
    default:
        fmt.Fprintf(w, "unknown request")
        break
    }

    fmt.Printf("%v", this.rooms)
}

func (this *Server) start() {
    http.HandleFunc("/room/", this.roomHandler)
    fmt.Println("Serving at 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    s := Server{}
    s.rooms = make(map[uuid.UUID]*Room)
    s.start()
}
