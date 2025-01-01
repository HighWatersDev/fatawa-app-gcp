package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FullAudio struct {
	ID         string
	Title      string
	Author     string
	FilePath   string
	Duration   int
	UploadTime time.Time
}

type AudioSegment struct {
	ID            int64
	FullAudioID   string
	StartTime     int
	EndTime       int
	Transcription string
	Processed     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type QAPost struct {
	ID             int64
	AudioSegmentID int64
	Question       string
	Answer         string
	PostTime       time.Time
}

type AudioDatabase struct {
	pool *pgxpool.Pool
}

var (
	instance *AudioDatabase
	once     sync.Once
)

func GetDB(ctx context.Context) (*AudioDatabase, error) {
	var err error
	once.Do(func() {
		connString := os.Getenv("DB_CONN_STRING")
		if connString == "" {
			err = fmt.Errorf("DB_CONN_STRING environment variable is not set")
			return
		}
		instance, err = NewAudioDatabase(connString)
	})
	return instance, err
}

// NewAudioDatabase creates a new database connection pool
func NewAudioDatabase(connString string) (*AudioDatabase, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &AudioDatabase{pool: pool}, nil
}

// CreateFullAudio inserts a new full audio record
func (db *AudioDatabase) CreateFullAudio(fa *FullAudio) error {
	ctx := context.Background()
	query := `INSERT INTO full_audios 
		(id, title, author, file_path, duration) 
		VALUES ($1, $2, $3, $4, $5)`
	
	_, err := db.pool.Exec(ctx, query, 
		fa.ID, 
		fa.Title, 
		fa.Author, 
		fa.FilePath, 
		fa.Duration,
	)

	return err
}

// CreateAudioSegment inserts a new audio segment
func (db *AudioDatabase) CreateAudioSegment(segment *AudioSegment) error {
	ctx := context.Background()
	query := `INSERT INTO audio_segments 
		(full_audio_id, start_time, end_time, transcription, processed, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	now := time.Now()
	err := db.pool.QueryRow(ctx, query, 
		segment.FullAudioID, 
		segment.StartTime, 
		segment.EndTime, 
		segment.Transcription, 
		segment.Processed,
		now,
		now,
	).Scan(&segment.ID)

	return err
}

// CreateQAPost inserts a new Q&A post
func (db *AudioDatabase) CreateQAPost(post *QAPost) error {
	ctx := context.Background()
	query := `INSERT INTO qa_posts 
		(audio_segment_id, question, answer) 
		VALUES ($1, $2, $3) RETURNING id`
	
	err := db.pool.QueryRow(ctx, query, 
		post.AudioSegmentID, 
		post.Question, 
		post.Answer,
	).Scan(&post.ID)

	return err
}

// GetFullAudioByID retrieves a full audio record by ID
func (db *AudioDatabase) GetFullAudioByID(id string) (*FullAudio, error) {
	ctx := context.Background()
	query := `SELECT id, title, author, file_path, duration, upload_time 
			  FROM full_audios WHERE id = $1`
	
	fa := &FullAudio{}
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&fa.ID, 
		&fa.Title, 
		&fa.Author, 
		&fa.FilePath, 
		&fa.Duration, 
		&fa.UploadTime,
	)

	if err != nil {
		return nil, err
	}

	return fa, nil
}

// GetAudioSegmentsByFullAudioID retrieves segments for a specific full audio
func (db *AudioDatabase) GetAudioSegmentsByFullAudioID(fullAudioID string) ([]AudioSegment, error) {
	ctx := context.Background()
	query := `SELECT id, full_audio_id, start_time, end_time, transcription, processed, created_at, updated_at 
			  FROM audio_segments WHERE full_audio_id = $1`
	
	rows, err := db.pool.Query(ctx, query, fullAudioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []AudioSegment
	for rows.Next() {
		var segment AudioSegment
		err := rows.Scan(
			&segment.ID, 
			&segment.FullAudioID, 
			&segment.StartTime, 
			&segment.EndTime, 
			&segment.Transcription, 
			&segment.Processed, 
			&segment.CreatedAt, 
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

// UpdateAudioSegmentProcessedStatus updates the processed status of a segment
func (db *AudioDatabase) UpdateAudioSegmentProcessedStatus(segmentID int64, processed bool) error {
	ctx := context.Background()
	query := `UPDATE audio_segments 
		SET processed = $2, updated_at = $3 
		WHERE id = $1`
	
	_, err := db.pool.Exec(ctx, query, 
		segmentID, 
		processed, 
		time.Now(),
	)

	return err
}

// GetQAPostsByAudioSegmentID retrieves Q&A posts for a specific audio segment
func (db *AudioDatabase) GetQAPostsByAudioSegmentID(audioSegmentID int64) ([]QAPost, error) {
	ctx := context.Background()
	query := `SELECT id, audio_segment_id, question, answer, post_time 
			  FROM qa_posts WHERE audio_segment_id = $1`
	
	rows, err := db.pool.Query(ctx, query, audioSegmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []QAPost
	for rows.Next() {
		var post QAPost
		err := rows.Scan(
			&post.ID, 
			&post.AudioSegmentID, 
			&post.Question, 
			&post.Answer, 
			&post.PostTime,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// DeleteFullAudio removes a full audio record by ID
func (db *AudioDatabase) DeleteFullAudio(id string) error {
	ctx := context.Background()
	query := `DELETE FROM full_audios WHERE id = $1`
	
	_, err := db.pool.Exec(ctx, query, id)
	return err
}

// Close terminates the database connection pool
func (db *AudioDatabase) Close() {
	db.pool.Close()
}
