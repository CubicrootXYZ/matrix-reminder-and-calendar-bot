package database

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

// Service offers an interface for a matrix related database.
type Service interface {
	GetRoomByID(roomID string) (*MatrixRoom, error)
	NewRoom(room *MatrixRoom) (*MatrixRoom, error)
	UpdateRoom(room *MatrixRoom) (*MatrixRoom, error)

	GetUserByID(userID string) (*MatrixUser, error)
	NewUser(user *MatrixUser) (*MatrixUser, error)
}

// MatrixRoom holds information about a room.
type MatrixRoom struct {
	gorm.Model                   // numeric ID required to match main database in- and outputs
	RoomID          string       `gorm:"unique"`
	Users           []MatrixUser `gorm:"many2many:matrix_rooms_matrix_users;"`
	LastCryptoEvent string
}

// MatrixUser holds information about an user.
type MatrixUser struct {
	ID    string       `gorm:"primary"`
	Rooms []MatrixRoom `gorm:"many2many:matrix_rooms_matrix_users;"`
}
