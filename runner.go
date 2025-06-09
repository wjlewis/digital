package digital

type runner interface {
	steps() <-chan struct{}
	cmds() <-chan string
	respond(string)
	log(string)
	start()
	close()
}
