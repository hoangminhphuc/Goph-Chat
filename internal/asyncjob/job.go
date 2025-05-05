package asyncjob

import "context"

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateSuccess
	StateFailed
	StateTimeout
)

// Fixed length
func (js JobState) String() string {
	return [...]string{"Init", "Running", "Success", "Failed", "Timeout"}[js]
}

type JobHandler func(ctx context.Context) error

type Job interface {
	Execute(ctx context.Context) error
}

type job struct {
	title 	string
	state 	JobState
	handler JobHandler
}

func NewJob(title string, handler JobHandler) *job {
	return &job{
		title: title,
		state: StateInit,
		handler: handler,
	}
}

func (j *job) Execute(ctx context.Context) error {
	j.state = StateRunning
	errChan := make(chan error, 1)

	go func() {
		errChan <- j.handler(ctx)
	}()

	select {
	case <-ctx.Done(): // context’s Done channel closes first
		j.state = StateTimeout
		return ctx.Err()
	case err := <-errChan: // handler finishes first
		if err != nil {
			j.state = StateFailed
			return err
		}
		j.state = StateSuccess
		return nil
	}
}