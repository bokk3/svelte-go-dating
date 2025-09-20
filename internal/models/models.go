package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	EmailVerified *time.Time `json:"email_verified_at" db:"email_verified_at"`
	Status        string     `json:"status" db:"status"`
	GDPRConsent   bool       `json:"gdpr_consent" db:"gdpr_consent"`
	GDPRConsentAt *time.Time `json:"gdpr_consent_at" db:"gdpr_consent_at"`
	LastActive    time.Time  `json:"last_active" db:"last_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	Profile       *Profile   `json:"profile,omitempty"`
}

type Profile struct {
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	DisplayName     string         `json:"display_name" db:"display_name"`
	Bio             *string        `json:"bio" db:"bio"`
	Age             *int           `json:"age" db:"age"`
	Gender          *string        `json:"gender" db:"gender"`
	InterestedIn    pq.StringArray `json:"interested_in" db:"interested_in"`
	LocationCity    *string        `json:"location_city" db:"location_city"`
	LocationCountry *string        `json:"location_country" db:"location_country"`
	Latitude        *float64       `json:"latitude" db:"latitude"`
	Longitude       *float64       `json:"longitude" db:"longitude"`
	AvatarURL       *string        `json:"avatar_url" db:"avatar_url"`
	IsVerified      bool           `json:"is_verified" db:"is_verified"`
	IsPremium       bool           `json:"is_premium" db:"is_premium"`
	Photos          []Photo        `json:"photos,omitempty"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

type Photo struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	URL          string    `json:"url" db:"url"`
	IsPrimary    bool      `json:"is_primary" db:"is_primary"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Swipe struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SwiperID  uuid.UUID `json:"swiper_id" db:"swiper_id"`
	SwipedID  uuid.UUID `json:"swiped_id" db:"swiped_id"`
	Liked     bool      `json:"liked" db:"liked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Match struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	MatchedAt time.Time `json:"matched_at" db:"matched_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	User1     *Profile  `json:"user1,omitempty"`
	User2     *Profile  `json:"user2,omitempty"`
}

type Message struct {
	ID          uuid.UUID `json:"id" db:"id"`
	MatchID     uuid.UUID `json:"match_id" db:"match_id"`
	SenderID    uuid.UUID `json:"sender_id" db:"sender_id"`
	Message     string    `json:"message" db:"message"`
	MessageType string    `json:"message_type" db:"message_type"`
	IsRead      bool      `json:"is_read" db:"is_read"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Subscription struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	UserID               uuid.UUID  `json:"user_id" db:"user_id"`
	PlanType             string     `json:"plan_type" db:"plan_type"`
	Status               string     `json:"status" db:"status"`
	StripeSubscriptionID *string    `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	CryptoTxHash         *string    `json:"crypto_tx_hash" db:"crypto_tx_hash"`
	StartsAt             time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt               *time.Time `json:"ends_at" db:"ends_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
}

// Location helper struct
type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
