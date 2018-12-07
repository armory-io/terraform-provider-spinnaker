package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func GetProject(client *gate.GatewayClient, projectName string, dest interface{}) error {
	// todo: find what can be used to get projects
	project, resp, err := client.ProjectControllerApi.GetProjectUsingGET(client.Context, projectName, map[string]interface{}{})
	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Project '%s' not found\n", projectName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Encountered an error getting project, status code: %d\n", resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}

	if err := mapstructure.Decode(project, dest); err != nil {
		return err
	}

	return nil
}

func CreateProject(client *gate.GatewayClient, projectName, email string) error {

	project := map[string]interface{}{
		"name":  projectName,
		"email": email,
	}

	createProjectTask := map[string]interface{}{
		"job":         []interface{}{map[string]interface{}{"type": "upsertProject", "project": project}},
		"project":     projectName,
		"description": fmt.Sprintf("Create Project: %s", projectName),
	}

	ref, _, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, createProjectTask)
	if err != nil {
		return err
	}

	toks := strings.Split(ref["ref"].(string), "/")
	id := toks[len(toks)-1]

	task, resp, err := client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
	attempts := 0
	for (task == nil || !taskCompleted(task)) && attempts < 5 {
		toks := strings.Split(ref["ref"].(string), "/")
		id := toks[len(toks)-1]

		task, resp, err = client.TaskControllerApi.GetTaskUsingGET1(client.Context, id)
		attempts++
		time.Sleep(time.Duration(attempts*attempts) * time.Second)
	}

	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Encountered an error saving project, status code: %d\n", resp.StatusCode)
	}
	if !taskSucceeded(task) {
		return fmt.Errorf("Encountered an error saving project, task output was: %v\n", task)
	}

	return nil
}

func DeleteProject(client *gate.GatewayClient, projectName string) error {
	jobSpec := map[string]interface{}{
		"type": "deleteProject",
		"project": map[string]interface{}{
			"name": projectName,
		},
	}

	deleteProjectTask := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"project":     projectName,
		"description": fmt.Sprintf("Delete Project: %s", projectName),
	}

	_, resp, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, deleteProjectTask)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error deleting project, status code: %d\n", resp.StatusCode)
	}

	return nil
}

// todo: consider abstracting away to be shared with api/application
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

// todo: consider abstracting away to be shared with api/application
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
