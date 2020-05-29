package redirect

import (
	"fmt"
	"html"
	"log"
	"net/http"

	jww "github.com/spf13/jwalterweatherman"
)

type Server struct {
	addr            string
	redirectHandler func(*Params) error
}

func Init(port uint, redirectHandler func(*Params) error) *Server {

	return &Server{
		addr:            fmt.Sprintf("localhost:%d", port),
		redirectHandler: redirectHandler}
}

func (s *Server) Serve() {
	http.HandleFunc("/callback", s.callbackFunc)

	log.Fatal(http.ListenAndServe(s.addr, nil))
}

func (s *Server) callbackFunc(w http.ResponseWriter, r *http.Request) {

	jww.DEBUG.Println("Received callback", *r.URL)
	redirect := DecodeRedirect(r.URL)
	err := s.redirectHandler(redirect)
	if err != nil {
		jww.ERROR.Println("Failed callback", err)
		fmt.Fprintf(w, "Error handling callback %q", html.EscapeString(r.URL.Path))
		return
	}
	fmt.Fprintf(w, "You are successfully authenticated %q", html.EscapeString(r.URL.Path))

}
