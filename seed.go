package main

import (
	"log"
	"time"
)

func Seed() map[string]any {
	// Clear existing data (optional - comment out if you want to preserve data)
	DB.Exec("TRUNCATE TABLE relationships, patient_records, user_roles, role_permissions, users, roles, permissions RESTART IDENTITY CASCADE")

	// Create permissions
	permRead := Permission{Name: "record:read"}
	permWrite := Permission{Name: "record:write"}
	if err := DB.Create(&permRead).Error; err != nil {
		log.Printf("Error creating permRead: %v", err)
	}
	if err := DB.Create(&permWrite).Error; err != nil {
		log.Printf("Error creating permWrite: %v", err)
	}

	// Create roles
	doctor := Role{Name: "doctor"}
	nurse := Role{Name: "nurse"}
	admin := Role{Name: "admin"}
	DB.Create(&doctor)
	DB.Create(&nurse)
	DB.Create(&admin)

	// Assign permissions to roles
	DB.Create(&RolePermission{RoleID: doctor.ID, PermissionID: permRead.ID})
	DB.Create(&RolePermission{RoleID: doctor.ID, PermissionID: permWrite.ID})
	DB.Create(&RolePermission{RoleID: nurse.ID, PermissionID: permRead.ID})
	DB.Create(&RolePermission{RoleID: admin.ID, PermissionID: permRead.ID})
	DB.Create(&RolePermission{RoleID: admin.ID, PermissionID: permWrite.ID})

	// Create users
	hashedPassword, _ := HashPassword("password")
	doc := User{
		Email:     "doc@example.com",
		Password:  hashedPassword,
		Name:      ptrString("Dr Alice"),
		CreatedAt: time.Now(),
	}
	nurseU := User{
		Email:     "nurse@example.com",
		Password:  hashedPassword,
		Name:      ptrString("Nora Nurse"),
		CreatedAt: time.Now(),
	}
	patient := User{
		Email:     "patient@example.com",
		Password:  hashedPassword,
		Name:      ptrString("Patient Paul"),
		CreatedAt: time.Now(),
	}
	outsider := User{
		Email:     "outsider@example.com",
		Password:  hashedPassword,
		Name:      ptrString("Outsider"),
		CreatedAt: time.Now(),
	}

	DB.Create(&doc)
	DB.Create(&nurseU)
	DB.Create(&patient)
	DB.Create(&outsider)

	// Assign roles to users
	DB.Create(&UserRole{UserID: doc.ID, RoleID: doctor.ID})
	DB.Create(&UserRole{UserID: nurseU.ID, RoleID: nurse.ID})

	// Create patient records
	rec1 := PatientRecord{
		PatientName: "Patient Paul",
		Data:        "Medical record A - Blood pressure: 120/80",
		OwnerID:     patient.ID,
	}
	rec2 := PatientRecord{
		PatientName: "Patient Paul",
		Data:        "Medical record B - Lab results: Normal",
		OwnerID:     patient.ID,
	}
	DB.Create(&rec1)
	DB.Create(&rec2)

	// Create relationships
	// Doctor is assigned to patient (indirect relationship)
	DB.Create(&Relationship{
		SubjectID: doc.ID,
		ObjectID:  patient.ID,
		Type:      "assigned_to",
	})

	// Also create a direct relationship: doctor -> record (for ReBAC example)
	DB.Create(&Relationship{
		SubjectID: doc.ID,
		ObjectID:  rec1.ID,
		Type:      "assigned_to",
	})

	out := map[string]any{
		"users": map[string]uint{
			"doctor":  doc.ID,
			"nurse":   nurseU.ID,
			"patient": patient.ID,
			"outsider": outsider.ID,
		},
		"records": []uint{rec1.ID, rec2.ID},
		"message": "Database seeded successfully. Use 'password' for all users.",
	}

	log.Println("Seed completed:", out)
	return out
}

func ptrString(s string) *string {
	return &s
}
