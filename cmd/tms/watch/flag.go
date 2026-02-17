package watch

import (
	"fmt"

	"github.com/spf13/cobra"
)

type flag int

const (
	Topic flag = iota
	CameraId
	GroupId
	ProjectId
	Stage
	All
)

func (f flag) String() string {
	switch f {
	case Topic:
		return "topic"
	case CameraId:
		return "camera-id"
	case GroupId:
		return "group-id"
	case ProjectId:
		return "project-id"
	case Stage:
		return "stage"
	case All:
		return "all"
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

	groupId, err := cmd.Flags().GetString(GroupId.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", GroupId.String(), err)
		return nil, err

	}

	if groupId != "" {
		return &GetFlagOutput{
			flag:  GroupId,
			stage: stage,
			value: groupId,
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

	all, err := cmd.Flags().GetBool(All.String())
	if err != nil {
		fmt.Printf("error getting %s flag: %s\n", All.String(), err)
		return nil, err
	}

	if all {
		return &GetFlagOutput{
			flag:  All,
			stage: stage,
			value: "",
		}, nil
	}

	return nil, nil
}

func flagToPtr(f flag) *flag {
	return &f
}
