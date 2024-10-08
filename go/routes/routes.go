package routes

import (
	"log"
	"net/http"
	"os"

	"g.rg-s.com/org/go/api/files"
	"g.rg-s.com/org/go/api/kine"
	"g.rg-s.com/org/go/api/nav"
	"g.rg-s.com/org/go/api/search"
	"g.rg-s.com/org/go/api/today"
	"g.rg-s.com/org/go/api/topics"
	"g.rg-s.com/org/go/api/view"
	"g.rg-s.com/org/go/helper/reqerr"
	"g.rg-s.com/org/go/index"

	"github.com/gorilla/mux"
)

func Routes(ix *index.Index) *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/file/").HandlerFunc(hix(ix, view.ServeStatic))

	r.HandleFunc("/api/today", hix(ix, today.GetToday))
	r.HandleFunc("/api/now", hix(ix, today.GetNow))

	r.PathPrefix("/api/nav").HandlerFunc(hix(ix, nav.ViewNav))

	r.PathPrefix("/api/view").HandlerFunc(hix(ix, view.ViewFile))
	r.PathPrefix("/api/copy").HandlerFunc(hix(ix, files.CopyFile))
	r.PathPrefix("/api/move").HandlerFunc(hix(ix, files.RenameFile))
	r.PathPrefix("/api/write").HandlerFunc(hix(ix, files.WriteSwitch))
	r.PathPrefix("/api/delete").HandlerFunc(hix(ix, files.DeleteFile))
	r.PathPrefix("/api/sort").HandlerFunc(hix(ix, files.WriteSort))
	r.PathPrefix("/api/search").HandlerFunc(hix(ix, search.SearchFiles))
	r.PathPrefix("/api/topics/").HandlerFunc(hix(ix, topics.ViewTopic))
	r.PathPrefix("/api/topics").HandlerFunc(hix(ix, topics.Topics))
	r.PathPrefix("/api/kines").HandlerFunc(hix(ix, kine.Kines))
	r.PathPrefix(kine.CreateAPI).HandlerFunc(hix(ix, kine.Create))
	r.PathPrefix(kine.UploadAPI).HandlerFunc(hix(ix, kine.Upload))
	r.Path(kine.TalkAPI).HandlerFunc(hix(ix, kine.Talk))
	r.Path("/api/kine/listen").HandlerFunc(hix(ix, kine.Listen))

	r.PathPrefix("/rl/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ix.Load()
	})

	r.PathPrefix("/").HandlerFunc(buildWrapper(ix, serveBuild))

	return r
}

func buildWrapper(ix *index.Index, fn func(*index.Index, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(ix, w, r)
	}
}

func serveBuild(ix *index.Index, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		file := ix.Build + r.URL.Path
		_, err := os.Stat(file)
		if err == nil {
			http.ServeFile(w, r, file)
			return
		}
	}
	http.ServeFile(w, r, ix.Build+"/index.html")
}

func hix(ix *index.Index, fn func(*index.Index, http.ResponseWriter, *http.Request) *reqerr.Err) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(ix, w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), err.Code)
		}
	}
}
