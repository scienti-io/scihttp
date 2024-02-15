package scihttp

import (
	"fmt"
	"net/http"
)

func StandardHTTPError(w http.ResponseWriter, sc int) {
	msg := fmt.Sprintf("%d %s", sc, http.StatusText(sc))
	http.Error(w, msg, sc)
}

func NotImplementedError(w http.ResponseWriter) {
	StandardHTTPError(w, http.StatusNotImplemented)
}
