package main

import (
  "net/http"
  "time"
  "encoding/json"
  "regexp"
  "strconv"
  "crypto/md5"
  "fmt"
)

type User struct {
  Name string `json:"name"`
}

func newUser(name string) User {
  return User{name}
}

type Message struct {
  User *User `json:"user"`
  Content string `json:"content"`
}

func newMessage(user *User, content string) Message {
  return Message{user, content}
}

type Room struct {
  Id int64 `json:"id"`
  Messages []*Message `json:"messages"`
  Creator *User `json:"creator"`
  Name string `json:"name"`
}

func newRoom(id int64, creator *User, name string) Room {
  return Room{id, []*Message{}, creator, name}
}

type chatServer struct {
  rooms map[int64]*Room
  users map[string]*User
}

func newChatServer() chatServer {
  users := make(map[string]*User)

  for _, name := range []string{"bruno mascarenhas", "Iago", "igorqs", "richardtrle", "sowter"} {
    md5 := md5.Sum([]byte(name))
    user := newUser(name)
    users[fmt.Sprintf("%x", md5)] = &user
  }

  return chatServer{map[int64]*Room{}, users}
}

func (t *chatServer) handleRequests() {
  http.HandleFunc("/chat/rooms", t.collectionHandler)
  http.HandleFunc("/chat/rooms/", t.itemHandler)
}

func (t *chatServer) collectionHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    t.listRooms(w, r)
    return;
  case "POST":
    t.createRoom(w, r)
    return;
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
    return;
  }
}

func (t *chatServer) listRooms(w http.ResponseWriter, r *http.Request) {
  rooms := make([]*Room, len(t.rooms))

  i := 0
  for _, room := range t.rooms {
    rooms[i] = room
    i++
  }

  jsonBytes, err := json.Marshal(rooms)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  }

  w.Header().Add("content-type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(jsonBytes)
}

func (t *chatServer) createRoom(w http.ResponseWriter, r *http.Request) {
  user := t.getUserFromMemory(w, r)
  if user == nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  var room Room
  err := json.NewDecoder(r.Body).Decode(&room)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  id := t.uniqid()
  room = newRoom(id, user, room.Name)
  t.rooms[id] = &room

  t.showRoom(w, r, &room)
}

func (t *chatServer) itemHandler(w http.ResponseWriter, r *http.Request) {
  room := t.getRoomFromMemory(w, r)
  if room == nil {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  switch r.Method {
  case "GET":
    t.readRoom(w, r, room)
    return;
  case "POST":
    t.createMessage(w, r, room)
    return;
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
    return;
  }
}

func (t *chatServer) readRoom(w http.ResponseWriter, r *http.Request, room *Room) {
  regex := regexp.MustCompile(`/chat/rooms/\d*$`)
  if !regex.MatchString(r.URL.String()) {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  t.showRoom(w, r, room)
}

func (t *chatServer) showRoom(w http.ResponseWriter, r *http.Request, room *Room) {
  jsonBytes, err := json.Marshal(room)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  }

  w.Header().Add("content-type", "application/json")
  w.WriteHeader(http.StatusCreated)
  w.Write(jsonBytes)
}

func (t *chatServer) createMessage(w http.ResponseWriter, r *http.Request, room *Room) {
  regex := regexp.MustCompile(`/chat/rooms/\d*/messages$`)
  if !regex.MatchString(r.URL.String()) {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  user := t.getUserFromMemory(w, r)
  if user == nil {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  var message Message
  err := json.NewDecoder(r.Body).Decode(&message)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  message = newMessage(user, message.Content)
  room.Messages = append(room.Messages, &message)

  t.showRoom(w, r, room)
}

func (t *chatServer) getRoomFromMemory(w http.ResponseWriter, r *http.Request) *Room {
  regex := regexp.MustCompile(`/chat/rooms/(\d*)`)
  matches := regex.FindStringSubmatch(r.URL.String())
  if len(matches) != 2 {
    return nil
  }

  id, err := strconv.ParseInt(matches[1], 10, 64)
  if err != nil {
    return nil
  }

  room, ok := t.rooms[id]
  if !ok {
    return nil
  }

  return room
}

func (t *chatServer) getUserFromMemory(w http.ResponseWriter, r *http.Request) *User {
  token := r.Header.Get("Authorization")

  user, ok := t.users[token]
  if !ok {
    return nil
  }

  return user
}

func (t *chatServer) uniqid() int64 {
  now := time.Now()
  return now.UnixNano()
}

