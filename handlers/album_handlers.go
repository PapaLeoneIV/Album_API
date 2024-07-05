package handlers


import (
	"api_mux/regex"
	"api_mux/types"
	"api_mux/db"
	"net/http"
	"encoding/json"
)

type AlbumHandler struct {
	Store *album_db.Db
}

func (h *AlbumHandler) Update(w http.ResponseWriter, req *http.Request){
	matches := regxp.GetAlbumRe.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}
	var buf types.Album
	if err := json.NewDecoder(req.Body).Decode(&buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}


	h.Store.Lock()
	defer h.Store.Unlock()

	if _, exist := h.Store.M[matches[1]]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Album not found!"))
		return
	}

	h.Store.M[matches[1]] = buf
	w.WriteHeader(http.StatusOK)
}

func (h *AlbumHandler) Delete(w http.ResponseWriter, req *http.Request){
	matches := regxp.GetAlbumRe.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}
	h.Store.Lock()
	
	if _, exist := h.Store.M[matches[1]]; !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Album not found!\n"))
		return
	}
	
	delete(h.Store.M, matches[1])
	h.Store.Unlock()
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deletion was successfull!"))
}

func (h *AlbumHandler) List(w http.ResponseWriter, req *http.Request) {
	h.Store.RLock()
	arr := make([]types.Album, 0, len(h.Store.M))
	for _, value := range h.Store.M {
		arr = append(arr, value)
	}
	h.Store.RUnlock()
	prettyJSON, err := json.Marshal(arr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)
}

func (h *AlbumHandler) Get(w http.ResponseWriter, req *http.Request) {
	matches := regxp.GetAlbumRe.FindStringSubmatch(req.URL.Path)
	if len(matches) < 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
	h.Store.RLock()
	album, ok := h.Store.M[matches[1]]
	h.Store.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("album not found"))
		return
	}
	prettyJSON, err := json.Marshal(album)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)
}

func (h *AlbumHandler) Add(w http.ResponseWriter, req *http.Request) {
	var buf types.Album
	if err := json.NewDecoder(req.Body).Decode(&buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	h.Store.Lock()
	h.Store.M[buf.ID] = buf
	h.Store.Unlock()
	prettyJSON, err := json.Marshal(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJSON)
}

func (h *AlbumHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case req.Method == "GET" && regxp.GetAlbumRe.MatchString(req.URL.Path):
		h.Get(w, req)
	case req.Method == "GET" && regxp.ListAlbumsRe.MatchString(req.URL.Path):
		h.List(w, req)
	case req.Method == "POST" && regxp.CreateAlbumRe.MatchString(req.URL.Path):
		h.Add(w, req)
	case req.Method == "DELETE" && regxp.DeleteAlbumRe.MatchString(req.URL.Path):
		h.Delete(w, req)
	case req.Method == "PUT" && regxp.UpdateAlbumre.MatchString(req.URL.Path):
		h.Update(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}
}
