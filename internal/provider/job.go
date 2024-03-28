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
	"strings"
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

	_, err := SendAndWaitForJobCompletionAndId(ctx, apiClient, job)

	return err
}

func SendAndWaitForJobCompletionAndId(ctx context.Context, apiClient *openapi.APIClient, job *openapi.Jobs) (int32, error) {
	funcName := "SendAndWaitForJobCompletionAndId"
	slog.Debug(funcName)

	var resp *http.Response
	var err error
	var jobResp JobResponseFields

	if resp, err = apiClient.JobsAPI.JobsPost(ctx).Jobs(*job).Execute(); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}
	slog.Debug(funcName, "Resp `Job`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}

	if err = json.Unmarshal(respBytes, &jobResp); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}
	//slog.Debug(funcName, "Job response", jobResp)

	// Wait for job to complete
	for true {
		// Calling Sleep method
		time.Sleep(10 * time.Second)

		checkJobStatus := NewCCJobForJobStatusCheck(jobResp.Job.Job_Id)
		if resp, err = apiClient.JobsAPI.JobsPost(ctx).Jobs(*checkJobStatus).Execute(); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}

		if respBytes, err = io.ReadAll(resp.Body); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}

		if err = json.Unmarshal(respBytes, &jobResp); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}
		slog.Debug(funcName, "Job response", jobResp)

		if strings.Contains(jobResp.Job.Status, JOB_STATUS_FINISHED) {
			break
		}

		//if jobResp.Job.Status != JOB_STATUS_RUNNING && jobResp.Job.Status != JOB_STATUS_DEFINED {
		// Sometimes the job status is RUNNING3. What the heck is this? But, it is indeed the REALITY. Go figure!
		if !strings.Contains(jobResp.Job.Status, JOB_STATUS_RUNNING) &&
			!strings.Contains(jobResp.Job.Status, JOB_STATUS_DEFINED) {
			err = errors.New(fmt.Sprintf("Job failed. (Status: %s)", jobResp.Job.Status))
			break
		}
	}

	return jobResp.Job.Job_Id, err
}

func CreateJobAndGetJobId(ctx context.Context, apiClient *openapi.APIClient, job *openapi.Jobs) (int32, error) {
	funcName := "CreateJobAndGetJobId"
	slog.Debug(funcName)

	var resp *http.Response
	var err error
	var jobResp JobResponseFields

	if resp, err = apiClient.JobsAPI.JobsPost(ctx).Jobs(*job).Execute(); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}
	slog.Debug(funcName, "Resp `Job`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}

	if err = json.Unmarshal(respBytes, &jobResp); err != nil {
		PrintError(err, resp)
		return jobResp.Job.Job_Id, err
	}
	//slog.Debug(funcName, "Job response", jobResp)

	// Wait for job to complete
	for true {
		// Calling Sleep method
		time.Sleep(10 * time.Second)

		checkJobStatus := NewCCJobForJobStatusCheck(jobResp.Job.Job_Id)
		if resp, err = apiClient.JobsAPI.JobsPost(ctx).Jobs(*checkJobStatus).Execute(); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}

		if respBytes, err = io.ReadAll(resp.Body); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}

		if err = json.Unmarshal(respBytes, &jobResp); err != nil {
			PrintError(err, resp)
			return jobResp.Job.Job_Id, err
		}
		slog.Debug(funcName, "Job response", jobResp)

		if strings.Contains(jobResp.Job.Status, JOB_STATUS_FINISHED) {
			break
		}

		//if jobResp.Job.Status != JOB_STATUS_RUNNING && jobResp.Job.Status != JOB_STATUS_DEFINED {
		// Sometimes the job status is RUNNING3. What the heck is this? But, it is indeed the REALITY. Go figure!
		if !strings.Contains(jobResp.Job.Status, JOB_STATUS_RUNNING) &&
			!strings.Contains(jobResp.Job.Status, JOB_STATUS_DEFINED) {
			err = errors.New(fmt.Sprintf("Job failed. (Status: %s)", jobResp.Job.Status))
			break
		}
	}

	return jobResp.Job.Job_Id, err
}
