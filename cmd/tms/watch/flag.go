package watch

import (
	"fmt"

	"github.com/spf13/cobra"
)

type flag int

const (
	Topic flag = iota
	CameraId
	ProjectId
	Stage
)

func (f flag) String() string {
	switch f {
	case Topic:
		return "topic"
	case CameraId:
		return "camera-id"
	case ProjectId:
		return "project-id"
	case Stage:
		return "stage"
	default:
		return "Unknown watchflag value received"
	}
}

type GetFlagOutput struct {
	flag         flag
	stage, value string
}

func getFlag(cmd *cobra.Command) (*GetFlagOutput, error) {
	stage, err := cmd.Flags().GetString(Stage.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", Stage.String(), err)
		return nil, err
	}

	topic, err := cmd.Flags().GetString(Topic.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", Topic.String(), err)
		return nil, err
	}

	if topic != "" {
		return &GetFlagOutput{
			flag:  Topic,
			stage: stage,
			value: topic,
		}, nil
	}

	cameraId, err := cmd.Flags().GetString(CameraId.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", CameraId.String(), err)
		return nil, err
	}

	if cameraId != "" {
		return &GetFlagOutput{
			flag:  CameraId,
			stage: stage,
			value: cameraId,
		}, nil
	}

	projectId, err := cmd.Flags().GetString(ProjectId.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", ProjectId.String(), err)
		return nil, err
	}

	if projectId != "" {
		return &GetFlagOutput{
			flag:  ProjectId,
			stage: stage,
			value: projectId,
		}, nil
	}

	return nil, nil
}
