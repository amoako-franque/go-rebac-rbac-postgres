package main

import (
	"time"

	"gorm.io/datatypes"
)

type Permission struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"unique;not null" json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Role struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"unique;not null" json:"name"`
	Permissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RolePermission struct {
	RoleID       uint       `gorm:"primaryKey" json:"roleId"`
	PermissionID uint       `gorm:"primaryKey" json:"permissionId"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	CreatedAt    time.Time
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	Name      *string    `json:"name,omitempty"`
	Roles     []UserRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type UserRole struct {
	UserID    uint `gorm:"primaryKey" json:"userId"`
	RoleID    uint `gorm:"primaryKey" json:"roleId"`
	User      User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt time.Time
}

func (UserRole) TableName() string {
	return "user_roles"
}

type PatientRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PatientName string    `gorm:"not null" json:"patientName"`
	Data        string    `json:"data"`
	OwnerID     uint      `gorm:"not null" json:"ownerId"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Relationship struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	SubjectID uint           `gorm:"not null;index" json:"subjectId"`
	ObjectID  uint           `gorm:"not null;index" json:"objectId"`
	Type      string         `gorm:"not null;index" json:"type"`
	Meta      datatypes.JSON `json:"meta,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
