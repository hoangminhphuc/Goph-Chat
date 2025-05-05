package asyncjob

import (
	"context"
	"sync"
)

type jobGroup struct {
	jobs []Job
	wg   *sync.WaitGroup
}

func NewJobGroup(jobs ...Job) *jobGroup {
	return &jobGroup{
		jobs: jobs,
		wg:   new(sync.WaitGroup),
	}
}

func (jg *jobGroup) Run(ctx context.Context) error {
	jg.wg.Add(len(jg.jobs))

	errChan := make(chan error, len(jg.jobs))

	for _, jb := range jg.jobs {
		go func(j Job) {
			defer jg.wg.Done()
			err := j.Execute(ctx)
			if err != nil {
				errChan <- err
			}
		}(jb)
	}

	jg.wg.Wait()

	for i := 1; i <= len(jg.jobs); i++ {
			if v := <-errChan; v != nil {
					return v
			}
	}

	return nil
}