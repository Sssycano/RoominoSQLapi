package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Username  string
	FirstName string
	LastName  string
	DOB       string
	Gender    int
	Email     string
	Phone     string
	Passwd    string
	Pets      []Pets
	Interests []Interests
}

type Pets struct {
	PetName  string
	PetType  string
	PetSize  string
	Username string
}
type ApartmentBuilding struct {
	CompanyName   string
	BuildingName  string
	AddrNum       int
	AddrStreet    string
	AddrCity      string
	AddrState     string
	AddrZipCode   string
	YearBuilt     int
	ApartmentUnit []ApartmentUnit
	PetPolicies   []PetPolicy
	Provides      []Provides
}

type ApartmentUnit struct {
	UnitRentID             int
	CompanyName            string
	BuildingName           string
	UnitNumber             string
	MonthlyRent            int
	SquareFootage          int
	AvailableDateForMoveIn string
	Rooms                  []Rooms
	AmenitiesIn            []AmenitiesIn
	Interests              []Interests
}

type Interests struct {
	Username    string
	UnitRentID  int
	RoommateCnt uint8
	MoveInDate  time.Time
}
type Rooms struct {
	Name          string
	SquareFootage int
	Description   string
	UnitRentID    uint
}
type PetPolicy struct {
	CompanyName     string
	BuildingName    string
	PetType         string
	PetSize         string
	IsAllowed       bool
	RegistrationFee int
	MonthlyFee      int
}

type Amenities struct {
	AType       string
	Description string
	AmenitiesIn []AmenitiesIn
	Provides    []Provides
}

type AmenitiesIn struct {
	AType      string
	UnitRentID uint
}

type Provides struct {
	AType        string
	CompanyName  string
	BuildingName string
	Fee          int
	WaitingList  int
}

const (
	PassWordCost = 12
)

func (user *Users) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Passwd = string(bytes)
	return nil
}

func (user *Users) CheckPassword(password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(password))
	if err != nil {
		fmt.Println("Password comparison failed:", err)
	}

	return err == nil
}
