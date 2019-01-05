package duci

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/logger"
	"github.com/pkg/errors"
	"time"
)

type duci struct {
	executor.Executor
	jobService jobService.Service
	github     github.GitHub
}

// New returns duci instance
func New() (executor.Executor, error) {
	jobService, err := jobService.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	github, err := github.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	builder, err := executor.DefaultExecutorBuilder()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	duci := &duci{
		jobService: jobService,
		github:     github,
	}
	duci.Executor = builder.
		StartFunc(duci.Start).
		EndFunc(duci.End).
		LogFunc(duci.AppendLog).
		Build()

	return duci, nil
}

// Start represents a function of start job
func (d *duci) Start(ctx context.Context) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	if err := d.jobService.Start(buildJob.ID); err != nil {
		if err := d.jobService.Append(buildJob.ID, job.LogLine{Timestamp: time.Now(), Message: err.Error()}); err != nil {
			logger.Error(err)
		}
		return
	}
	if err := d.github.CreateCommitStatus(ctx, github.CommitStatus{
		TargetSource: buildJob.TargetSource,
		State:        github.PENDING,
		Description:  "pending",
		Context:      buildJob.TaskName,
		TargetURL:    buildJob.TargetURL,
	}); err != nil {
		logger.Error(err)
	}
}

// AppendLog is a function that print and store log
func (d *duci) AppendLog(ctx context.Context, log job.Log) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
		println(line.Message)
		if err := d.jobService.Append(buildJob.ID, *line); err != nil {
			logger.Error(err)
		}
	}
}

// End represents a function
func (d *duci) End(ctx context.Context, e error) {
	buildJob, err := application.BuildJobFromContext(ctx)
	if err != nil {
		// TODO: output error message
		return
	}
	if err := d.jobService.Finish(buildJob.ID); err != nil {
		if err := d.jobService.Append(buildJob.ID, job.LogLine{Timestamp: time.Now(), Message: err.Error()}); err != nil {
			logger.Error(err)
		}
		return
	}

	switch e {
	case nil:
		if err := d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.SUCCESS,
			Description:  "success",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}); err != nil {
			logger.Error(err)
		}
	case runner.ErrFailure:
		if err := d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.FAILURE,
			Description:  "failure",
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}); err != nil {
			logger.Error(err)
		}
	default:
		if err := d.github.CreateCommitStatus(ctx, github.CommitStatus{
			TargetSource: buildJob.TargetSource,
			State:        github.ERROR,
			Description:  github.Description(fmt.Sprintf("error: %s", e.Error())),
			Context:      buildJob.TaskName,
			TargetURL:    buildJob.TargetURL,
		}); err != nil {
			logger.Error(err)
		}
	}
}