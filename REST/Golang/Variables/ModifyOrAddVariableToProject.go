package main

import (
	"fmt"
	"log"

	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func main() {

	apiURL, err := url.Parse("https://YourURL")
	if err != nil {
		log.Println(err)
	}
	APIKey := "API-YourAPIKey"
	spaceName := "Default"
	variable := octopusdeploy.NewVariable("MyVariable")
	variable.IsSensitive = false
	variable.Type = "String"
	variable.Value = "MyValue"
	projectName := "MyProject"

	// Get reference to space
	space := GetSpace(apiURL, APIKey, spaceName)

	// Get project reference
	project := GetProject(apiURL, APIKey, space, projectName)

	// Get project variables
	projectVariables := GetProjectVariables(apiURL, APIKey, space, project)
	variableFound := false

	for i := 0; i < len(projectVariables.Variables); i++ {
		if projectVariables.Variables[i].Name == variable.Name {
			projectVariables.Variables[i].IsSensitive = variable.IsSensitive
			projectVariables.Variables[i].Type = variable.Type
			projectVariables.Variables[i].Value = variable.Value

			variableFound = true
			break
		}
	}

	if !variableFound {
		projectVariables.Variables = append(projectVariables.Variables, variable)
	}

	// Update target
	client := octopusAuth(apiURL, APIKey, space.ID)
	client.Variables.Update(project.ID, projectVariables)
}

func octopusAuth(octopusURL *url.URL, APIKey, space string) *octopusdeploy.Client {
	client, err := octopusdeploy.NewClient(nil, octopusURL, APIKey, space)
	if err != nil {
		log.Println(err)
	}

	return client
}

func GetSpace(octopusURL *url.URL, APIKey string, spaceName string) *octopusdeploy.Space {
	client := octopusAuth(octopusURL, APIKey, "")

	spaceQuery := octopusdeploy.SpacesQuery{
		Name: spaceName,
	}

	// Get specific space object
	spaces, err := client.Spaces.Get(spaceQuery)

	if err != nil {
		log.Println(err)
	}

	for _, space := range spaces.Items {
		if space.Name == spaceName {
			return space
		}
	}

	return nil
}

func GetProject(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, projectName string) *octopusdeploy.Project {
	// Create client
	client := octopusAuth(octopusURL, APIKey, space.ID)

	projectsQuery := octopusdeploy.ProjectsQuery {
		Name: projectName,
	}

	// Get specific project object
	projects, err := client.Projects.Get(projectsQuery)

	if err != nil {
		log.Println(err)
	}

	for _, project := range projects.Items {
		if project.Name == projectName {
			return project
		}
	}

	return nil
}

func GetProjectVariables(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, project *octopusdeploy.Project) octopusdeploy.VariableSet {
	// Create client
	client := octopusAuth(octopusURL, APIKey, space.ID)

	// Get project variables
	projectVariables, err := client.Variables.GetAll(project.ID)

	if err != nil {
		log.Println(err)
	}

	return projectVariables
}