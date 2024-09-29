package watch

import (
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

var _ *fsnotify.Op

func Test_all(t *testing.T) {
	watch, err := NewWatch()
	if err != nil {
		t.Fatal(err)
	}
	defer watch.Close()
	tdir := os.TempDir()
	watch.Monitor(tdir, func(event fsnotify.Event) {
		t.Log(event.Name, event.Op.String())
	})

	tfile, err := os.CreateTemp(tdir, "*.tmp")
	if err != nil {
		t.Fatal(err)
	}
	b := []byte("123")
	tfile.Write(b)
	// tfile.Chmod(0o777)
	tfile.Close()
	os.Rename(tfile.Name(), tfile.Name()+"-")
	os.RemoveAll(tfile.Name() + "-")
	time.Sleep(time.Millisecond * 100)
}
