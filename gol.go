package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ogier/pflag"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"./auth"
	_ "./auth/insecure"
	_ "./auth/ldap"
	"./post"
	"./storage"
	_ "./storage/gol"
	_ "./storage/json"
	_ "./storage/memory"
	_ "./storage/multi"
	_ "./storage/sqlite"
	"./templates"
)

func toByteSlice(data interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		log.Println("Error: binary.Write failed:", err)
		return []byte{}
	}
	return buf.Bytes()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func renderPosts(templates *template.Template, w http.ResponseWriter, posts []post.Post) {
	m := make(map[string]interface{})
	m["title"] = "gol"
	m["posts"] = posts
	templates.ExecuteTemplate(w, "posts", m)
}

func createPost(title, content string) post.Post {
	now := time.Now()
	return post.Post{
		Id:      fmt.Sprintf("%x", md5.Sum(toByteSlice(now.UnixNano()))),
		Title:   title,
		Content: content,
		Created: now,
	}
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func notImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

func refererRedirectPath(r *http.Request, defaultPath string) string {
	referer := r.Referer()
	if referer != "" {
		refererUrl, err := url.Parse(referer)
		if err != nil {
			return defaultPath
		}
		if refererUrl.Host != r.URL.Host {
			return defaultPath
		}
		return refererUrl.Path
	} else {
		return defaultPath
	}
}

func urlHasQuery(u *url.URL) bool {
	q := u.Query()
	if len(q) == 0 {
		return false
	}

	queryParams := []string{"id", "title", "start", "end", "sort", "reverse", "match", "range"}
	for _, p := range queryParams {
		if _, ok := q[p]; ok {
			return true
		}
	}

	return false
}

func queryFromURL(u *url.URL, store storage.Store) ([]post.Post, error) {
	defaultQuery, _ := storage.Query().Reverse().Build()
	if !urlHasQuery(u) {
		return store.Find(*defaultQuery)
	}

	q, err := storage.QueryFromURL(u)
	if err != nil {
		return nil, err
	}
	return store.Find(*q)
}

func newSession(sessions map[string]string, username string) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println(err)
	}
	sessionId := base64.StdEncoding.EncodeToString(randomBytes)
	sessions[sessionId] = username
	return sessionId
}

// returns (user, true) if valid session id
func hasSession(sessions map[string]string, sessionId string) (string, bool) {
	if username, ok := sessions[sessionId]; ok {
		return username, true
	}

	return "", false
}

func isLoggedIn(sessions map[string]string, r *http.Request) bool {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	_, ok := hasSession(sessions, sessionCookie.Value)
	return ok
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("/login?redirect_to=%s", url.QueryEscape(r.URL.Path))
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func getTemplateBasePath(customLocation string) (string, error) {
	userHome := getEnv("HOME", "")
	workingDir := getEnv("PWD", "")
	basePaths := []string{customLocation,
		path.Join(userHome, "/.local/share/gol/templates"),
		"/usr/share/gol/templates",
		path.Join(workingDir, "templates")}
	basePathFound := ""
	for _, basePath := range basePaths {
		_, err := os.Stat(basePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return "", nil
			}
		} else {
			basePathFound = basePath
			break
		}
	}

	return basePathFound, nil
}

var Environment = getEnv("ENVIRONMENT", "development")
var Version = "master"
var templateBase = pflag.String("templates",
	"",
	"templates path")
var assetBase = pflag.String("assets",
	fmt.Sprintf("https://cdn.rawgit.com/klingtnet/gol/%s/assets", Version),
	"assets path")
var ssl = pflag.String("ssl",
	"",
	"enable ssl (give server.crt,server.key as value)")
var storageUrl = pflag.String("storage",
	"json://posts.json",
	"the storage to connect to")
var authUrl = pflag.String("authentication",
	"",
	"the authentication method to use")

func init() {
	if Environment == "production" {
		//TODO
	}

	fmt.Printf("gol - %s (%s)\n", Version, Environment)
}

func main() {
	pflag.Parse()

	var store storage.Store
	storageUrls := strings.Split(*storageUrl, ",")
	if len(storageUrls) > 1 {
		multiUrl := fmt.Sprintf("multi://?primary=%s", url.QueryEscape(storageUrls[0]))
		for _, storageUrl := range storageUrls[1:] {
			multiUrl = fmt.Sprintf("%s&secondary=%s", multiUrl, url.QueryEscape(storageUrl))
		}
		var err error
		store, err = storage.Open(multiUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		store, err = storage.Open(*storageUrl)
		if err != nil {
			log.Fatal(err)
		}
	}

	var authenticator auth.Auth
	if authUrl != nil && *authUrl != "" {
		a, err := auth.Open(*authUrl)
		if err != nil {
			log.Fatal(err)
		}
		authenticator = a
	}

	// username -> session
	sessions := map[string]string{}

	templBasePath, err := getTemplateBasePath(*templateBase)
	if err != nil || templBasePath == "" {
		log.Print("Could not get template base path!")
		log.Fatal(err)
	}

	fmt.Printf("Using template base path: %s\n", templBasePath)
	templates := templates.Templates(templBasePath, *assetBase)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		posts, err := queryFromURL(r.URL, store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			renderPosts(templates, w, posts)
		}
	})

	if authenticator != nil {
		router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
			if isLoggedIn(sessions, r) {
				redirectPath := refererRedirectPath(r, "/")
				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
				return
			}

			if r.Method == "GET" {
				templates.ExecuteTemplate(w, "login", map[string]string{"title": "Login"})
			} else if r.Method == "POST" {
				username := r.FormValue("username")
				password := r.FormValue("password")
				err := authenticator.Login(username, password)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
				} else {
					http.SetCookie(w, &http.Cookie{
						Name:  "session",
						Value: newSession(sessions, username),
					})

					redirectPath := r.URL.Query().Get("redirect_to")
					if redirectPath == "" {
						redirectPath = "/"
					}
					http.Redirect(w, r, redirectPath, http.StatusSeeOther)
				}
			} else {
				notImplemented(w)
			}
		})

		router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				sessionCookie, err := r.Cookie("session")
				if err == nil {
					delete(sessions, sessionCookie.Value)
					http.SetCookie(w, &http.Cookie{
						Name:   "session",
						Value:  "",
						MaxAge: -1,
					})
				}

				redirectPath := refererRedirectPath(r, "/")
				http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			}
		})
	}

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		posts, err := queryFromURL(r.URL, store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(posts)
		}
	}).Methods("GET").Headers("Content-Type", "application/json")

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			posts, err := queryFromURL(r.URL, store)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				renderPosts(templates, w, posts)
			}
		} else if r.Method == "POST" { // POST creates a new post
			isJson := strings.Contains(r.Header.Get("Content-Type"), "application/json")

			if authenticator != nil && !isLoggedIn(sessions, r) {
				redirectToLogin(w, r)
				return
			}

			var post post.Post
			if isJson {
				json.NewDecoder(r.Body).Decode(&post)
				post = createPost(post.Title, post.Content)
			} else {
				post = createPost(r.FormValue("title"), r.FormValue("content"))
			}

			err := store.Create(post)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			if isJson {
				w.WriteHeader(http.StatusAccepted)
				writeJson(w, post)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/new", func(w http.ResponseWriter, r *http.Request) {
		if authenticator != nil && !isLoggedIn(sessions, r) {
			redirectToLogin(w, r)
			return
		}

		templates.ExecuteTemplate(w, "post_form", map[string]string{"title": "Write a new post!"})
	})

	router.HandleFunc("/posts/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // use post to receive content in body
			var post post.Post
			err := json.NewDecoder(r.Body).Decode(&post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				templates.ExecuteTemplate(w, "post", post)
			}
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		p, _ := store.FindById(id)
		if p == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if r.Method == "GET" {
			if r.Header.Get("Content-Type") == "application/json" {
				writeJson(w, p)
			} else {
				m := make(map[string]interface{})
				m["title"] = p.Title
				m["posts"] = []post.Post{*p}
				templates.ExecuteTemplate(w, "posts", m)
			}
		} else if r.Method == "HEAD" {
			// already handle by p == nil above
		} else if r.Method == "POST" {
			if authenticator != nil && !isLoggedIn(sessions, r) {
				redirectToLogin(w, r)
				return
			}

			var newPost post.Post
			if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				newPost.Title = r.FormValue("title")
				newPost.Content = r.FormValue("content")

				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else { // assume it's JSON
				err := json.NewDecoder(r.Body).Decode(&newPost)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusAccepted)
				writeJson(w, newPost)
			}

			if newPost.Title != "" {
				p.Title = newPost.Title
			}
			if newPost.Content != "" {
				p.Content = newPost.Content
			}
			store.Update(*p)
			json.NewEncoder(w).Encode(p)
		} else if r.Method == "DELETE" {
			if authenticator != nil && !isLoggedIn(sessions, r) {
				redirectToLogin(w, r)
				return
			}

			err := store.Delete(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		if authenticator != nil && !isLoggedIn(sessions, r) {
			redirectToLogin(w, r)
			return
		}

		id := mux.Vars(r)["id"]
		post, _ := store.FindById(id)
		if post != nil {
			m := make(map[string]interface{})
			m["title"] = "Edit post"
			m["post"] = post
			templates.ExecuteTemplate(w, "post_form", m)
		} else {
			http.NotFound(w, r)
		}
	})

	if Environment == "development" {
		// in development, serve local assets
		router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	}

	http.Handle("/", router)

	host := getEnv("HOST", "localhost")
	port := getEnv("PORT", "5000")
	addr := fmt.Sprintf("%s:%s", host, port)
	if *ssl == "" {
		fmt.Printf("Listening on http://%s\n", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	} else {
		certAndKey := strings.Split(*ssl, ",")
		if len(certAndKey) != 2 {
			fmt.Println("Error: -ssl needs server.crt,server.key as arguments")
			os.Exit(1)
		}
		fmt.Printf("Listening on https://%s\n", addr)
		log.Fatal(http.ListenAndServeTLS(addr, certAndKey[0], certAndKey[1], nil))
	}
}
