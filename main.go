package main

import (
	album_db "api_mux/db"
	"api_mux/handlers"
	"api_mux/types"
	"net/http"
	"sync"
)

func main() {

	//fake database
	var albums = []types.Album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	//creating the map
	store := &album_db.Db{
		M: make(map[string]types.Album),
	}
	//init the mutex
	store.RWMutex = &sync.RWMutex{}
	//filling the map
	for _, album := range albums {
		store.M[album.ID] = album
	}

	/** init http request multiplexer
	* in dictionary:
	* Multiplexer. A device that enables the simultaneous
	* transmission of several messages or signals over 
	* one communications channel.
	*
	* 
	* in GO:
	* 	In go ServeMux is an HTTP request multiplexer.
	* 	It matches the URL of each incoming request against a list of
	* 	registered patterns and calls the handler
	* 	for the pattern that most closely matches the URL.
	*
		type ServeMux struct {
	    mu    sync.RWMutex
	    m     map[string]muxEntry
	    es    []muxEntry
	    hosts bool
	}
	
	type muxEntry struct {
	    h       Handler
	    pattern string
	}
*/
	mux := http.NewServeMux()
	//adding handler to array
	mux.Handle("/album/", &handlers.AlbumHandler{Store: store})
	//accept and responde to client
	http.ListenAndServe(":8080", mux)
}
