package models

import "fmt"

type (
	Gx1 struct {
		Pk string `json:"-" dynamo:"gsi1pk"`
		Sk string `json:"-" dynamo:"gsi1sk"`
	}
	Gx2 struct {
		Pk string `json:"-" dynamo:"gsi2pk"`
		Sk string `json:"-" dynamo:"gsi2sk"`
	}
	Gx3 struct {
		Pk string `json:"-" dynamo:"gsi3pk"`
		Sk string `json:"-" dynamo:"gsi3sk"`
	}
	Gx4 struct {
		Pk string `json:"-" dynamo:"gsi4pk"`
		Sk string `json:"-" dynamo:"gsi4sk"`
	}
	Gx5 struct {
		Pk string `json:"-" dynamo:"gsi5pk"`
		Sk string `json:"-" dynamo:"gsi5sk"`
	}
	Gx6 struct {
		Pk string `json:"-" dynamo:"gsi6pk"`
		Sk string `json:"-" dynamo:"gsi6sk"`
	}
	Gx7 struct {
		Pk string `json:"-" dynamo:"gsi7pk"`
		Sk string `json:"-" dynamo:"gsi7sk"`
	}
)

func BuildGsi1Pk(entityType string) string {
	return fmt.Sprintf("t#%v", entityType)
}

func BuildGsi1Sk(createdAt string) string {
	return fmt.Sprintf("c#%v", createdAt)
}

func NewGx1(entityType, createdAt string) Gx1 {
	return Gx1{
		Pk: BuildGsi1Pk(entityType),
		Sk: BuildGsi1Sk(createdAt),
	}
}

func NewGx2(pk, sk string) Gx2 {
	return Gx2{
		Pk: pk,
		Sk: sk,
	}
}

func NewGx3(pk, sk string) Gx3 {
	return Gx3{
		Pk: pk,
		Sk: sk,
	}
}

func NewGx4(pk, sk string) Gx4 {
	return Gx4{
		Pk: pk,
		Sk: sk,
	}
}

func NewGx5(entityType, createdAt string) Gx5 {
	return Gx5{
		Pk: "int",
		Sk: fmt.Sprintf("t#%v#c#%s", entityType, createdAt),
	}
}

func NewGx6(pk, sk string) Gx6 {
	return Gx6{
		Pk: pk,
		Sk: sk,
	}
}

func NewGx7(pk, sk string) Gx7 {
	return Gx7{
		Pk: pk,
		Sk: sk,
	}
}
