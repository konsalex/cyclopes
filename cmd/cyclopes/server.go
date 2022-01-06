package cyclopes

import (
	"net/http"

	"github.com/pterm/pterm"
)

func Server(path string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	pterm.Info.Println("Serving path: `" + path + "` on port :3000")
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			panic(err)
		}
	}()
}
