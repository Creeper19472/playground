package services

import (
	"testing"

	"github.com/Creeper19472/playground/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOpinionTestDB(t *testing.T) *gorm.DB {
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

	return db
}

func TestOpinionService_CreateOpinion(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	opinion, err := service.CreateOpinion(1, 1, "This is my opinion")
	if err != nil {
		t.Fatalf("Failed to create opinion: %v", err)
	}

	if opinion.ID == 0 {
		t.Error("Expected opinion ID to be set")
	}

	if opinion.Content != "This is my opinion" {
		t.Errorf("Expected content 'This is my opinion', got '%s'", opinion.Content)
	}

	if opinion.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", opinion.Username)
	}
}

func TestOpinionService_ListOpinionsByIssue(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	service.CreateOpinion(1, 1, "Opinion 1")
	service.CreateOpinion(2, 1, "Opinion 2")

	opinions, err := service.ListOpinionsByIssue(1)
	if err != nil {
		t.Fatalf("Failed to list opinions: %v", err)
	}

	if len(opinions) != 2 {
		t.Errorf("Expected 2 opinions, got %d", len(opinions))
	}
}

func TestOpinionService_SupportRelationship(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	op1, _ := service.CreateOpinion(1, 1, "Opinion A supports Opinion B")
	op2, _ := service.CreateOpinion(2, 1, "Opinion B is the main point")

	// op1 supports op2
	err := service.AddSupport(op1.ID, op2.ID)
	if err != nil {
		t.Fatalf("Failed to add support: %v", err)
	}

	// Check what op1 supports
	supports, err := service.GetSupports(op1.ID)
	if err != nil {
		t.Fatalf("Failed to get supports: %v", err)
	}
	if len(supports) != 1 || supports[0].ID != op2.ID {
		t.Errorf("Expected op1 to support op2")
	}

	// Check what supports op2
	supportedBy, err := service.GetSupportedBy(op2.ID)
	if err != nil {
		t.Fatalf("Failed to get supported by: %v", err)
	}
	if len(supportedBy) != 1 || supportedBy[0].ID != op1.ID {
		t.Errorf("Expected op2 to be supported by op1")
	}

	// Remove the support
	err = service.RemoveSupport(op1.ID, op2.ID)
	if err != nil {
		t.Fatalf("Failed to remove support: %v", err)
	}

	supports, _ = service.GetSupports(op1.ID)
	if len(supports) != 0 {
		t.Errorf("Expected no supports after removal, got %d", len(supports))
	}
}

func TestOpinionService_MultipleSupports(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	opMain, _ := service.CreateOpinion(1, 1, "Main opinion")
	opSup1, _ := service.CreateOpinion(1, 1, "Supporting opinion 1")
	opSup2, _ := service.CreateOpinion(2, 1, "Supporting opinion 2")

	// Both support the main opinion
	service.AddSupport(opSup1.ID, opMain.ID)
	service.AddSupport(opSup2.ID, opMain.ID)

	supportedBy, err := service.GetSupportedBy(opMain.ID)
	if err != nil {
		t.Fatalf("Failed to get supported by: %v", err)
	}
	if len(supportedBy) != 2 {
		t.Errorf("Expected 2 supporters, got %d", len(supportedBy))
	}

	// One opinion supports multiple opinions
	opMain2, _ := service.CreateOpinion(1, 1, "Another main opinion")
	service.AddSupport(opSup1.ID, opMain2.ID)

	supports, err := service.GetSupports(opSup1.ID)
	if err != nil {
		t.Fatalf("Failed to get supports: %v", err)
	}
	if len(supports) != 2 {
		t.Errorf("Expected opSup1 to support 2 opinions, got %d", len(supports))
	}
}

func TestOpinionService_UpdateOpinion(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	opinion, _ := service.CreateOpinion(1, 1, "Original content")
	updated, err := service.UpdateOpinion(opinion.ID, "Updated content")
	if err != nil {
		t.Fatalf("Failed to update opinion: %v", err)
	}

	if updated.Content != "Updated content" {
		t.Errorf("Expected 'Updated content', got '%s'", updated.Content)
	}
}

func TestOpinionService_DeleteOpinion(t *testing.T) {
	db := setupOpinionTestDB(t)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	service := NewOpinionService(db)

	opinion, _ := service.CreateOpinion(1, 1, "To be deleted")
	err := service.DeleteOpinion(opinion.ID)
	if err != nil {
		t.Fatalf("Failed to delete opinion: %v", err)
	}

	_, err = service.GetOpinionByID(opinion.ID)
	if err == nil {
		t.Error("Expected error getting deleted opinion, got nil")
	}
}
