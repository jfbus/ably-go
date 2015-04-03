package realtime

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ably/ably-go/config"
	"github.com/ably/ably-go/protocol"
	"github.com/ably/ably-go/rest"

	"github.com/ably/ably-go/Godeps/_workspace/src/code.google.com/p/go.net/websocket"
)

type ConnState int

const (
	ConnStateInitialized ConnState = iota
	ConnStateConnecting
	ConnStateConnected
	ConnStateDisconnected
	ConnStateSuspended
	ConnStateClosed
	ConnStateFailed
)

type ConnListener func()

type Conn struct {
	config.Params

	ID        string
	State     ConnState
	stateChan chan ConnState
	Ch        chan *protocol.ProtocolMessage
	Err       chan error
	ws        *websocket.Conn
	mtx       sync.RWMutex
	stateMtx  sync.RWMutex

	listeners   map[ConnState][]ConnListener
	listenerMtx sync.RWMutex
}

func NewConn(params config.Params) *Conn {
	c := &Conn{
		Params:    params,
		State:     ConnStateInitialized,
		stateChan: make(chan ConnState),
		Ch:        make(chan *protocol.ProtocolMessage),
		Err:       make(chan error),
	}
	go c.watchConnectionState()
	return c
}

func (c *Conn) isActive() bool {
	// Are you sure a lock is necessary in this function ?
	c.stateMtx.RLock()
	defer c.stateMtx.RUnlock()

	switch c.State {
	case ConnStateConnecting,
		ConnStateConnected:
		return true
	default:
		return false
	}
}

func (c *Conn) send(msg *protocol.ProtocolMessage) error {
	// I would add a read lock on c.mtx and check the connection state
	return websocket.JSON.Send(c.ws, msg)
}

func (c *Conn) Close() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.setState(ConnStateFailed)
	c.setConnectionID("")
	c.ws.Close()
}

func (c *Conn) websocketUrl(token *rest.Token) (*url.URL, error) {
	u, err := url.Parse(c.Params.RealtimeEndpoint + "?access_token=" + token.ID + "&binary=false&timestamp=" + strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c *Conn) Connect() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.isActive() {
		return nil
	}

	c.setState(ConnStateConnecting)

	restClient := rest.NewClient(c.Params)
	token, err := restClient.Auth.RequestToken(60*60, &rest.Capability{"*": []string{"*"}})
	if err != nil {
		return fmt.Errorf("Error fetching token: %s", err)
	}

	u, err := c.websocketUrl(token)
	if err != nil {
		return err
	}

	return c.dial(u)
}

func (c *Conn) dial(u *url.URL) error {
	ws, err := websocket.Dial(u.String(), "", "https://"+u.Host)
	if err != nil {
		return err
	}

	c.Ch = make(chan *protocol.ProtocolMessage)
	c.Err = make(chan error)
	c.ws = ws

	go c.read()

	return nil
}

func (c *Conn) read() {
	for {
		msg := &protocol.ProtocolMessage{}
		err := websocket.JSON.Receive(c.ws, &msg)
		if err != nil {
			c.Close()
			c.Err <- fmt.Errorf("Failed to get websocket frame: %s", err)
			return
		}
		c.handle(msg)
	}
}

func (c *Conn) watchConnectionState() {
	for connState := range c.stateChan {
		c.trigger(connState)
	}
}

func (c *Conn) trigger(connState ConnState) {
	for i := range c.listeners[connState] {
		go c.listeners[connState][i]()
	}
}

func (c *Conn) On(connState ConnState, listener ConnListener) {
	c.listenerMtx.Lock()
	defer c.listenerMtx.Unlock()

	if c.listeners == nil {
		c.listeners = make(map[ConnState][]ConnListener)
	}

	if c.listeners[connState] == nil {
		c.listeners[connState] = []ConnListener{}
	}

	c.listeners[connState] = append(c.listeners[connState], listener)
}

func (c *Conn) setState(state ConnState) {
	// We could change the semantics of watchConnectionState to manageConnectionState
	// and move c.State = state to manageConnectionState
	// this could lead to the removal of stateMtx.
	c.stateMtx.Lock()
	defer c.stateMtx.Unlock()

	c.State = state
	c.stateChan <- state
}

func (c *Conn) setConnectionID(id string) {
	// This could be merged with setState to remove one Lock
	// or move to manageConnectionState and try to completely
	// remove the mutex
	c.stateMtx.Lock()
	defer c.stateMtx.Unlock()

	c.ID = id
}

func (c *Conn) handle(msg *protocol.ProtocolMessage) {
	switch msg.Action {
	case protocol.ActionHeartbeat, protocol.ActionAck, protocol.ActionNack:
		return
	case protocol.ActionError:
		if msg.Channel != "" {
			c.Ch <- msg
			return
		}
		c.Close()
		c.Err <- msg.Error
		return
	case protocol.ActionConnected:
		c.setState(ConnStateConnected)
		c.setConnectionID(msg.ConnectionId)
		return
	case protocol.ActionDisconnected:
		c.setState(ConnStateDisconnected)
		c.setConnectionID("")
		return
	default:
		c.Ch <- msg
	}
}
