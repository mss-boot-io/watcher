package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mss-boot-io/watcher"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	w := watcher.New()

	// Uncomment to use SetMaxEvents set to 1 to allow at most 1 event to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	// w.SetMaxEvents(1)

	// Uncomment to only notify rename and move events.
	// w.FilterOps(watcher.Rename, watcher.Move)

	// Uncomment to filter files based on a regular expression.
	//
	// Only files that match the regular expression during file listing
	// will be watched.
	// r := regexp.MustCompile("^abc$")
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add("example"); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive("example/test_folder"); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		log.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	// Trigger 2 events after watcher started.
	go func() {
		w.Wait()
		err := w.Ignore("example/test_folder/node_modules")
		if err != nil {
			log.Fatalln(err)
		}
		w.TriggerEvent(watcher.Create, nil)
		w.TriggerEvent(watcher.Remove, nil)
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
