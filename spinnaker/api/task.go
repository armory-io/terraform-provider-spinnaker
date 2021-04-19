package api

import (
	"fmt"
	"net/http"

	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func EnableApplicationCanary(client *gate.GatewayClient, applicationName string) error {
	jobSpec := map[string]interface{}{
		"type": "updateApplication",
		"application": map[string]interface{}{
			"name": applicationName,
			"datasources": map[string]interface{}{
				"enabled":  "canaryConfigs",
				"disabled": "",
			},
			"user": "[anonymous]",
		},
	}
	enableAppCanaryTask := map[string]interface{}{
		"job":         []interface{}{jobSpec},
		"application": applicationName,
		"description": fmt.Sprintf("Enable canary for: %s", applicationName),
	}

	_, resp, err := client.TaskControllerApi.TaskUsingPOST1(client.Context, enableAppCanaryTask)

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Encountered an error enabling canary, status code: %d\n", resp.StatusCode)
	}
	return nil
}
