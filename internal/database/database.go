package database

import (
    "fmt"
    
    "github.com/jmoiron/sqlx"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
    
    "dating-svelte/internal/models"
)

type DB struct {
    *sqlx.DB
}

func New(dsn string) (*DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{db}, nil
}

// User methods
func (db *DB) GetUser(id uuid.UUID) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE id = $1 AND status != 'deleted'`
    err := db.Get(&user, query, id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    query := `SELECT * FROM users WHERE email = $1 AND status != 'deleted'`
    err := db.Get(&user, query, email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (db *DB) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (id, email, password_hash, status, gdpr_consent, gdpr_consent_at)
        VALUES (:id, :email, :password_hash, :status, :gdpr_consent, :gdpr_consent_at)
    `
    _, err := db.NamedExec(query, user)
    return err
}

func (db *DB) UpdateUser(user *models.User) error {
    query := `
        UPDATE users 
        SET email = :email, status = :status, last_active = :last_active, updated_at = NOW()
        WHERE id = :id
    `
    _, err := db.NamedExec(query, user)
    return err
}

// Profile methods
func (db *DB) GetProfile(userID uuid.UUID) (*models.Profile, error) {
    var profile models.Profile
    query := `SELECT * FROM profiles WHERE user_id = $1`
    err := db.Get(&profile, query, userID)
    if err != nil {
        return nil, err
    }
    
    // Get photos
    photos, err := db.GetUserPhotos(userID)
    if err == nil {
        profile.Photos = photos
    }
    
    return &profile, nil
}

func (db *DB) CreateProfile(profile *models.Profile) error {
    query := `
        INSERT INTO profiles (user_id, display_name, bio, age, gender, interested_in, 
                             location_city, location_country, latitude, longitude, avatar_url)
        VALUES (:user_id, :display_name, :bio, :age, :gender, :interested_in,
                :location_city, :location_country, :latitude, :longitude, :avatar_url)
    `
    _, err := db.NamedExec(query, profile)
    return err
}

func (db *DB) UpdateProfile(profile *models.Profile) error {
    query := `
        UPDATE profiles 
        SET display_name = :display_name, bio = :bio, age = :age, gender = :gender,
            interested_in = :interested_in, location_city = :location_city,
            location_country = :location_country, latitude = :latitude, 
            longitude = :longitude, avatar_url = :avatar_url, updated_at = NOW()
        WHERE user_id = :user_id
    `
    _, err := db.NamedExec(query, profile)
    return err
}

// Photo methods
func (db *DB) GetUserPhotos(userID uuid.UUID) ([]models.Photo, error) {
    var photos []models.Photo
    query := `SELECT * FROM photos WHERE user_id = $1 ORDER BY display_order`
    err := db.Select(&photos, query, userID)
    return photos, err
}

func (db *DB) AddPhoto(photo *models.Photo) error {
    query := `
        INSERT INTO photos (id, user_id, url, is_primary, display_order)
        VALUES (:id, :user_id, :url, :is_primary, :display_order)
    `
    _, err := db.NamedExec(query, photo)
    return err
}

// Swipe methods
func (db *DB) CreateSwipe(swipe *models.Swipe) error {
    query := `
        INSERT INTO swipes (id, swiper_id, swiped_id, liked)
        VALUES (:id, :swiper_id, :swiped_id, :liked)
        ON CONFLICT (swiper_id, swiped_id) DO UPDATE SET
        liked = :liked, created_at = NOW()
    `
    _, err := db.NamedExec(query, swipe)
    return err
}

func (db *DB) CheckForMatch(userID1, userID2 uuid.UUID) (bool, error) {
    var count int
    query := `
        SELECT COUNT(*) FROM swipes 
        WHERE ((swiper_id = $1 AND swiped_id = $2) OR (swiper_id = $2 AND swiped_id = $1))
        AND liked = true
    `
    err := db.Get(&count, query, userID1, userID2)
    return count == 2, err
}

func (db *DB) CreateMatch(match *models.Match) error {
    // Ensure user1_id < user2_id for consistency
    if match.User2ID.String() < match.User1ID.String() {
        match.User1ID, match.User2ID = match.User2ID, match.User1ID
    }
    
    query := `
        INSERT INTO matches (id, user1_id, user2_id, matched_at)
        VALUES (:id, :user1_id, :user2_id, :matched_at)
    `
    _, err := db.NamedExec(query, match)
    return err
}

// Match methods
func (db *DB) GetUserMatches(userID uuid.UUID) ([]models.Match, error) {
    var matches []models.Match
    query := `
        SELECT m.*, 
               p1.display_name as user1_display_name, p1.avatar_url as user1_avatar,
               p2.display_name as user2_display_name, p2.avatar_url as user2_avatar
        FROM matches m
        JOIN profiles p1 ON m.user1_id = p1.user_id
        JOIN profiles p2 ON m.user2_id = p2.user_id
        WHERE (m.user1_id = $1 OR m.user2_id = $1) AND m.is_active = true
        ORDER BY m.matched_at DESC
    `
    err := db.Select(&matches, query, userID)
    return matches, err
}

// Message methods
func (db *DB) GetMatchMessages(matchID uuid.UUID) ([]models.Message, error) {
    var messages []models.Message
    query := `
        SELECT * FROM messages 
        WHERE match_id = $1 
        ORDER BY created_at ASC
    `
    err := db.Select(&messages, query, matchID)
    return messages, err
}

func (db *DB) CreateMessage(message *models.Message) error {
    query := `
        INSERT INTO messages (id, match_id, sender_id, message, message_type)
        VALUES (:id, :match_id, :sender_id, :message, :message_type)
    `
    _, err := db.NamedExec(query, message)
    return err
}

// Discovery methods
func (db *DB) GetPotentialMatches(userID uuid.UUID, limit int) ([]models.Profile, error) {
    var profiles []models.Profile
    query := `
        SELECT p.* FROM profiles p
        JOIN users u ON p.user_id = u.id
        WHERE p.user_id != $1 
        AND u.status = 'active'
        AND p.user_id NOT IN (
            SELECT swiped_id FROM swipes WHERE swiper_id = $1
        )
        AND p.user_id NOT IN (
            SELECT CASE 
                WHEN user1_id = $1 THEN user2_id 
                ELSE user1_id 
            END FROM matches WHERE (user1_id = $1 OR user2_id = $1)
        )
        ORDER BY RANDOM()
        LIMIT $2
    `
    err := db.Select(&profiles, query, userID, limit)
    return profiles, err
}