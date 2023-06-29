package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"dictionnaire/dictionary"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func main() {
	dict, err := dictionary.NewDictionary("./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dict.Close()

	router := mux.NewRouter()

	router.HandleFunc("/word/{word}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		word := params["word"]

		switch r.Method {
		case http.MethodGet:
			entry, err := dict.GetWord(word)
			if err != nil {
				if err == bolt.ErrBucketNotFound {
					http.NotFound(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(entry)

		case http.MethodDelete:
			err := dict.DeleteWord(word)
			if err != nil {
				if err == bolt.ErrBucketNotFound {
					http.NotFound(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet, http.MethodDelete)

	router.HandleFunc("/word", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var entry dictionary.Entry
			err := json.NewDecoder(r.Body).Decode(&entry)
			if err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
			}

			entry.CreatedAt = time.Now()

			err = dict.AddWord(entry.Word, entry.Definition, entry.CreatedAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(entry)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost)

	router.HandleFunc("/words", func(w http.ResponseWriter, r *http.Request) {
		entries, err := dict.GetAllWords()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(entries)
	}).Methods(http.MethodGet)

	log.Println("Listening on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", router))
}
