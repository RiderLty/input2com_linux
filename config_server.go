package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"sync"
)

//go:embed server/build
var staticFS embed.FS

var (
	mouseConfigDict    = make(map[byte]string)
	keyboardConfigDict = make(map[byte]string)
	preConfigDict      = make(map[string][2]map[byte]string)
	mousedictMutex     sync.RWMutex
	keyboarddictMutex  sync.RWMutex
)

func serve(port int) {
	logger.Infof(fmt.Sprintf("Starting config server on port %d", port))
	webFS, err := fs.Sub(staticFS, "server/build")
	if err != nil {
		logger.Errorf("无法加载静态文件: %v", err)
		return
	}
	http.Handle("/", http.FileServer(http.FS(webFS)))

	http.HandleFunc("/api/get/macros", func(w http.ResponseWriter, r *http.Request) {
		keyboarddictMutex.RLock()
		defer keyboarddictMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(macros); err != nil {
			http.Error(w, "Failed to encode keyboard config", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/get/preConfig", func(w http.ResponseWriter, r *http.Request) {
		mousedictMutex.RLock()
		defer mousedictMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(preConfigDict); err != nil {
			http.Error(w, "Failed to encode preConfig", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/get/mouse", func(w http.ResponseWriter, r *http.Request) {
		//return josn mouseConfigDict
		mousedictMutex.RLock()
		defer mousedictMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mouseConfigDict); err != nil {
			http.Error(w, "Failed to encode mouse config", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/get/keyboard", func(w http.ResponseWriter, r *http.Request) {
		keyboarddictMutex.RLock()
		defer keyboarddictMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(keyboardConfigDict); err != nil {
			http.Error(w, "Failed to encode keyboard config", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/set/mouse", func(w http.ResponseWriter, r *http.Request) {
		mousedictMutex.Lock()
		defer mousedictMutex.Unlock()
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "CLEAR_ALL" {
			mouseConfigDict = make(map[byte]string)
			logger.Infof("clear mouse config")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			if _, ok := MouseValidKeys[key]; !ok {
				http.Error(w, "Invalid key", http.StatusBadRequest)
				return
			}
			if value == "CLEAR_FUNCTION" {
				bkey, _ := strconv.ParseUint(key, 10, 8)
				logger.Infof("clear mouse config: %d", bkey)
				delete(mouseConfigDict, byte(bkey))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			} else {
				if _, ok := macros[value]; !ok {
					http.Error(w, "Invalid macro Name", http.StatusBadRequest)
					return
				}
				bkey, _ := strconv.ParseUint(key, 10, 8)
				logger.Infof("Set mouse config: %d -> %s", bkey, value)
				mouseConfigDict[byte(bkey)] = value
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			}
		}
	})

	http.HandleFunc("/api/set/keyboard", func(w http.ResponseWriter, r *http.Request) {
		mousedictMutex.Lock()
		defer mousedictMutex.Unlock()
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "CLEAR_ALL" {
			keyboardConfigDict = make(map[byte]string)
			logger.Infof("clear keyboard config")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			if _, ok := KeyboardValidKeys[key]; !ok {
				http.Error(w, "Invalid key", http.StatusBadRequest)
				return
			}
			if value == "CLEAR_FUNCTION" {
				bkey, _ := strconv.ParseUint(key, 10, 8)
				logger.Infof("clear keyboard config: %d", bkey)
				delete(keyboardConfigDict, byte(bkey))
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			} else {
				if _, ok := macros[value]; !ok {
					http.Error(w, "Invalid macro Name", http.StatusBadRequest)
					return
				}
				bkey, _ := strconv.ParseUint(key, 10, 8)
				logger.Infof("Set keyboard config: %d -> %s", bkey, value)
				keyboardConfigDict[byte(bkey)] = value
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			}
		}

	})

	logger.Infof("Config server started ")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
