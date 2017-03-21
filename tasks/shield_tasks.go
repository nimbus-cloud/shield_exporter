package tasks

import (
	"github.com/starkandwayne/shield/api"
	"os"
	"log"
	"time"
)

func init() {
	os.Setenv("SHIELD_SKIP_SSL_VERIFY", "true")
	api.Cfg = &api.Config{
		Backend: "default",
		Backends: map[string]string{},
		Aliases:  map[string]string{},
	}
}

// backend should not have trailing slash, ie: https://10.92.246.91
func GetFailedBackups(backend, user, password string) ([]api.Task, error) {

	if err := setupClient(backend,user, password); err != nil {
		log.Printf("Error setting up api client: %s\n", err)
		return nil, err
	}

	allTasks, err := api.GetTasks(api.TaskFilter{})
	if err != nil {
		log.Printf("Error getting tasks: %s\n", err)
		return nil, err
	}

	filteredTasks := filterTasks(allTasks, func(t api.Task) bool {
		dayAgo := time.Now().Add(-24 * time.Hour)
		return t.Op == "backup" && t.StartedAt.Time().After(dayAgo)
	})

	lookupTable, err := jobLookupTable()
	if err != nil {
		log.Printf("Error getting job lookup table: %s\n", err)
		return nil, err
	}

	for i := range filteredTasks {
		filteredTasks[i].JobUUID = lookupTable[filteredTasks[i].JobUUID]
	}

	return filteredTasks, nil
}

func jobLookupTable() (map[string]string, error)  {
	jobs, err := api.GetJobs(api.JobFilter{})
	if err != nil {
		log.Printf("Error getting jobs: %s\n", err)
		return nil, err
	}

	lookupTable := make(map[string]string)
	for _, job := range jobs {
		lookupTable[job.UUID] = job.Name
	}
	return lookupTable, nil
}

func setupClient(backend, user, password string) error {
	err := api.Cfg.AddBackend(backend, "default")
	if err != nil {
		log.Printf("Error setting backend: %s\n", err)
		return err
	}

	token := api.BasicAuthToken(user, password)
	err = api.Cfg.UpdateBackend(backend, token)
	if err != nil {
		log.Printf("Error updating backend: %s\n", err)
		return err
	}

	return nil
}

func filterTasks(vs []api.Task, f func(api.Task) bool) []api.Task {
	vsf := make([]api.Task, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}