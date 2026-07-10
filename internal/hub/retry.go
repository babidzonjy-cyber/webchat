package hub

import "time"

func retryOp(op func() error, attempts int, delay time.Duration) error {
	var err error

	for i := range attempts {
		if err = op(); err == nil {
			return nil
		}
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}
	return err
}
