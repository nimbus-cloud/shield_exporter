package tasks

import (
	"github.com/starkandwayne/shield/api"
	"os"
	"log"
)

func init() {
	os.Setenv("SHIELD_SKIP_SSL_VERIFY", "true")
}

// backend should not have trailing slash, ie: https://10.92.246.91
func GetFailedBackups(backend, user, password string) ([]api.Task, error) {

	api.Cfg = &api.Config{
		Backend: "default",
		Backends: map[string]string{},
		Aliases:  map[string]string{},
	}

	err := api.Cfg.AddBackend(backend, "default")
	if err != nil {
		log.Printf("Error setting backend: %s\n", err)
		return nil, err
	}

	token := api.BasicAuthToken(user, password)
	err = api.Cfg.UpdateBackend(backend, token)
	if err != nil {
		log.Printf("Error updating backend: %s\n", err)
		return nil, err
	}

	failedTasks, err := api.GetTasks(api.TaskFilter{
		Limit:  "200",
		Status: "failed",	// possible values: 'pending', 'running', 'canceled', 'failed', 'done'
	})

	if err != nil {
		log.Printf("Error getting tasks: %s\n", err)
		return nil, err
	}

	// TODO: return all tasks (purge & backup) for the last 24 hours
	// TODO: lookup job details for tasks

	return failedTasks, nil
}
