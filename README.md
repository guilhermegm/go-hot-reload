# GoHR - Go Hot Reload

GoHR is a lightweight utility for Go, created using `fsnotify`, for personal use to streamline my development process. It watches for file changes in your project and automatically triggers a command of your choice. Though originally designed for hot-reloading Go applications, it's versatile enough for various tasks.

Installation
---

Install GoHR with:

```bash
go install github.com/guilhermegm/go-hot-reload@latest
```

Usage
---

Basic usage:

```bash
gohr -cmd "YOUR_COMMAND" -dir "YOUR_DIRECTORY"
```

For instance, if you're working on a Go project in the `cmd/api` directory and want to run migrations every time there's a change:

```bash
gohr -cmd "go run . -migrate" -dir "cmd/api"
```

### Examples of Usage

#### 1. Basic Go Project

Folder Structure:
```
/my-go-project
|-- main.go
|-- go.mod
|-- go.sum
```

Command:
```
gohr -cmd "go run ."
```

#### 2. Go Project with Multiple Services

Folder Structure:
```
/my-multi-service-project
|-- api
|   |-- main.go
|-- worker
|   |-- main.go
|-- go.mod
|-- go.sum
```

Command:
```
gohr -cmd "go run ." -dir "worker"

gohr -cmd "go run ." -dir "api"
```

#### 3. Frontend Development

Folder Structure:
```
/my-react-app
|-- public
|-- src
|   |-- App.js
|   |-- index.js
|-- package.json
```

Command:
```
gohr -cmd "npm start"
```

Flags
---

- __cmd__: The command you want to run on file changes.
- __dir__ (optional): The directory where the specified command will be executed. GoHR, however, will be diligently watching for file changes across the entire project.

License
---

### Acknowledgments

This project uses the [fsnotify](https://github.com/fsnotify/fsnotify) library for file system notifications.
