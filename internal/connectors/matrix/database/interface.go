package database

import "gorm.io/gorm"

// Service offers an interface for a matrix related database.
type Service interface {
	GetRoomByID(roomID string) (*MatrixRoom, error)
	GetUserByID(userID string) (*MatrixUser, error)

	UpdateRoom(room *MatrixRoom) (*MatrixRoom, error)
}

// MatrixRoom holds information about a room.
type MatrixRoom struct {
	gorm.Model             // numeric ID required to match main database in- and outputs
	RoomID          string `gorm:"unique"`
	Users           []MatrixUser
	LastCryptoEvent string
}

// MatrixUser holds information about an user.
type MatrixUser struct {
	ID    string `gorm:"primary"`
	Rooms []MatrixRoom
}
