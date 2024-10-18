package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"

    "github.com/go-chi/chi/v5"
)

// Definição da estrutura User
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var (
    users = make(map[string]User) // Mapa para armazenar usuários
    mu    sync.Mutex               // Mutex para acesso seguro
)

func main() {
    r := chi.NewRouter()

    // Rotas CRUD
    r.Post("/users", createUser)
    r.Get("/users", getUsers)
    r.Get("/users/{id}", getUser)
    r.Put("/users/{id}", updateUser)
    r.Delete("/users/{id}", deleteUser)

    // Iniciar o servidor
    http.ListenAndServe(":8080", r)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    mu.Lock()
    users[user.ID] = user
    mu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    var userList []User
    for _, user := range users {
        userList = append(userList, user)
    }

    json.NewEncoder(w).Encode(userList)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    mu.Lock()
    user, exists := users[id]
    mu.Unlock()

    if !exists {
        http.NotFound(w, r)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    mu.Lock()
    user, exists := users[id]
    mu.Unlock()

    if !exists {
        http.NotFound(w, r)
        return
    }

    var updatedUser User
    if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    updatedUser.ID = id
    mu.Lock()
    users[id] = updatedUser
    mu.Unlock()

    json.NewEncoder(w).Encode(updatedUser)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    mu.Lock()
    _, exists := users[id]
    if exists {
        delete(users, id)
    }
    mu.Unlock()

    if !exists {
        http.NotFound(w, r)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
