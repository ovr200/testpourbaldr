package main

import "time"

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Pseudo    string `gorm:"size:255"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string
	Grade     int
}

type Album struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Name        string `gorm:"size:255"`
	Description string `gorm:"size:255"`
	Image       string `gorm:"size:255"`
	Genre       string `gorm:"size:255"`
	Years       int    `gorm:"size:4"`
}

//Les favoris des users
type Favorite struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserId    uint
	Album     uint
}
