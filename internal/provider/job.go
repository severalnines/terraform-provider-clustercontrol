package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func NewCCJob(clusterOperation string) *openapi.Jobs {
	funcName := "NewCCJob"
	slog.Debug(funcName)

	jobs := openapi.NewJobs(clusterOperation)

	job := openapi.NewJobsJob()
	job.SetClassName(CMON_JOB_CLASS_NAME)

	jobSpec := openapi.NewJobsJobJobSpec()

	jobData := openapi.NewJobsJobJobSpecJobData()

	jobSpec.SetJobData(*jobData)
	job.SetJobSpec(*jobSpec)
	jobs.SetJob(*job)

	return jobs
}

func NewCCJobForJobStatusCheck(jobId int32) *openapi.Jobs {
	funcName := "NewCCJobForJobStatusCheck"
	slog.Debug(funcName)

	jobs := openapi.NewJobs(CMON_JOB_GET_JOB)
	jobs.SetJobId(jobId)
	return jobs
}

func SendAndWaitForJobCompletion(ctx context.Context, apiClient *openapi.APIClient, job *openapi.Jobs) error {
	funcName := "SendAndWaitForJobCompletion"
	slog.Debug(funcName)

	var resp *http.Response
	var err error

	//request, _ := json.Marshal(job)
	//slog.Info(string(request))
	//fmt.Fprintf(os.Stderr, string(request))
	// os.Stderr.Write(request)

	if resp, err = apiClient.JobsApi.JobsPost(ctx).Jobs(*job).Execute(); err != nil {
		PrintError(err, resp)
		return err
	}
	slog.Debug(funcName, "Resp `Job`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, resp)
		return err
	}

	var jobResp ResponseJobJson
	if err = json.Unmarshal(respBytes, &jobResp); err != nil {
		PrintError(err, resp)
		return err
	}
	//slog.Debug(funcName, "Job response", jobResp)

	// Wait for job to complete
	for true {
		// Calling Sleep method
		time.Sleep(10 * time.Second)

		checkJobStatus := NewCCJobForJobStatusCheck(jobResp.Job.Job_Id)
		if resp, err = apiClient.JobsApi.JobsPost(ctx).Jobs(*checkJobStatus).Execute(); err != nil {
			PrintError(err, resp)
			return err
		}

		if respBytes, err = io.ReadAll(resp.Body); err != nil {
			PrintError(err, resp)
			return err
		}

		if err = json.Unmarshal(respBytes, &jobResp); err != nil {
			PrintError(err, resp)
			return err
		}
		slog.Debug(funcName, "Job response", jobResp)

		if jobResp.Job.Status == JOB_STATUS_FINISHED {
			break
		}

		if jobResp.Job.Status != JOB_STATUS_RUNNING && jobResp.Job.Status != JOB_STATUS_DEFINED {
			err = errors.New(fmt.Sprintf("Job failed. (Status: %s)", jobResp.Job.Status))
			break
		}
	}

	return err
}
