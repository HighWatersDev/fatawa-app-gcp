package router

import (
	"crypto/rand"
	"encoding/base64"
	"fatawa-app-gcp/backend/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CreateAudioRequest struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	FilePath string `json:"filePath"`
	Duration int    `json:"duration"`
}

type CreateSegmentRequest struct {
	FullAudioID   string `json:"fullAudioId"`
	StartTime     int    `json:"startTime"`
	EndTime       int    `json:"endTime"`
	Transcription string `json:"transcription"`
}

type CreateQARequest struct {
	AudioSegmentID int64  `json:"audioSegmentId"`
	Question       string `json:"question"`
	Answer         string `json:"answer"`
}

// GenerateRandomSuffix generates a random string of a specified length.
func GenerateRandomSuffix(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func CreateAudio(c *gin.Context) {
	var req CreateAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	suffix, err := GenerateRandomSuffix(8)
	if err != nil {
		log.Printf("Error generating random suffix: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
		return
	}

	audioID := fmt.Sprintf("audio_%s", suffix)
	audio := &db.FullAudio{
		ID:         audioID,
		Title:      req.Title,
		Author:     req.Author,
		FilePath:   req.FilePath,
		Duration:   req.Duration,
		UploadTime: time.Now(),
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	if err := dbInstance.CreateFullAudio(audio); err != nil {
		log.Printf("Error creating audio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create audio"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Audio created successfully", "audioId": audioID})
}

func CreateSegment(c *gin.Context) {
	var req CreateSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	segment := &db.AudioSegment{
		FullAudioID:   req.FullAudioID,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Transcription: req.Transcription,
		Processed:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	if err := dbInstance.CreateAudioSegment(segment); err != nil {
		log.Printf("Error creating segment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create segment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Segment created successfully", "segmentId": segment.ID})
}

func GetAudio(c *gin.Context) {
	audioID := c.Param("id")
	if audioID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio ID is required"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	audio, err := dbInstance.GetFullAudioByID(audioID)
	if err != nil {
		log.Printf("Error retrieving audio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audio"})
		return
	}

	segments, err := dbInstance.GetAudioSegmentsByFullAudioID(audioID)
	if err != nil {
		log.Printf("Error retrieving segments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve segments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audio":    audio,
		"segments": segments,
	})
}

func DeleteAudio(c *gin.Context) {
	audioID := c.Param("id")
	if audioID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio ID is required"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	if err := dbInstance.DeleteFullAudio(audioID); err != nil {
		log.Printf("Error deleting audio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete audio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audio deleted successfully"})
}

func GetAudioSegments(c *gin.Context) {
	audioID := c.Param("id")
	if audioID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio ID is required"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	segments, err := dbInstance.GetAudioSegmentsByFullAudioID(audioID)
	if err != nil {
		log.Printf("Error retrieving segments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve segments"})
		return
	}

	c.JSON(http.StatusOK, segments)
}

func UpdateSegmentProcessed(c *gin.Context) {
	segmentID := c.Param("id")
	if segmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Segment ID is required"})
		return
	}

	var req struct {
		Processed bool `json:"processed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	segmentIDInt, err := strconv.ParseInt(segmentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}

	if err := dbInstance.UpdateAudioSegmentProcessedStatus(segmentIDInt, req.Processed); err != nil {
		log.Printf("Error updating segment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update segment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment status updated successfully"})
}

func CreateQA(c *gin.Context) {
	var req CreateQARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	post := &db.QAPost{
		AudioSegmentID: req.AudioSegmentID,
		Question:       req.Question,
		Answer:         req.Answer,
		PostTime:       time.Now(),
	}

	if err := dbInstance.CreateQAPost(post); err != nil {
		log.Printf("Error creating QA post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create QA post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "QA post created successfully", "postId": post.ID})
}

func GetSegmentQA(c *gin.Context) {
	segmentID := c.Param("id")
	if segmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Segment ID is required"})
		return
	}

	dbInstance, err := db.GetDB(c.Request.Context())
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}

	segmentIDInt, err := strconv.ParseInt(segmentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}

	posts, err := dbInstance.GetQAPostsByAudioSegmentID(segmentIDInt)
	if err != nil {
		log.Printf("Error retrieving QA posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve QA posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}