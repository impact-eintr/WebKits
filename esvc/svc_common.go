package esvc

import (
	"context"
	"os"
	"syscall"
)

func Run(svc Service, sig ...os.Signal) error {
	if err := svc.Init(); err != nil {
		return err
	}

	if err := svc.Start(); err != nil {
		return err
	}

	if len(sig) == 0 {
		sig = []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
		}
	}

	sigCh := make(chan os.Signal, 1)
	signalNotify(sigCh, sig...)

	var ctx context.Context
	if s, ok := svc.(Context); ok {
		ctx = s.Context()
	} else {
		ctx = context.Background()
	}

	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	return svc.Stop()

}
