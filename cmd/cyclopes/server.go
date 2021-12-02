package cyclopes

import (
	"log"
	"net/http"
)

func Server(path string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	log.Println("Serving path: `" + path + "` on port :3000")

	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
