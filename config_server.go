package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed server/build
var staticFS embed.FS

var (
	mouseConfigDict    = make(map[byte]string)
	keyboardConfigDict = make(map[byte]string)
	preConfigDict      = make(map[string][2]map[byte]string)
	activePreConfig    string
	mousedictMutex     sync.RWMutex
	keyboarddictMutex  sync.RWMutex
)

func serve(port int) {
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

	http.HandleFunc("/api/get/active", func(w http.ResponseWriter, r *http.Request) {
		mousedictMutex.RLock()
		defer mousedictMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"active": activePreConfig})
	})

	http.HandleFunc("/api/set/active", func(w http.ResponseWriter, r *http.Request) {
		mousedictMutex.Lock()
		defer mousedictMutex.Unlock()
		activePreConfig = r.URL.Query().Get("name")
		logger.Infof("Set active preConfig: %s", activePreConfig)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/get/recoilInput", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./config/手动设置压枪数据.txt")
		if err != nil {
			if os.IsNotExist(err) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Write([]byte(""))
				return
			}
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})

	http.HandleFunc("/api/set/recoilInput", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		err = os.WriteFile("./config/手动设置压枪数据.txt", body, 0644)
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}
		logger.Infof("更新自定义压枪数据: %d bytes", len(body))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
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
		keyboarddictMutex.Lock()
		defer keyboarddictMutex.Unlock()
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

	// 配置文件管理 API
	profilesDir := "./config/profiles"
	os.MkdirAll(profilesDir, 0755)

	http.HandleFunc("/api/get/configs", func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir(profilesDir)
		if err != nil {
			http.Error(w, "Failed to read profiles directory", http.StatusInternalServerError)
			return
		}
		var names []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
				names = append(names, strings.TrimSuffix(e.Name(), ".yaml"))
			}
		}
		sort.Strings(names)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(names)
	})

	http.HandleFunc("/api/get/config", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing name parameter", http.StatusBadRequest)
			return
		}
		data, err := os.ReadFile(filepath.Join(profilesDir, name+".yaml"))
		if err != nil {
			http.Error(w, "Failed to read config", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})

	http.HandleFunc("/api/set/config", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing name parameter", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		err = os.WriteFile(filepath.Join(profilesDir, name+".yaml"), body, 0644)
		if err != nil {
			http.Error(w, "Failed to write config", http.StatusInternalServerError)
			return
		}
		logger.Infof("保存配置: %s", name)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/delete/config", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing name parameter", http.StatusBadRequest)
			return
		}
		err := os.Remove(filepath.Join(profilesDir, name+".yaml"))
		if err != nil {
			http.Error(w, "Failed to delete config", http.StatusInternalServerError)
			return
		}
		logger.Infof("删除配置: %s", name)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/apply/config", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing name parameter", http.StatusBadRequest)
			return
		}
		data, err := os.ReadFile(filepath.Join(profilesDir, name+".yaml"))
		if err != nil {
			http.Error(w, "Config not found", http.StatusNotFound)
			return
		}
		// 写入 using.yaml（程序读取）
		if err := os.WriteFile("./config/using.yaml", data, 0644); err != nil {
			http.Error(w, "Failed to apply config", http.StatusInternalServerError)
			return
		}
		// 记录当前配置名
		os.WriteFile("./config/using.txt", []byte(name), 0644)
		logger.Infof("应用配置: %s", name)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/api/get/activeConfig", func(w http.ResponseWriter, r *http.Request) {
		nameBytes, err := os.ReadFile("./config/using.txt")
		if err != nil {
			http.Error(w, "Failed to read active config name", http.StatusInternalServerError)
			return
		}
		name := strings.TrimSpace(string(nameBytes))
		data, err := os.ReadFile(filepath.Join(profilesDir, name+".yaml"))
		if err != nil {
			http.Error(w, "Failed to read active config", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(data)
	})

	http.HandleFunc("/api/get/activeConfigName", func(w http.ResponseWriter, r *http.Request) {
		nameBytes, err := os.ReadFile("./config/using.txt")
		if err != nil {
			http.Error(w, "Failed to read active config name", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(strings.TrimSpace(string(nameBytes))))
	})

	http.HandleFunc("/api/restart", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("收到重启请求，程序将在100ms后退出")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		go func() {
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}()
	})

	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	logger.Info("可从以下网址访问控制后台:")
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ipv4 := ipNet.IP.To4()
			if ipv4 == nil {
				continue // 跳过非 IPv4 地址
			}
			if !ipv4.IsLoopback() {
				logger.Infof("http://%s:%v", ipv4, port)
			}
		}

	}
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
