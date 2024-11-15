package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OCRModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	CompanyName string             `json:"company_name" bson:"company_name"`
	ModelName   string             `json:"model" bson:"model"`
}

type OCR struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	CompanyName string `json:"company_name" bson:"company_name"`
	ModelName   string `json:"model" bson:"model"`
}
