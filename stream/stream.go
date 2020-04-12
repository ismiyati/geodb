package stream

import (
	"context"
	api "github.com/autom8ter/geodb/gen/go/geodb"
	"github.com/gofrs/uuid"
	"sync"
)

var objectChan = make(chan *api.Object)
var eventChan = make(chan *api.Event)

type Hub struct {
	objectClients map[string]chan *api.Object
	eventClients  map[string]chan *api.Event
	mu            *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		objectClients: map[string]chan *api.Object{},
		eventClients:  map[string]chan *api.Event{},
		mu:            &sync.Mutex{},
	}
}

func (h *Hub) StartObjectStream(ctx context.Context) error {
	for {
		select {

		case obj := <-objectChan:
			if h.objectClients == nil {
				h.objectClients = map[string]chan *api.Object{}
			}

			for _, channel := range h.objectClients {
				if channel != nil {
					channel <- obj
				}
			}
		case <-ctx.Done():
			break
		}
	}
}

func (h *Hub) StartEventStream(ctx context.Context) error {
	for {
		select {

		case event := <-eventChan:
			if h.eventClients == nil {
				h.eventClients = map[string]chan *api.Event{}
			}

			for _, channel := range h.eventClients {
				if channel != nil {
					channel <- event
				}
			}
		case <-ctx.Done():
			break
		}
	}
}

func (h *Hub) AddObjectStreamClient(clientID string) string {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.objectClients == nil {
		h.objectClients = map[string]chan *api.Object{}
	}
	if clientID == "" {
		id, _ := uuid.NewV4()
		clientID = id.String()
	}
	h.objectClients[clientID] = make(chan *api.Object)
	return clientID
}

func (h *Hub) RemoveObjectStreamClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.objectClients[id]; ok {
		close(h.objectClients[id])
		delete(h.objectClients, id)
	}
}

func (h *Hub) GetClientObjectStream(id string) chan *api.Object {
	if _, ok := h.objectClients[id]; ok {
		return h.objectClients[id]
	}
	return nil
}

func (h *Hub) PublishObject(obj *api.Object) {
	objectChan <- obj
}

func (h *Hub) AddEventStreamClient(clientID string) string {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.eventClients == nil {
		h.eventClients = map[string]chan *api.Event{}
	}
	if clientID == "" {
		id, _ := uuid.NewV4()
		clientID = id.String()
	}
	h.eventClients[clientID] = make(chan *api.Event)
	return clientID
}

func (h *Hub) RemoveEventStreamClient(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.eventClients[id]; ok {
		close(h.eventClients[id])
		delete(h.eventClients, id)
	}
}

func (h *Hub) GetClientEventStream(id string) chan *api.Event {
	if _, ok := h.objectClients[id]; ok {
		return h.eventClients[id]
	}
	return nil
}

func (h *Hub) PublishEvent(event *api.Event) {
	eventChan <- event
}
