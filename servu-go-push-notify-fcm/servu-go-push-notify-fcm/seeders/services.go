package seeders

import (
	"github.com/baryogenesis2025/servu-go/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedServices(db *gorm.DB) error {

	services := []models.Service{
		*models.NewService("Auto Repair Service", "license", "Auto Repair Service"),
		*models.NewService("Home Cleaning Service", "license", "Home Cleaning Service"),
		*models.NewService("HVAC Maintenance", "license", "HVAC Maintenance"),
		*models.NewService("Lawn Care & Gardening", "license", "Lawn Care & Gardening"),
		*models.NewService("Moving Assistance", "license", "Moving Assistance"),
		*models.NewService("Electrical Services", "license", "Electrical Services"),
		*models.NewService("Plumbing Services", "license", "Plumbing Services"),
		*models.NewService("Pest Control", "license", "Pest Control"),
		*models.NewService("Locksmith Services", "license", "Locksmith Services"),
		*models.NewService("Advocate Services", "license", "Advocate Services"),
		*models.NewService("Dog Care", "license", "Dog Care"),
		*models.NewService("Excavation Services", "license", "Excavation Services"),
		*models.NewService("Nanny", "license", "Nanny"),
		*models.NewService("Photographer", "license", "Photographer"),
		*models.NewService("Baby Sitter", "license", "Baby Sitter"),
		*models.NewService("Carpenter", "license", "Carpenter"),
		*models.NewService("Towing Services", "license", "Towing Services"),
		*models.NewService("Laundry Services", "license", "Laundry Services"),
		*models.NewService("Driving Services", "license", "Driving Services"),
	}

	return db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "title"},
			},
			DoNothing: true,
		}).
		Create(&services).Error
}
