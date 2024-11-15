package helpers

import (
	"context"
	"errors"

	"github.com/MingSkyRocker/HMG/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.OCR()

// Fetch OCR form database
func FetchOCRDB(db *mongo.Collection) ([]models.Todo, error) {
	// var todoM OCR
	var ocrs []models.OCRModel
	ocrList := []models.OCR{}
	cur, err := db.Find(ctx, bson.D{})
	if err != nil {
		defer cur.Close(ctx)
		return ocrList, errors.New("failed to fetch ocr")
	}

	if err = cur.All(ctx, &ocrs); err != nil {
		return ocrList, errors.New("failed to load data")
	}

	for _, t := range ocrs {
		ocrList = append(ocrList, models.OCR{
			ID:          t.ID.Hex(),
			Title:       t.Title,
			CompanyName: t.company_name,
			ModelName:   t.model,
		})
	}

	return ocrList, nil
}
