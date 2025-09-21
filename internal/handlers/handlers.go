package handlers

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "github.com/google/uuid"
    
    "dating-svelte/internal/auth"
    "dating-svelte/internal/database"
    "dating-svelte/internal/models"
)

var db *database.DB

// InitializeHandlers sets up the database connection for handlers
func InitializeHandlers(database *database.DB) {
    db = database
}

// Auth handlers
type RegisterRequest struct {
    Email       string `json:"email"`
    Password    string `json:"password"`
    DisplayName string `json:"display_name"`
    Age         int    `json:"age"`
    Gender      string `json:"gender"`
    GDPRConsent bool   `json:"gdpr_consent"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    // Validation
    if req.Email == "" || req.Password == "" || req.DisplayName == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Email, password, and display name are required"})
    }
    
    if !req.GDPRConsent {
        return c.Status(400).JSON(fiber.Map{"error": "GDPR consent is required"})
    }
    
    if req.Age < 18 {
        return c.Status(400).JSON(fiber.Map{"error": "Must be 18 or older"})
    }
    
    // Check if user already exists
    _, err := db.GetUserByEmail(req.Email)
    if err == nil {
        return c.Status(409).JSON(fiber.Map{"error": "User already exists"})
    }
    
    // Hash password
    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to process password"})
    }
    
    // Create user
    userID := uuid.New()
    now := time.Now()
    user := &models.User{
        ID:            userID,
        Email:         req.Email,
        PasswordHash:  hashedPassword,
        Status:        "active",
        GDPRConsent:   req.GDPRConsent,
        GDPRConsentAt: &now,
        CreatedAt:     now,
        UpdatedAt:     now,
    }
    
    if err := db.CreateUser(user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
    }
    
    // Create profile
    profile := &models.Profile{
        UserID:      userID,
        DisplayName: req.DisplayName,
        Age:         &req.Age,
        Gender:      &req.Gender,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
    
    if err := db.CreateProfile(profile); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create profile"})
    }
    
    // Generate tokens
    tokens, err := auth.GenerateTokenPair(userID, req.Email, false)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
    }
    
    return c.Status(201).JSON(fiber.Map{
        "user":   user,
        "tokens": tokens,
    })
}

func Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    // Get user by email
    user, err := db.GetUserByEmail(req.Email)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
    }
    
    // Check password
    if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
    }
    
    // Update last active
    user.LastActive = time.Now()
    db.UpdateUser(user)
    
    // Get profile for premium status
    profile, _ := db.GetProfile(user.ID)
    isPremium := profile != nil && profile.IsPremium
    
    // Generate tokens
    tokens, err := auth.GenerateTokenPair(user.ID, user.Email, isPremium)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to generate tokens"})
    }
    
    return c.JSON(fiber.Map{
        "user":   user,
        "tokens": tokens,
    })
}

func GetCurrentUser(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    user, err := db.GetUser(userID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }
    
    return c.JSON(fiber.Map{"user": user})
}

func RefreshToken(c *fiber.Ctx) error {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    tokens, err := auth.RefreshAccessToken(req.RefreshToken)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
    }
    
    return c.JSON(fiber.Map{"tokens": tokens})
}

// Profile handlers
func GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    profile, err := db.GetProfile(userID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Profile not found"})
    }
    
    return c.JSON(profile)
}

func UpdateProfile(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    var profile models.Profile
    if err := c.BodyParser(&profile); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    profile.UserID = userID
    profile.UpdatedAt = time.Now()
    
    if err := db.UpdateProfile(&profile); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update profile"})
    }
    
    return c.JSON(profile)
}

// Match handlers
func GetMatches(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    // Get actual matches for this user
    matches, err := db.GetUserMatches(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get matches"})
    }
    
    // Populate each match with user profiles
    var enrichedMatches []fiber.Map
    for _, match := range matches {
        // Get profiles for both users
        user1Profile, _ := db.GetProfile(match.User1ID)
        user2Profile, _ := db.GetProfile(match.User2ID)
        
        // Determine which is the other user
        var otherUser *models.Profile
        if match.User1ID == userID {
            otherUser = user2Profile
        } else {
            otherUser = user1Profile
        }
        
        // Get last message for this match
        messages, _ := db.GetMatchMessages(match.ID)
        var lastMessage string
        var lastMessageTime *time.Time
        if len(messages) > 0 {
            lastMessage = messages[len(messages)-1].Message
            lastMessageTime = &messages[len(messages)-1].CreatedAt
        }
        
        enrichedMatch := fiber.Map{
            "id":               match.ID,
            "matched_at":       match.MatchedAt,
            "is_active":        match.IsActive,
            "other_user":       otherUser,
            "last_message":     lastMessage,
            "last_message_at":  lastMessageTime,
            "unread_count":     0, // TODO: Implement unread count
        }
        
        enrichedMatches = append(enrichedMatches, enrichedMatch)
    }
    
    return c.JSON(enrichedMatches)
}

func GetPotentialMatches(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    // Get potential matches for swiping
    profiles, err := db.GetPotentialMatches(userID, 10)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get potential matches"})
    }
    
    return c.JSON(profiles)
}

type SwipeRequest struct {
    TargetUserID uuid.UUID `json:"target_user_id"`
    Liked        bool      `json:"liked"`
}

func Swipe(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    var req SwipeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }
    
    // Create swipe record
    swipe := &models.Swipe{
        ID:        uuid.New(),
        SwiperID:  userID,
        SwipedID:  req.TargetUserID,
        Liked:     req.Liked,
        CreatedAt: time.Now(),
    }
    
    if err := db.CreateSwipe(swipe); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to record swipe"})
    }
    
    response := fiber.Map{"matched": false}
    
    // If liked, check for mutual match
    if req.Liked {
        isMatch, err := db.CheckForMatch(userID, req.TargetUserID)
        if err == nil && isMatch {
            // Create match record
            match := &models.Match{
                ID:        uuid.New(),
                User1ID:   userID,
                User2ID:   req.TargetUserID,
                MatchedAt: time.Now(),
                IsActive:  true,
                CreatedAt: time.Now(),
            }
            
            if err := db.CreateMatch(match); err == nil {
                response["matched"] = true
                response["match_id"] = match.ID
            }
        }
    }
    
    return c.JSON(response)
}

// WebSocket handler for real-time messaging
func WebSocketHandler(c *websocket.Conn) {
    defer c.Close()
    
    // TODO: Implement proper authentication for WebSocket
    // TODO: Handle real-time messaging, typing indicators, etc.
    
    for {
        messageType, msg, err := c.ReadMessage()
        if err != nil {
            break
        }
        
        // Echo back for now - implement proper message handling
        if err := c.WriteMessage(messageType, msg); err != nil {
            break
        }
    }
}

// Message handlers
func GetMessages(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    matchIDStr := c.Params("matchId")
    
    matchID, err := uuid.Parse(matchIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid match ID"})
    }
    
    // Verify user is part of this match
    matches, err := db.GetUserMatches(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get matches"})
    }
    
    validMatch := false
    for _, match := range matches {
        if match.ID == matchID {
            validMatch = true
            break
        }
    }
    
    if !validMatch {
        return c.Status(403).JSON(fiber.Map{"error": "Access denied to this match"})
    }
    
    // Get messages for this match
    messages, err := db.GetMatchMessages(matchID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get messages"})
    }
    
    return c.JSON(fiber.Map{
        "messages": messages,
        "match_id": matchID,
    })
}

func GetMatchDetails(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    matchIDStr := c.Params("matchId")
    
    matchID, err := uuid.Parse(matchIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid match ID"})
    }
    
    // Get match details
    matches, err := db.GetUserMatches(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get matches"})
    }
    
    var targetMatch *models.Match
    for _, match := range matches {
        if match.ID == matchID {
            targetMatch = &match
            break
        }
    }
    
    if targetMatch == nil {
        return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
    }
    
    // Get the other user's profile
    var otherUserID uuid.UUID
    if targetMatch.User1ID == userID {
        otherUserID = targetMatch.User2ID
    } else {
        otherUserID = targetMatch.User1ID
    }
    
    otherProfile, err := db.GetProfile(otherUserID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get other user's profile"})
    }
    
    return c.JSON(fiber.Map{
        "match": targetMatch,
        "other_user": otherProfile,
    })
}

// Payment handlers (placeholders)
func CreateSubscription(c *fiber.Ctx) error {
    // TODO: Implement Stripe integration
    return c.JSON(fiber.Map{"message": "Create subscription endpoint"})
}

func ProcessCryptoPayment(c *fiber.Ctx) error {
    // TODO: Implement cryptocurrency payment processing
    return c.JSON(fiber.Map{"message": "Process crypto payment endpoint"})
}

// GDPR handlers
func ExportData(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    // TODO: Implement data export functionality
    // Collect all user data: profile, photos, matches, messages, etc.
    
    return c.JSON(fiber.Map{
        "message": "Data export initiated",
        "user_id": userID,
    })
}

func DeleteAccount(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uuid.UUID)
    
    // TODO: Implement account deletion
    // Mark user as deleted, anonymize data, etc.
    
    return c.JSON(fiber.Map{
        "message": "Account deletion initiated",
        "user_id": userID,
    })
}