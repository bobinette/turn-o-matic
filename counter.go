package turnomatic

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Server struct {
	Upgrader websocket.Upgrader
	counter  map[string]chan int
}

func (s *Server) Display(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/display" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "display.html")
}

func (s *Server) HandleDisplay(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("Empty token"))
		return
	}

	ch, ok := s.counter[token]
	if !ok {
		conn.WriteMessage(websocket.CloseMessage, []byte("Invalid token"))
		return
	}

	go func() {
		for i := range ch {
			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(strconv.Itoa(i)))

			if err := w.Close(); err != nil {
				return
			}
		}
	}()
}

func (s *Server) Desk(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/desk" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "desk.html")
}

func (s *Server) HandleDesk(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("Empty token"))
		return
	}

	if _, ok := s.counter[token]; ok {
		conn.WriteMessage(websocket.CloseMessage, []byte("Token already in use"))
		return
	}

	if s.counter == nil {
		s.counter = make(map[string]chan int)
	}

	ch := make(chan int, 0)
	s.counter[token] = ch

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				break
			}

			if n, err := strconv.Atoi(string(msg)); err != nil {
				log.Println(err)
				w.Write([]byte(err.Error()))
			} else {
				log.Println(string(msg))
				w.Write(msg)
				ch <- n
			}

			if err := w.Close(); err != nil {
				log.Println(err)
			}
		}
	}()
}
