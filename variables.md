base-url -> ex"localhost:8080"
serve: bool
serve-dir: if serve is enabled we serve the application from this directory with simple server
images-dir -> ex"images/" (default: "cyclopes")
multithreading: bool

---

## list

path -> ex"/blog"
device: (desktop)/mobile/both

---

1. Parse yml
2. Check if we need a server or not
3. If so open a server and serve the desired directory (new parallel thread)
4. Create chromedb instsance
5. Use this instance and start screenshoting in paraller
6. When all screenshots are done, close the server and exit applications

---

```
GOOS=linux GOARCH=amd64 go build -v -o cyclopes-linux-amd64 cmd/eidetic/*.go
```

```
apk add chromium
```

```
docker exec -it testing -v ${PWD}:/var/folder alpine /bin/bash
```

### Feedback

1. May have Cookies pop-up we want to disable
2. Check if directory already exists to save images else create
3. If the previous page is for example "/" and the next is "/#sda" the page is not reloading the Javascript script is raising errors like "scrollToBottomStepper" has already been declared
