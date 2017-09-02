// go run counter.go -auto-created

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	autoCreate bool
	dbFile     string
	dbLockFile string
	lockValue  string
	saveMin    int

	counter Counter
)

const (
	pathIndex = "/"
	pathCount = "/count"
)

type Counter struct {
	count int64
	m     sync.RWMutex
}

func (c *Counter) Init(i int64) {
	c.m.Lock()
	defer c.m.Unlock()

	c.count = i
}

func (c *Counter) Inc() int64 {
	c.m.Lock()
	defer c.m.Unlock()

	c.count++
	return c.count
}

func (c *Counter) Get() int64 {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.count
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != pathIndex {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET", "HEAD", "OPTIONS":
	default:
		methodNotAllow(w, r)
		return
	}

	noCache(w.Header())
	w.Write([]byte(`Simple counter service`))
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != pathCount {
		http.NotFound(w, r)
		return
	}

	var c int64
	switch r.Method {
	case "GET", "HEAD", "OPTIONS":
		c = counter.Get()
	case "POST":
		c = counter.Inc()

	default:
		methodNotAllow(w, r)
		return
	}

	s := strconv.FormatInt(c, 10)
	w.Write([]byte(s))
}

func methodNotAllow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("405 method not allowed"))
	return
}

func noCache(h http.Header) {
	h.Add("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Add("Pragma", "no-cache")
	h.Add("Expires", "0")
}

func wrap(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			e := recover()
			if e != nil {
				log.Printf("Panic (recovered): %v", e)
				debug.PrintStack()
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 internal server error"))
			}
		}()
		h(w, r)
	}
}

func load() {
	var err error
	dbFile, err = filepath.Abs(dbFile)
	if err != nil {
		log.Fatalf("Invalid file path: %v", err)
	}
	dbLockFile = dbFile + ".lock"
	if _, err = os.Stat(dbLockFile); err == nil {
		log.Fatalf("Counter is corrupted: `%v.lock` exists.\nRemove the file to continue.", dbFile)
	}

	var c int64
	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		if !autoCreate {
			log.Fatalf("File not found `%v`.\nIf it is the first time, consider run with `-auto-create` argument.", dbFile)
		}

		log.Printf("Auto created file at `%v`", dbFile)

	} else {
		data, err := ioutil.ReadFile(dbFile)
		if err != nil {
			log.Fatalf("Unable to read file `%v`", dbFile)
		}
		c, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			log.Fatalf("Unable to load file `%v`", dbFile)
		}
		log.Printf("Loaded from file `%v`", dbFile)
		counter.Init(c)
	}

	save(true)
	log.Printf("Start with count %v", c)
}

func randString() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(b)
}

func save(withLock bool) {
	isLockExist := false
	if _, err := os.Stat(dbLockFile); err == nil {
		isLockExist = true
		data, err := ioutil.ReadFile(dbLockFile)
		if err != nil {
			log.Fatalf("Counter is corrupted: Unable to read file `%v`.\nRemove the file to continue.", dbLockFile)
		}
		if string(data) != lockValue {
			log.Fatalf("Counter is corrupted: Unmatched value at `%v`.\nRemove the file to continue.", dbLockFile)
		}
	}

	if withLock {
		lockValue = randString()
		err := ioutil.WriteFile(dbLockFile, []byte(lockValue), os.ModePerm)
		if err != nil {
			log.Fatalf("Unable to write to file `%v` (%v)", dbLockFile, err)
		}
		log.Print("Wrote to lock file")

	} else if isLockExist {
		defer func() {
			err := os.Remove(dbLockFile)
			if err != nil {
				log.Fatalf("Unable to remove lock file `%v` (%v)", dbLockFile, err)
			}
			log.Print("Removed lock file")
		}()
	}

	// Retrieve the counter value
	c := counter.Get()
	s := strconv.FormatInt(c, 10)
	err := ioutil.WriteFile(dbFile, []byte(s), os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to write to file `%v` (%v)", dbFile, err)
	}
	log.Println("Saved to storage")
}

func main() {
	flListen := flag.String("listen", ":8901", "Address to listen on")
	flag.StringVar(&dbFile, "file", ".counter.db", "Path to counter file")
	flag.BoolVar(&autoCreate, "auto-create", false, "Auto create file if not found")
	flag.IntVar(&saveMin, "save-min", 10, "Save to storage for each x minutes")
	flag.Parse()

	if saveMin <= 0 || saveMin >= 1000 {
		log.Fatalf("Invalid minutes: %v", saveMin)
	}

	load()
	defer save(false)

	// Listen to os signal
	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()

		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("Received OS signal %v", <-osSignal)
	}()

	// Save each 10 minutes
	go func() {
		t := time.NewTicker(time.Duration(saveMin) * time.Minute)
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			save(true)
		}
	}()

	http.HandleFunc(pathIndex, wrap(indexHandler))
	http.HandleFunc(pathCount, wrap(countHandler))

	// Start HTTP server
	svr := &http.Server{Addr: *flListen}
	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := svr.Shutdown(ctx2)
	if err != nil {
		log.Fatal(err)
	}
}
