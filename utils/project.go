package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var defaultGCloudConfig = "$HOME/.config/gcloud/configurations/config_default"

// DetectDefaultProject searches for gcloud configuration to see what a
// default GCE project could be.
func DetectDefaultProject() string {
	path := os.ExpandEnv(defaultGCloudConfig)
	if _, err := os.Stat(path); err == nil {
		if raw, err := ioutil.ReadFile(path); err == nil {
			match := regexp.MustCompile(`project\s*=\s*([^\s]+)`).FindSubmatch(raw)
			if len(match) >= 1 {
				return string(match[1])
			}
		}
	}
	return ""
}

// DetermineProject will return the project part from the given topic resource
// or the first non-empty project.
//
// Also see: https://cloud.google.com/pubsub/docs/overview#names
func DetermineProject(resource string, projects ...string) (projectID, topicID string, err error) {
	resource = strings.Trim(resource, "/ ")
	parts := strings.Split(resource, "/")
	if len(parts) >= 4 && parts[0] == "projects" && parts[2] == "topics" {
		return parts[1], parts[3], nil
	} else if len(parts) >= 2 {
		return "", "", fmt.Errorf("invalid resource name: %s", resource)
	}
	topicID = resource
	for _, p := range projects {
		p = strings.TrimSpace(p)
		if p != "" {
			return p, topicID, nil
		}
	}
	defaultProject := DetectDefaultProject()
	if defaultProject == "" {
		return "", "", errors.New("unable to determine project")
	}
	return defaultProject, topicID, nil
}
