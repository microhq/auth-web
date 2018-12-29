package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	proto "github.com/micro/auth-srv/proto/account"
)

const (
	alertId = "_ar"
)

type Alert struct {
	Type, Message string
}

var (
	store = sessions.NewCookieStore([]byte("auth"))
)

type sortedRecords struct {
	records   []*proto.Record
	reverse bool
}

func (s sortedRecords) Len() int {
	return len(s.records)
}

func (s sortedRecords) Less(i, j int) bool {
	if s.reverse {
		return s.records[i].Created > s.records[j].Created
	}
	return s.records[i].Created < s.records[j].Created
}

func (s sortedRecords) Swap(i, j int) {
	s.records[i], s.records[j] = s.records[j], s.records[i]
}

func distanceOfTime(minutes float64) string {
	switch {
	case minutes < 1:
		return fmt.Sprintf("%d secs", int(minutes*60))
	case minutes < 59:
		return fmt.Sprintf("%d minutes", int(minutes))
	case minutes < 90:
		return "about an hour"
	case minutes < 120:
		return "almost 2 hours"
	case minutes < 1080:
		return fmt.Sprintf("%d hours", int(minutes/60))
	case minutes < 1680:
		return "about a day"
	case minutes < 2160:
		return "more than a day"
	case minutes < 2520:
		return "almost 2 days"
	case minutes < 2880:
		return "about 2 days"
	default:
		return fmt.Sprintf("%d days", int(minutes/1440))
	}

	return ""
}

func getAlert(w http.ResponseWriter, r *http.Request) *Alert {
	sess, err := store.Get(r, alertId)
	if err != nil {
		return nil
	}
	defer sess.Save(r, w)

	for _, i := range []string{"info", "error", "success"} {
		f := sess.Flashes(i)
		if f != nil {
			if i == "error" {
				i = "danger"
			}

			return &Alert{
				Type:    i,
				Message: f[0].(string),
			}
		}
	}
	return nil
}

func setAlert(w http.ResponseWriter, r *http.Request, msg string, typ string) {
	sess, err := store.Get(r, alertId)
	if err != nil {
		return
	}
	sess.AddFlash(msg, typ)
	sess.Save(r, w)
}

func timeAgo(t int64) string {
	d := time.Unix(t, 0)
	timeAgo := ""
	startDate := time.Now().Unix()
	deltaMinutes := float64(startDate-d.Unix()) / 60.0
	if deltaMinutes <= 523440 { // less than 363 days
		timeAgo = fmt.Sprintf("%s ago", distanceOfTime(deltaMinutes))
	} else {
		timeAgo = d.Format("2 Jan")
	}

	return timeAgo
}

func hostPath(r *http.Request) string {
	if path := r.Header.Get("X-Micro-Web-Base-Path"); len(path) > 0 {
		return path
	}
	return "/"
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/accounts", Accounts)
	r.HandleFunc("/search", Search)
	r.HandleFunc("/delete/account/{id}", DeleteAccount)
	r.HandleFunc("/edit/account/{id}", EditAccount)
//	r.HandleFunc("/latest", Latest)
//	r.HandleFunc("/trace/{id}", Trace)
	return r
}
