package asyncjob

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type jobGroup struct {
	jobs []Job
}

func NewJobGroup(jobs ...Job) *jobGroup {
	return &jobGroup{jobs: jobs}
}

func (jg *jobGroup) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx) // errgroup wraps waitgroup
	for _, jb := range jg.jobs {
		jb := jb
		g.Go(func() error {
			return jb.Execute(ctx)
		})
	}
	return g.Wait() // returns the first error encountered, or nil if all jobs succeeded
}

// func (jg *jobGroup) Run(ctx context.Context) error {
// 	jg.wg.Add(len(jg.jobs))

// 	errChan := make(chan error, len(jg.jobs))

// 	for _, jb := range jg.jobs {
// 		go func(j Job) {
// 			defer jg.wg.Done()
// 			err := j.Execute(ctx)
// 			if err != nil {
// 				errChan <- err
// 			}
// 		}(jb)
// 	}

// 	jg.wg.Wait()

// 	for i := 1; i <= len(jg.jobs); i++ {
// 			if v := <-errChan; v != nil {
// 					return v
// 			}
// 	}

// 	return nil
// }

