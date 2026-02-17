package transport

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketConn struct {
	conn websocket.Conn
}

func DialWebsocket(url string) (*WebsocketConn, error) {
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &WebsocketConn{conn: *wsConn}, nil
}

func (w *WebsocketConn) Read(b []byte) (n int, err error) {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	copy(b, message)
	return len(message), nil
}

func (w *WebsocketConn) Write(b []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *WebsocketConn) Close() error {
	return w.conn.Close()
}

func (w *WebsocketConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *WebsocketConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *WebsocketConn) SetDeadline(t time.Time) error {
	return nil
}

func (w *WebsocketConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *WebsocketConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}
