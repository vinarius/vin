package watch

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/guregu/dynamo/v2"
	"github.com/vinarius/vin/constants"
	"github.com/vinarius/vin/models"
)

func getTopicsFromProjectIdFlag(projectIdsCsv string, dynamoClient dynamo.Table, ctx context.Context, stage string) []string {
	topics := []string{}
	projectIds := strings.Split(projectIdsCsv, ",")

	if len(projectIds) == 0 || (len(projectIds) == 1 && projectIds[0] == "") {
		return topics
	}

	keys := make([]dynamo.Keyed, 0, len(projectIds))
	for _, id := range projectIds {
		trimmedId := strings.TrimSpace(id)

		if trimmedId == "" {
			continue
		}

		pk := fmt.Sprintf("pId#%s", trimmedId)
		sk := fmt.Sprintf("pId#%s", trimmedId)
		keys = append(keys, models.PrimaryKey{Pk: pk, Sk: sk})
	}

	var projects []models.Project
	err := dynamoClient.Batch(constants.PK, constants.SK).Get(keys...).All(ctx, &projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to batch get projects from database: %s\n", err)
		os.Exit(1)
	}

	if len(projects) < len(keys) {
		fmt.Fprintf(os.Stderr, "Warning: Requested %d projects, but found %d. Some project IDs may have been invalid.\n", len(keys), len(projects))
	}

	for _, project := range projects {
		projectID := strings.TrimSpace(project.ProjectID)

		topic := fmt.Sprintf("iot-%s-ue2/%s", stage, projectID)
		topics = append(topics, topic)

		topic = fmt.Sprintf("iot-%s-ue2/%s/#", stage, projectID)
		topics = append(topics, topic)
	}

	return topics
}
