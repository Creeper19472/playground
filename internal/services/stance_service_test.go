package services

import (
	"testing"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupStanceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Issue{}, &models.Opinion{}, &models.Reference{}, &models.Stance{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	db.Create(&models.User{Username: "testuser", Email: "test@example.com"})
	db.Create(&models.User{Username: "user2", Email: "user2@example.com"})

	issueService := NewIssueService(db)
	issueService.CreateIssue(1, "issue", "Test Issue", "A test issue")

	opinionService := NewOpinionService(db)
	opinionService.CreateOpinion(1, 1, "Test opinion")

	refService := NewReferenceService(db)
	refService.CreateReference(1, 1, nil, "https://example.com", "Example", "A reference")

	return db
}

func TestStanceService_CreateStance_OnOpinion(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	opinionID := uint(1)
	stance, err := service.CreateStance(2, "doubt", "I don't think this is correct", &opinionID, nil)
	if err != nil {
		t.Fatalf("Failed to create stance: %v", err)
	}

	if stance.ID == 0 {
		t.Error("Expected stance ID to be set")
	}

	if stance.StanceType != "doubt" {
		t.Errorf("Expected stance type 'doubt', got '%s'", stance.StanceType)
	}

	if stance.Reason != "I don't think this is correct" {
		t.Errorf("Expected reason to be set correctly")
	}

	if stance.OpinionID == nil || *stance.OpinionID != 1 {
		t.Error("Expected opinion_id to be 1")
	}
}

func TestStanceService_CreateStance_OnReference(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	refID := uint(1)
	stance, err := service.CreateStance(2, "support", "This is a reliable source", nil, &refID)
	if err != nil {
		t.Fatalf("Failed to create stance: %v", err)
	}

	if stance.ReferenceID == nil || *stance.ReferenceID != 1 {
		t.Error("Expected reference_id to be 1")
	}
}

func TestStanceService_CreateStance_InvalidType(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	opinionID := uint(1)
	_, err := service.CreateStance(2, "invalid", "reason", &opinionID, nil)
	if err == nil {
		t.Error("Expected error for invalid stance type, got nil")
	}
}

func TestStanceService_CreateStance_NoTarget(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	_, err := service.CreateStance(2, "support", "reason", nil, nil)
	if err == nil {
		t.Error("Expected error when no target specified, got nil")
	}
}

func TestStanceService_ListStancesByOpinion(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	opinionID := uint(1)
	service.CreateStance(1, "support", "I agree", &opinionID, nil)
	service.CreateStance(2, "doubt", "I disagree", &opinionID, nil)

	stances, err := service.ListStancesByOpinion(1)
	if err != nil {
		t.Fatalf("Failed to list stances: %v", err)
	}

	if len(stances) != 2 {
		t.Errorf("Expected 2 stances, got %d", len(stances))
	}
}

func TestStanceService_ListStancesByReference(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	refID := uint(1)
	service.CreateStance(1, "support", "Reliable source", nil, &refID)

	stances, err := service.ListStancesByReference(1)
	if err != nil {
		t.Fatalf("Failed to list stances: %v", err)
	}

	if len(stances) != 1 {
		t.Errorf("Expected 1 stance, got %d", len(stances))
	}
}

func TestStanceService_UpdateStance(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	opinionID := uint(1)
	stance, _ := service.CreateStance(1, "support", "Original reason", &opinionID, nil)

	updated, err := service.UpdateStance(stance.ID, "oppose", "Changed my mind")
	if err != nil {
		t.Fatalf("Failed to update stance: %v", err)
	}

	if updated.StanceType != "oppose" {
		t.Errorf("Expected stance type 'oppose', got '%s'", updated.StanceType)
	}

	if updated.Reason != "Changed my mind" {
		t.Errorf("Expected updated reason")
	}
}

func TestStanceService_DeleteStance(t *testing.T) {
	db := setupStanceTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewStanceService(db)

	opinionID := uint(1)
	stance, _ := service.CreateStance(1, "support", "reason", &opinionID, nil)

	err := service.DeleteStance(stance.ID)
	if err != nil {
		t.Fatalf("Failed to delete stance: %v", err)
	}

	_, err = service.GetStanceByID(stance.ID)
	if err == nil {
		t.Error("Expected error getting deleted stance, got nil")
	}
}
