package digital

type testRunner struct {
	stepsChan     chan struct{}
	cmdsChan      chan string
	responsesData []string
	logData       []string
}

func newTestRunner() *testRunner {
	return &testRunner{
		stepsChan: make(chan struct{}),
		cmdsChan:  make(chan string),
		logData:   make([]string, 0),
	}
}

func (r *testRunner) steps() <-chan struct{} {
	return r.stepsChan
}

func (r *testRunner) cmds() <-chan string {
	return r.cmdsChan
}

func (r *testRunner) respond(msg string) {
	r.responsesData = append(r.responsesData, msg)
}

func (r *testRunner) log(msg string) {
	r.logData = append(r.logData, msg)
}

func (r *testRunner) start() {}

func (r *testRunner) close() {
	close(r.stepsChan)
	close(r.cmdsChan)
}

func (r *testRunner) step() {
	r.stepsChan <- struct{}{}
}

func (r *testRunner) sendCmd(cmd string) {
	r.cmdsChan <- cmd
}
