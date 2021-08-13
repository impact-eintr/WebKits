package esvc

type testProgram struct {
	start func() error
	stop  func() error
	init  func() error
}

func (p *testProgram) Start() error {
	return p.start()
}

func (p *testProgram) Stop() error {
	return p.stop()
}

func (p *testProgram) Init() error {
	return p.init()
}

func makeProgram(startClassed, stopCalled, initCalled *int) *testProgram {
	return &testProgram{
		start: func() error {
			*startClassed++
			return nil
		},
		stop: func() error {
			*stopCalled++
			return nil
		},
		init: func() error {
			*initCalled++
			return nil
		},
	}
}
