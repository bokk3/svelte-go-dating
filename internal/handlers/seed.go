package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"dating-svelte/internal/models"
)

// SeedData creates test users, profiles, and matches for development
func SeedData(c *fiber.Ctx) error {
	// Check if we're in development mode
	if c.Query("confirm") != "yes" {
		return c.JSON(fiber.Map{
			"message": "Add ?confirm=yes to actually seed the database",
			"warning": "This will create test data in your database",
		})
	}

	// Create test users and profiles
	testUsers := []struct {
		Email       string
		DisplayName string
		Age         int
		Bio         string
		Location    string
		Gender      string
		Interests   []string
	}{
		{
			Email:       "alice.test@dating.com",
			DisplayName: "Alice",
			Age:         25,
			Bio:         "Love hiking, photography, and good coffee. Looking for someone who shares my passion for adventure!",
			Location:    "Chicago, IL",
			Gender:      "female",
			Interests:   []string{"hiking", "photography", "coffee", "travel"},
		},
		{
			Email:       "maria.test@dating.com",
			DisplayName: "Maria",
			Age:         22,
			Bio:         "Artist and dreamer. I spend my time painting, reading, and exploring new places.",
			Location:    "Miami, FL",
			Gender:      "female",
			Interests:   []string{"art", "reading", "museums", "yoga"},
		},
		{
			Email:       "jessica.test@dating.com",
			DisplayName: "Jessica",
			Age:         24,
			Bio:         "Music lover and concert-goer. Always up for trying new restaurants and live shows!",
			Location:    "Austin, TX",
			Gender:      "female",
			Interests:   []string{"music", "concerts", "food", "dancing"},
		},
		{
			Email:       "sophia.test@dating.com",
			DisplayName: "Sophia",
			Age:         26,
			Bio:         "Tech enthusiast and fitness lover. Building the future one line of code at a time.",
			Location:    "San Francisco, CA",
			Gender:      "female",
			Interests:   []string{"technology", "fitness", "coding", "startups"},
		},
		{
			Email:       "emma.test@dating.com",
			DisplayName: "Emma",
			Age:         23,
			Bio:         "Bookworm and cat lover. Perfect evening is a good book, tea, and cozy blanket.",
			Location:    "Portland, OR",
			Gender:      "female",
			Interests:   []string{"books", "cats", "tea", "writing"},
		},
	}

	createdUserIDs := make([]uuid.UUID, 0, len(testUsers))

	// Create users and profiles
	for _, userData := range testUsers {
		// Check if user already exists
		existingUser, _ := db.GetUserByEmail(userData.Email)
		if existingUser != nil {
			continue // Skip if user already exists
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
		if err != nil {
			continue
		}

		// Create user
		user := &models.User{
			ID:           uuid.New(),
			Email:        userData.Email,
			PasswordHash: string(hashedPassword),
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = db.CreateUser(user)
		if err != nil {
			continue
		}

		// Create profile
		age := userData.Age
		bio := userData.Bio
		gender := userData.Gender
		locationCity := userData.Location

		profile := &models.Profile{
			UserID:       user.ID,
			DisplayName:  userData.DisplayName,
			Age:          &age,
			Bio:          &bio,
			Gender:       &gender,
			LocationCity: &locationCity,
			InterestedIn: userData.Interests,
			Photos:       []models.Photo{}, // No photos for now
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = db.CreateProfile(profile)
		if err != nil {
			continue
		}

		createdUserIDs = append(createdUserIDs, user.ID)
	}

	// Get the current user (the one making the request)
	currentUserID := c.Locals("user_id").(uuid.UUID)

	// Create matches between current user and the test users
	matchesCreated := 0
	for _, testUserID := range createdUserIDs {
		// Create mutual swipes
		swipe1 := &models.Swipe{
			ID:        uuid.New(),
			SwiperID:  currentUserID,
			SwipedID:  testUserID,
			Liked:     true,
			CreatedAt: time.Now().Add(-time.Duration(matchesCreated+1) * time.Hour),
		}

		swipe2 := &models.Swipe{
			ID:        uuid.New(),
			SwiperID:  testUserID,
			SwipedID:  currentUserID,
			Liked:     true,
			CreatedAt: time.Now().Add(-time.Duration(matchesCreated+1) * time.Hour),
		}

		err := db.CreateSwipe(swipe1)
		if err != nil {
			continue
		}

		err = db.CreateSwipe(swipe2)
		if err != nil {
			continue
		}

		// Create match
		match := &models.Match{
			ID:        uuid.New(),
			User1ID:   currentUserID,
			User2ID:   testUserID,
			CreatedAt: time.Now().Add(-time.Duration(matchesCreated+1) * time.Hour),
		}

		err = db.CreateMatch(match)
		if err != nil {
			continue
		}

		// Create some sample messages
		messages := []string{
			"Hey there! I really liked your profile! 😊",
			"Thanks! Your interests look really cool too!",
			"Would you like to chat sometime?",
		}

		for i, messageText := range messages {
			senderID := testUserID
			if i%2 == 1 { // Alternate between users
				senderID = currentUserID
			}

			message := &models.Message{
				ID:        uuid.New(),
				MatchID:   match.ID,
				SenderID:  senderID,
				Message:   messageText,
				IsRead:    i < 2, // Mark first two as read
				CreatedAt: time.Now().Add(-time.Duration(matchesCreated+1)*time.Hour + time.Duration(i)*time.Minute*10),
			}

			db.CreateMessage(message)
		}

		matchesCreated++
	}

	return c.JSON(fiber.Map{
		"message":         "Seed data created successfully!",
		"users_created":   len(createdUserIDs),
		"matches_created": matchesCreated,
	})
}
