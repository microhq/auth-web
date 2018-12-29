package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"golang.org/x/net/context"

	//proto "github.com/micro/go-platform/auth/proto"
	account "github.com/micro/auth-srv/proto/account"
)

var (
	opts        *ace.Options
	accClient account.AccountClient
)

func Init(dir string, t account.AccountClient) {
	accClient = t

	opts = ace.InitializeOptions(nil)
	opts.BaseDir = dir
	opts.DynamicReload = true
	opts.FuncMap = template.FuncMap{
		"TimeAgo": func(t int64) string {
			return timeAgo(t)
		},
		"Timestamp": func(t int64) string {
			return time.Unix(t, 0).Format("02 Jan 06 15:04:05 MST")
		},
	}
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	basePath := hostPath(r)

	opts.FuncMap["URL"] = func(path string) string {
		return filepath.Join(basePath, path)
	}

	tpl, err := ace.Load("layout", tmpl, opts)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
		return
	}

        if data == nil {
                data = make(map[string]interface{})
        }

        data["Alert"] = getAlert(w, r)

	if err := tpl.Execute(w, data); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
	}
}

// The index page
func Index(w http.ResponseWriter, r *http.Request) {
	rsp, err := accClient.Search(context.TODO(), &account.SearchRequest{
//		Reverse: true,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	sort.Sort(sortedRecords{records: rsp.Accounts, reverse: true})

	render(w, r, "index", map[string]interface{}{
		"Latest": rsp.Accounts,
	})
}

func Accounts(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limit := 15

	page, err := strconv.Atoi(r.Form.Get("p"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	offset := (page * limit) - limit

	rsp, err := accClient.Search(context.TODO(), &account.SearchRequest{
//		Reverse: true,
		Limit:   int64(limit),
		Offset:  int64(offset),
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	var less, more int
	if len(rsp.Accounts) == limit {
		more = page + 1
	}

	if page > 1 {
		less = page - 1
	}

	sort.Sort(sortedRecords{records: rsp.Accounts, reverse: false})

	render(w, r, "accounts", map[string]interface{}{
		"Accounts": rsp.Accounts,
		"Less":   less,
		"More":   more,
	})
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id := r.Form.Get("id")
		typ := r.Form.Get("type")
		rsp, err := accClient.Search(context.TODO(), &account.SearchRequest{
			ClientId: id,
			Type: typ,
//			Reverse: true,
		})
		if err != nil {
			http.Redirect(w, r, filepath.Join(hostPath(r), "search"), 302)
			return
		}

		sort.Sort(sortedRecords{records: rsp.Accounts, reverse: false})

		render(w, r, "results", map[string]interface{}{
			"Results": rsp.Accounts,
		})

		return
	}
	render(w, r, "search", map[string]interface{}{})
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		vars := mux.Vars(r)
		id := vars["id"]
		if len(id) == 0 {
			return
		}
		_, err := accClient.Delete(context.TODO(), &account.DeleteRequest{
			Id: id,
		})
		if err != nil {
			setAlert(w, r, err.Error(), "error")
			http.Redirect(w, r, r.Referer(), 302)
			return
		}
		http.Redirect(w, r, r.Referer(), 302)
	}
}

func EditAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if r.Method == "POST" {
		r.ParseForm()

		typ := r.Form.Get("type")
		clientId := r.Form.Get("client_id")
		clientSecret := r.Form.Get("client_secret")

		var metadata map[string]string

		json.Unmarshal([]byte(r.Form.Get("metadata")), &metadata)

		_, err := accClient.Update(context.TODO(), &account.UpdateRequest{
			Account: &account.Record{
				Id: id,
				Type: typ,
				ClientId: clientId,
				ClientSecret: clientSecret,
				Metadata: metadata,
			},
		})

		if err != nil {
			setAlert(w, r, err.Error(), "error")
			http.Redirect(w, r, r.Referer(), 302)
			return
		}

		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	rsp, err := accClient.Read(context.TODO(), &account.ReadRequest{
		Id: id,
	})
	if err != nil {
		http.Redirect(w, r, r.Referer(), 302)
		return
	}

	render(w, r, "editAccount", map[string]interface{}{
		"Account": rsp.Account,
	})
}
/*

func Auth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}
	// TODO: limit/offset
	rsp, err := authClient.Read(context.TODO(), &auth.ReadRequest{
		Id: id,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	sort.Sort(sortedSpans{spans: rsp.Spans, reverse: true})

	for _, span := range rsp.Spans {
		if len(span.Annotations) == 0 {
			continue
		}
		sort.Sort(sortedAnns{span.Annotations})
	}

	render(w, r, "auth", map[string]interface{}{
		"Id":    id,
		"Spans": rsp.Spans,
	})
}
*/
