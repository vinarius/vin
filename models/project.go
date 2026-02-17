package models

import "time"

type Project struct {
	PK            string    `dynamo:"pk" json:"pk"`
	SK            string    `dynamo:"sk" json:"sk"`
	Approvers     Approvers `dynamo:"approvers" json:"approvers"`
	CameraCount   int       `dynamo:"cameraCount" json:"cameraCount"`
	ConfigVersion string    `dynamo:"configVersion" json:"configVersion"`
	CreatedAt     time.Time `dynamo:"createdAt" json:"createdAt"`
	Creator       string    `dynamo:"creator" json:"creator"`
	GSI1PK        string    `dynamo:"gsi1pk" json:"gsi1pk"`
	GSI1SK        string    `dynamo:"gsi1sk" json:"gsi1sk"`
	GSI2PK        string    `dynamo:"gsi2pk" json:"gsi2pk"`
	GSI2SK        string    `dynamo:"gsi2sk" json:"gsi2sk"`
	GSI3PK        string    `dynamo:"gsi3pk" json:"gsi3pk"`
	GSI3SK        string    `dynamo:"gsi3sk" json:"gsi3sk"`
	HWVersion     string    `dynamo:"hwVersion" json:"hwVersion"`
	LTEVersion    string    `dynamo:"lteVersion" json:"lteVersion"`
	MCUVersion    string    `dynamo:"mcuVersion" json:"mcuVersion"`
	ModelName     string    `dynamo:"modelName" json:"modelName"`
	OEM           string    `dynamo:"oem" json:"oem"`
	OSVersion     string    `dynamo:"osVersion" json:"osVersion"`
	Product       string    `dynamo:"product" json:"product"`
	Project       string    `dynamo:"project" json:"project"`
	ProjectID     string    `dynamo:"projectId" json:"projectId"`
	Region        string    `dynamo:"region" json:"region"`
	SMSCtrl       int       `dynamo:"smsCtrl" json:"smsCtrl"`
	Type          string    `dynamo:"type" json:"type"`
	UpdatedAt     time.Time `dynamo:"updatedAt" json:"updatedAt"`
	Version       string    `dynamo:"version" json:"version"`
}

type Approvers struct {
	Director []string `dynamo:"director" json:"director"`
	PM       []string `dynamo:"pm" json:"pm"`
	PO       []string `dynamo:"po" json:"po"`
}
