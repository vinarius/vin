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

func getTopicsFromCameraIdFlag(cameraIdsCsv string, dynamoClient dynamo.Table, ctx context.Context, stage string) []string {
	topics := []string{}
	cameraIds := strings.Split(cameraIdsCsv, ",")

	if len(cameraIds) == 0 || (len(cameraIds) == 1 && cameraIds[0] == "") {
		return topics
	}

	keys := make([]dynamo.Keyed, 0, len(cameraIds))
	for _, id := range cameraIds {
		trimmedId := strings.TrimSpace(id)

		if trimmedId == "" {
			continue
		}

		pk := fmt.Sprintf("cId#%s", trimmedId)
		sk := fmt.Sprintf("cId#%s", trimmedId)
		keys = append(keys, models.PrimaryKey{Pk: pk, Sk: sk})
	}

	var cameras []models.Camera
	err := dynamoClient.Batch(constants.PK, constants.SK).Get(keys...).All(ctx, &cameras)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to batch get cameras from database: %s\n", err)
		os.Exit(1)
	}

	if len(cameras) < len(keys) {
		fmt.Fprintf(os.Stderr, "Warning: Requested %d cameras, but found %d. Some camera IDs may have been invalid.\n", len(keys), len(cameras))
	}

	for _, camera := range cameras {
		projectID := strings.TrimSpace(camera.ProjectID)
		groupID := strings.TrimSpace(camera.GroupID)
		cameraID := strings.TrimSpace(camera.CameraID)

		topic := fmt.Sprintf("iot-%s-ue2/%s/%s/%s", stage, projectID, groupID, cameraID)
		topics = append(topics, topic)

		topic = fmt.Sprintf("iot-%s-ue2/%s/%s/%s/#", stage, projectID, groupID, cameraID)
		topics = append(topics, topic)
	}

	return topics
}
