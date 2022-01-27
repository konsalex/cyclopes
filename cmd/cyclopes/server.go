package cyclopes

import (
	"errors"
	"net/http"
	"os"

	"github.com/pterm/pterm"
)

// Server spins-up a local server on :3000 to
// serve the static files, if the Website is not published
// yet nor served from a local server
func Server(path string) *http.Server {
	// If path does not exist, throw
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	fs := http.FileServer(http.Dir(path))

	pterm.Info.Println("Serving path: `" + path + "` on port :3000")
	server := &http.Server{Addr: ":3000", Handler: fs}

	go func(server *http.Server) {
		defer recover()
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}(server)

	return server
}
