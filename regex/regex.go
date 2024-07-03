package regxp

import "regexp"

var (
    ListAlbumsRe  = regexp.MustCompile(`^\/album[\/]*$`)
    GetAlbumRe    = regexp.MustCompile(`^\/album\/(\d+)$`)
    CreateAlbumRe = regexp.MustCompile(`^\/album[\/]*$`)
)
