package watcher

import (
	"io/fs"
	"os"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
)

type Op watcher.Op

const (
	FILE_CREATED  Op = Op(watcher.Create)
	FILE_DELETED  Op = Op(watcher.Remove)
	FILE_MODIFIED Op = Op(watcher.Write)
)

var filewatcher *fileWatcher = nil

type Config struct {
	PathToWatch string
	Ops         []Op // TODO: use map to avoid duplicates
	// Frequency represents the polling cycle which repeats every specified duration until Close is called.
	Frequency time.Duration
}

type Event struct {
	watcher.Event
}

func NewEvent(op Op, path string, oldPath string, f os.FileInfo) *Event {
	return &Event{
		watcher.Event{
			Op:       watcher.Op(op),
			Path:     path,
			OldPath:  oldPath,
			FileInfo: f,
		},
	}
}

type fileWatcher struct {
	w      *watcher.Watcher
	config Config
}

// NewFileWatcher creates a new FileWatcher instance if it doesn't already exist
func NewFileWatcher(config Config) *fileWatcher {
	var lock = &sync.Mutex{}

	if filewatcher == nil {
		lock.Lock()
		defer lock.Unlock()
		if filewatcher == nil {

			// create a new file watcher
			w := watcher.New()

			// adding events to the watcher
			{
				filterOps := []watcher.Op{}
				for _, e := range config.Ops {
					filterOps = append(filterOps, watcher.Op(e))
				}

				w.FilterOps(filterOps...)
			}

			// adding filedir to be watched
			w.AddRecursive(config.PathToWatch)

			// TODO: get path from config
			w.Ignore("D:\\Projects\\github.com\\pramod-janardhana\\dir-watcher\\file_event.sqlite")
			w.Ignore("D:\\Projects\\github.com\\pramod-janardhana\\dir-watcher\\job_runs.sqlite")

			filewatcher = &fileWatcher{
				w:      w,
				config: config,
			}
		}
	}

	return filewatcher
}

func GetFileWatcher() *fileWatcher {
	return filewatcher
}

// Start begins the polling cycle which repeats every specified duration until Close is called.
func (w *fileWatcher) Start() error {
	return w.w.Start(w.config.Frequency)
}

// Event returns an event channel
func (w *fileWatcher) Event() chan watcher.Event {
	return w.w.Event
}

func (w *fileWatcher) Error() chan error {
	return w.w.Error
}

// Close stops a Watcher and unlocks its mutex
func (w *fileWatcher) Close() {
	w.w.Close()
	// <-w.w.Closed
	filewatcher = nil
}

// Closed returns a channel that signals when the watcher is closed.
func (w *fileWatcher) Closed() chan struct{} {
	return w.w.Closed
}

// WatchedFiles returns a map of files/dirs added to a Watcher.
func (w *fileWatcher) WatchedFiles() map[string]fs.FileInfo {
	return w.w.WatchedFiles()
}

// ReloadConfig closes the previous file watcher and starts a new one with the given configuration.
func (w *fileWatcher) ReloadConfig(config Config) *fileWatcher {
	filewatcher.w.Close()
	filewatcher = nil
	return NewFileWatcher(config)
}

// func main() {
// 	w := watcher.New()

// 	// SetMaxEvents to 1 to allow at most 1 event's to be received
// 	// on the Event channel per watching cycle.
// 	//
// 	// If SetMaxEvents is not set, the default is to send all events.
// 	// w.SetMaxEvents(1)

// 	// Only notify rename and move events.
// 	w.FilterOps(watcher.Remove, watcher.Create, watcher.Write)

// 	// Only files that match the regular expression during file listings
// 	// will be watched.
// 	r := regexp.MustCompile(".*")
// 	w.AddFilterHook(watcher.RegexFilterHook(r, true))

// 	go func() {
// 		for {
// 			select {
// 			case event := <-w.Event:
// 				read(event.Path)
// 				fmt.Println(event) // Print the event's info.
// 			case err := <-w.Error:
// 				log.Fatalln(err)
// 			case <-w.Closed:
// 				return
// 			}
// 		}
// 	}()

// 	// Watch this folder for changes.
// 	if err := w.AddRecursive("."); err != nil {
// 		log.Fatalln(err)
// 	}

// 	// Print a list of all of the files and folders currently
// 	// being watched and their paths.
// 	for path, f := range w.WatchedFiles() {
// 		fmt.Printf("%s: %s\n", path, f.Name())
// 	}

// 	fmt.Println()

// 	// Trigger 2 events after watcher started.
// 	go func() {
// 		w.Wait()
// 		w.TriggerEvent(watcher.Create, nil)
// 		w.TriggerEvent(watcher.Remove, nil)
// 	}()

// 	// Start the watching process - it'll check for changes every 100ms.
// 	if err := w.Start(time.Millisecond * 100); err != nil {
// 		log.Fatalln(err)
// 	}
// }
