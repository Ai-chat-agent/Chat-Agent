package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/Ai-chat-agent/Chat-Agent.git/internal/config"
	"github.com/Ai-chat-agent/Chat-Agent.git/pkg/logger"
	"github.com/Ai-chat-agent/Chat-Agent.git/pkg/models"
)

// Database wraps the database connection
type Database struct {
	DB     *gorm.DB
	logger logger.Logger
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig, log logger.Logger) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	// Configure GORM logger
	gormLogger := gormlogger.Default
	if log != nil {
		gormLogger = gormlogger.Discard // Use our custom logger instead
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database := &Database{
		DB:     db,
		logger: log,
	}

	log.Info("Database connection established",
		logger.F("host", cfg.Host),
		logger.F("database", cfg.DBName),
	)

	return database, nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.ChatSession{},
		&models.ChatMessage{},
	)

	if err != nil {
		d.logger.Error("Database migration failed", logger.F("error", err.Error()))
		return fmt.Errorf("database migration failed: %w", err)
	}

	d.logger.Info("Database migration completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks the database connection
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Repository interfaces and implementations

// UserRepository handles user-related database operations
type UserRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("Failed to create user", logger.F("error", err.Error()))
		return err
	}
	r.logger.Info("User created", logger.F("user_id", user.ID))
	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get user", logger.F("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// ChatRepository handles chat-related database operations
type ChatRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewChatRepository creates a new chat repository
func NewChatRepository(db *gorm.DB, logger logger.Logger) *ChatRepository {
	return &ChatRepository{
		db:     db,
		logger: logger,
	}
}

// CreateMessage creates a new chat message
func (r *ChatRepository) CreateMessage(message *models.ChatMessage) error {
	if err := r.db.Create(message).Error; err != nil {
		r.logger.Error("Failed to create message", logger.F("error", err.Error()))
		return err
	}
	r.logger.Info("Message created", logger.F("message_id", message.ID))
	return nil
}

// GetMessagesByUserID retrieves messages for a specific user
func (r *ChatRepository) GetMessagesByUserID(userID string, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&messages).Error; err != nil {
		r.logger.Error("Failed to get messages", logger.F("error", err.Error()))
		return nil, err
	}

	return messages, nil
}

// DeleteMessage deletes a message by ID
func (r *ChatRepository) DeleteMessage(messageID string) error {
	if err := r.db.Where("id = ?", messageID).Delete(&models.ChatMessage{}).Error; err != nil {
		r.logger.Error("Failed to delete message", logger.F("error", err.Error()))
		return err
	}
	r.logger.Info("Message deleted", logger.F("message_id", messageID))
	return nil
}

