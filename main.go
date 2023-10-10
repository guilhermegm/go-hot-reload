package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var cmdStr string
	flag.StringVar(&cmdStr, "cmd", "", "Command to run")
	var dirStr string
	flag.StringVar(&dirStr, "dir", "", "Directory to run the command in")
	flag.Parse()

	cmds := []*Cmd{
		{
			Ch:         make(chan error, 1),
			RawCommand: cmdStr,
			Directory:  dirStr,
		},
	}

	startCommands(cmds)

	handleSignals(cmds)

	watcher, err := setupFileWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watchForChanges(watcher, cmds)
}

func startCommands(cmds []*Cmd) {
	for _, cmd := range cmds {
		go func(c *Cmd) {
			c.Prepare()
			err := c.Cmd.Run()
			c.Ch <- err
		}(cmd)
	}
}

func handleSignals(cmds []*Cmd) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh
		terminateCommands(cmds)
		os.Exit(1)
	}()
}

func setupFileWatcher() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".git") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return watcher, nil
}

func watchForChanges(watcher *fsnotify.Watcher, cmds []*Cmd) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				terminateCommands(cmds)
				startCommands(cmds)
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				fileInfo, err := os.Stat(event.Name)
				if err != nil {
					log.Println("Error fetching info:", err)
					continue
				}

				if fileInfo.IsDir() {
					err = watcher.Add(event.Name)
					if err != nil {
						log.Println("Error adding created directory to watcher:", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if ok {
				log.Println("error:", err)
			}
		}
	}
}

func terminateCommands(cmds []*Cmd) {
	for _, cmd := range cmds {
		if err := syscall.Kill(-cmd.Cmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Printf("Failed to kill process group: %v", err)
		} else {
			fmt.Println("Reloading...")
		}
	}
}

type Cmd struct {
	Cmd        *exec.Cmd
	Ch         chan error
	RawCommand string
	Directory  string
}

func (c *Cmd) Prepare() {
	c.Cmd = exec.Command("sh", "-c", c.RawCommand)
	c.Cmd.Dir = c.Directory
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr
	c.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
