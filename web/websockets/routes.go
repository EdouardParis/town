package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	chirender "github.com/go-chi/render"
	"github.com/pkg/errors"
	melody "gopkg.in/olahol/melody.v1"

	"github.com/EdouardParis/town/app"
	"github.com/EdouardParis/town/failures"
	"github.com/EdouardParis/town/logging"
	"github.com/EdouardParis/town/managers"
	"github.com/EdouardParis/town/resources"
)

var cache = struct {
	counter  int
	sessions map[string]*melody.Session
	uuids    map[*melody.Session]string
	sync.RWMutex
}{
	sessions: make(map[string]*melody.Session),
	uuids:    make(map[*melody.Session]string),
}

func NewRouter(a *app.App) http.Handler {
	mrouter := melody.New()
	r := chi.NewRouter()
	r.Get("/checkout", handleError(a.Logger, mrouter.HandleRequest))

	mrouter.HandleConnect(func(s *melody.Session) {
		cache.Lock()
		defer cache.Unlock()

		resource, err := managers.CreateCharge(a.Config)
		if err != nil {
			s.Close()
			a.Logger.Error("Error during charge creation", logging.Error(err))
			return
		}

		rsc, err := json.Marshal(resource)
		if err != nil {
			s.Close()
			a.Logger.Error("Error during charge resource marshalling", logging.Error(err))
			return
		}

		err = s.Write(rsc)
		if err != nil {
			s.Close()
			a.Logger.Error("Error during charge resources writing", logging.Error(err))
			return
		}

		cache.sessions[resource.OrderID] = s
		cache.uuids[s] = resource.OrderID
		cache.counter += 1
		a.Logger.Info("New websocket connection", logging.Int("total_connected", cache.counter))
	})

	mrouter.HandleDisconnect(func(s *melody.Session) {
		cache.Lock()
		defer cache.Unlock()
		id := cache.uuids[s]
		delete(cache.uuids, s)
		delete(cache.sessions, id)
		cache.counter -= 1
		a.Logger.Info("websocket disconnected", logging.Int("total_connected", cache.counter))
	})
	return r
}

func SendChargeAndCloseSession(resource *resources.Charge) error {
	cache.Lock()
	defer cache.Unlock()

	session, ok := cache.sessions[resource.OrderID]
	if !ok {
		return fmt.Errorf("cannot find websocket connection with id %s", resource.OrderID)
	}

	rsc, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	err = session.Write(rsc)
	if err != nil {
		return err
	}

	session.Close()
	delete(cache.uuids, session)
	delete(cache.sessions, resource.OrderID)
	cache.counter -= 1

	return nil
}

func handleError(logger logging.Logger, fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err == nil {
			return
		}

		var status int
		switch cerr := errors.Cause(err).(type) {
		case failures.Error:
			status = cerr.Code
			err = cerr
		default:
			logger.Error(cerr.Error())
			status = http.StatusInternalServerError
		}

		chirender.Status(r, status)
		chirender.JSON(w, r, err)
	}
}
