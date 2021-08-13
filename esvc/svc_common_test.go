package esvc

import (
	"os"
	"syscall"
	"testing"
)

func TestDefaultSignalHandle(t *testing.T) {
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	for _, s := range sigs {
		testSignalNotify(t, s)
	}
}

func TestUserDefineSignalHandle(t *testing.T) {
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP}
	for _, s := range sigs {
		testSignalNotify(t, s, sigs...)
	}

}

func testSignalNotify(t *testing.T, signal os.Signal, sig ...os.Signal) {
	sigCh := make(chan os.Signal)

	var startCalled, stopCalled, initCalled int
	prg := makeProgram(&startCalled, &stopCalled, &initCalled)

	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		if c == nil {
			panic("os/signal: Notify using nil channel")
		}

		go func() {
			for val := range sigCh {
				for _, registeredSig := range sig {
					if val == registeredSig {
						c <- val
					}
				}
			}
		}()
	}

	go func() {
		sigCh <- signal // 发送信号
	}()

	if err := Run(prg, sig...); err != nil {
		t.Fatal(err)
	}

	// assert
	if startCalled != 1 {
		t.Errorf("startCalled[%d]", startCalled)
	}
	if stopCalled != 1 {
		t.Errorf("stopCalled[%d]", stopCalled)
	}
	if initCalled != 1 {
		t.Errorf("initCalled[%d]", initCalled)
	}

}
