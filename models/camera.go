package models

import "time"

type (
	Camera struct {
		PK                 string           `dynamo:"pk,hash" json:"pk"`
		SK                 string           `dynamo:"sk,range" json:"sk"`
		CameraID           string           `dynamo:"cameraId" json:"cameraId"`
		CertLocation       string           `dynamo:"certLocation" json:"certLocation"`
		CreatedAt          time.Time        `dynamo:"createdAt" json:"createdAt"`
		CreationMethod     string           `dynamo:"creationMethod" json:"creationMethod"`
		FirmwareVersions   FirmwareVersions `dynamo:"firmwareVersions" json:"firmwareVersions"`
		GroupID            string           `dynamo:"groupId" json:"groupId"`
		GSI1PK             string           `dynamo:"gsi1pk" json:"gsi1pk"`
		GSI1SK             string           `dynamo:"gsi1sk" json:"gsi1sk"`
		GSI2PK             string           `dynamo:"gsi2pk" json:"gsi2pk"`
		GSI2SK             string           `dynamo:"gsi2sk" json:"gsi2sk"`
		GSI3PK             string           `dynamo:"gsi3pk" json:"gsi3pk"`
		GSI3SK             string           `dynamo:"gsi3sk" json:"gsi3sk"`
		GSI4PK             string           `dynamo:"gsi4pk" json:"gsi4pk"`
		GSI4SK             string           `dynamo:"gsi4sk" json:"gsi4sk"`
		IsOnline           int              `dynamo:"isOnline" json:"isOnline"`
		MigrationAudit     []any            `dynamo:"migrationAudit" json:"migrationAudit"`
		NetworkRegion      string           `dynamo:"networkRegion" json:"networkRegion"`
		PrivateKeyLocation string           `dynamo:"privateKeyLocation" json:"privateKeyLocation"`
		Product            string           `dynamo:"product" json:"product"`
		ProjectID          string           `dynamo:"projectId" json:"projectId"`
		PublicKeyLocation  string           `dynamo:"publicKeyLocation" json:"publicKeyLocation"`
		Type               string           `dynamo:"type" json:"type"`
		UpdatedAt          time.Time        `dynamo:"updatedAt" json:"updatedAt"`
		Vendor             string           `dynamo:"vendor" json:"vendor"`
	}
	FirmwareVersions struct {
		Config string `dynamo:"config" json:"config"`
		LTE    string `dynamo:"lte" json:"lte"`
		MCU    string `dynamo:"mcu" json:"mcu"`
		OS     string `dynamo:"os" json:"os"`
	}
)
