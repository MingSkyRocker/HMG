package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MingSkyRocker/HMG/helpers"
	"github.com/MingSkyRocker/HMG/models"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var db *mongo.Collection
var ctx = context.OCR()

// connect the database
func init() {
	rnd = renderer.New()

	// load env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("There is no env file")
	}

	// get env variable
	mongoUri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))
}

// Home Handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ocrs, err := helpers.FetchOCRsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err,
		})

		return
	}

	data := models.TemplateData{
		CSRFToken: nosurf.Token(r),
		Ocrs:      ocrs,
	}
	err = rnd.Template(w, http.StatusOK, []string{"views/index.html"}, data)
	if err != nil {
		log.Fatal(err)
	}
}

// Fetch all ocrs
func FetchOCRs(w http.ResponseWriter, r *http.Request) {
	ocrs, err := helpers.FetchOCRsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err,
		})

		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": ocrs,
	})
}

// Create ocr
func CreateOCR(w http.ResponseWriter, r *http.Request) {
	var t models.OCRModel
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title is required",
		})

		return
	}

	tm := models.OCRModel{
		ID:          t.ID,
		Title:       t.Title,
		CompanyName: t.CompanyName,
		ModelName:   t.ModelName,
	}

	_, err := db.InsertOne(ctx, &tm)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save",
			"error":   err,
		})
		return
	}

	ocrs, err := helpers.FetchOCRsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err,
		})

		return
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "OCR is created successfully",
		"ocrs":    ocrs,
	})
}

// Update OCR
func UpdateOCR(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var t models.OCRModel

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	// simple validation
	if t.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The title field is requried",
		})
		return
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed}}
	_, err = db.UpdateOne(ctx, filter, update)

	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to update OCR",
			"error":   err,
		})
		return
	}

	ocrs, err := helpers.FetchOCRsFormDB(db)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "OCR updated successfully",
		"ocrs":    ocrs,
	})
}

// Delete One OCR
func DeleteOneOCR(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	_, err = db.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete OCR",
			"error":   err,
		})
		return
	}

	ocrs, err1 := helpers.FetchOCRsFormDB(db)
	if err1 != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err1,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "OCR is successfully deleted",
		"ocr":     ocrs,
	})
}

// Delete completed OCR
func DeleteCompleted(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{
		"completed": bson.M{
			"$eq": true,
		},
	}
	_, err := db.DeleteMany(ctx, filter)
	if err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete completed OCR",
			"error":   err,
		})
		return
	}

	ocrs, err1 := helpers.FetchOCRsFormDB(db)
	if err1 != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch OCR",
			"error":   err1,
		})

		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "OCR is successfully deleted",
		"ocrs":    ocrs,
	})
}
