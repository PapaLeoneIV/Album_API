package album_db

import (
	"api_mux/types"
	"sync"
)

type Db struct {
	M map[string]types.Album
	*sync.RWMutex
}