package api

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func GetApplication(client *gate.GatewayClient, applicationName string, dest interface{}) error {
	app, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.ApplicationControllerApi.GetApplicationUsingGET(client.Context, applicationName, map[string]interface{}{})
	})

	if resp != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Application '%s' not found\n", applicationName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting application, status code: %d\n", resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}

	if err := mapstructure.Decode(app, dest); err != nil {
		return err
	}

	return nil
}

func CreateApplication(client *gate.GatewayClient, applicationData *schema.ResourceData) error {
	applicationName := applicationData.Get("application").(string)
	app := map[string]interface{}{
		"instancePort":   80,
		"name":           applicationName,
		"email":          applicationData.Get("email").(string),
		"repoType":       applicationData.Get("repo_type").(string),
		"repoProjectKey": applicationData.Get("repo_project_key").(string),
		"repoSlug":       applicationData.Get("repo_slug").(string),
	}

	createAppTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "createApplication", "application": app}},
		"application": applicationName,
		"description": fmt.Sprintf("Create Application: %s", applicationName),
	}

	ref, _, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.TaskControllerApi.TaskUsingPOST1(client.Context, createAppTask)
	})

	if err != nil {
		return err
	}

	toks := strings.Split(ref["ref"].(string), "/")
	id := toks[len(toks)-1]

	task, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
	})

	attempts := 0
	for (task == nil || !taskCompleted(task)) && attempts < 5 {
		toks := strings.Split(ref["ref"].(string), "/")
		id := toks[len(toks)-1]

		task, resp, err = retry(func() (map[string]interface{}, *http.Response, error) {
			return client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
		})
		attempts += 1
		time.Sleep(time.Duration(attempts*attempts) * time.Second)
	}

	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Encountered an error saving application, status code: %d\n", resp.StatusCode)
	}
	if !taskSucceeded(task) {
		return fmt.Errorf("Encountered an error saving application, task output was: %v\n", task)
	}

	return nil
}

func DeleteAppliation(client *gate.GatewayClient, applicationName string) error {
	jobSpec := map[string]interface{}{
		"type": "deleteApplication",
		"application": map[string]interface{}{
			"name": applicationName,
		},
	}

	deleteAppTask := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"application": applicationName,
		"description": fmt.Sprintf("Delete Application: %s", applicationName),
	}

	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.TaskControllerApi.TaskUsingPOST1(client.Context, deleteAppTask)
	})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting application, status code: %d\n", resp.StatusCode)
	}

	return nil
}

func taskCompleted(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	COMPLETED := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED", "TERMINAL", "FAILED_CONTINUE"}
	for _, status := range COMPLETED {
		if taskStatus == status {
			return true
		}
	}
	return false
}

func taskSucceeded(task map[string]interface{}) bool {
	taskStatus, exists := task["status"]
	if !exists {
		return false
	}

	SUCCESSFUL := [...]string{"SUCCEEDED", "STOPPED", "SKIPPED"}
	for _, status := range SUCCESSFUL {
		if taskStatus == status {
			return true
		}
	}
	return false
}
