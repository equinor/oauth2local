package ipc

func HasSovereign() (bool, error) {
	resp := make(chan string)
	send("ping", "", resp)
	return "pong" == <-resp, nil

}

func SendRedirect(url string) {
	resp := make(chan string)
	send("redir", url, resp)
	<-resp
}
