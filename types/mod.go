package types

import (
	"errors"
	"time"
)

type UserServiceReq struct {
	Username  string `json:"username" binding:"required,max=20"`
	FirstName string `json:"first_name" binding:"omitempty,max=20"`
	LastName  string `json:"last_name" binding:"omitempty,max=20"`
	DOB       string `json:"dob" binding:"omitempty,max=20"`
	Gender    int    `json:"gender" binding:"omitempty"`
	Email     string `json:"email" binding:"omitempty,email,max=50"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
	Passwd    string `json:"passwd" binding:"omitempty,max=200"`
}
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type UserResp struct {
	UserName string `json:"user_name" form:"user_name" example:"FanOne"`
}

type UnitInforeq struct {
	CompanyName  string `json:"company_name" binding:"required,max=20"`
	BuildingName string `json:"building_name" binding:"required,max=20"`
}
type UnitInfoResp struct {
	UnitRentID             int    `json:"unit_rent_id" form:"unit_rent_id"`
	MonthlyRent            int    `json:"monthly_rent" form:"monthly_rent"`
	SquareFootage          int    `json:"square_footage" form:"square_footage"`
	AvailableDateForMoveIn string `json:"available_date_for_move_in" form:"available_date_for_move_in" binding:"required"`
	IsPetAllowed           bool   `json:"is_pet_allowed" form:"is_pet_allowed"`
}

type UpdatePets struct {
	CurrentPetName string `json:"current_pet_name" binding:"required,max=50"`
	CurrentPetType string `json:"current_pet_type" binding:"required,max=20"`
	NewPetName     string `json:"new_pet_name" binding:"required,max=50"`
	NewPetType     string `json:"new_pet_type" binding:"required,max=50"`
	NewPetSize     string `json:"new_pet_size" binding:"required,max=20"`
}
type GetPets struct {
	CurrentPetName string `json:"current_pet_name" binding:"required,max=50"`
	CurrentPetType string `json:"current_pet_type" binding:"required,max=20"`
	CurrentPetSize string `json:"current_pet_size" binding:"required,max=20"`
}
type InterestResp struct {
	Username    string    `json:"username"`
	UnitRentID  int       `json:"unit_rent_id"`
	RoommateCnt uint8     `json:"roommate_cnt"`
	MoveInDate  time.Time `json:"move_in_date"`
}
type UnitRentIDReq struct {
	UnitRentID int `json:"unit_rent_id" form:"unit_rent_id"`
}
type PostInterestReq struct {
	UnitRentID  int       `json:"unit_rent_id"`
	RoommateCnt uint8     `json:"roommate_cnt"`
	MoveInDate  time.Time `json:"move_in_date"`
}

type ComplexUnitinfo struct {
	CompanyName       string   `json:"company_name" binding:"required"`
	BuildingName      string   `json:"building_name" binding:"required"`
	Address           string   `json:"address"`
	YearBuilt         int      `json:"year_built" binding:"required"`
	Amenitiesinbuding []string `json:"amenities_building"`
	Amenitiesinunit   []string `json:"amenities_unit"`
	AvailableUnits    int      `json:"available_units"`

	UnitRentID             int       `json:"unit_rent_id"`
	MonthlyRent            int       `json:"monthly_rent" binding:"required"`
	SquareFootage          int       `json:"square_footage"`
	AvailableDateForMoveIn time.Time `json:"available_date_for_move_in"`

	BedroomNum    int `json:"bedroom_num"`
	BathroomNum   int `json:"bathroom_num"`
	LivingRoomNum int `json:"living_room_num"`
}

type UserProfile struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Gender    int    `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = str[1 : len(str)-1]
	parsedTime, err := time.Parse("2006-01-02", str)
	if err != nil {
		return errors.New("invalid date format")
	}
	t.Time = parsedTime
	return nil
}

func (t CustomTime) ToTime() time.Time {
	return t.Time
}

type InteresCondReq struct {
	UnitRentID  int        `json:"unit_rent_id"`
	RoommateCnt uint8      `json:"roommate_cnt" binding:"omitempty"`
	MoveInDate  CustomTime `json:"move_in_date" binding:"omitempty"`
}

type AveragePriceReq struct {
	AddrZipCode string `json:"addr_zip_code" binding:"required,max=5"`
	BedroomNum  int    `json:"bedroom_num"`
	BathroomNum int    `json:"bathroom_num"`
}
type AveragePriceResp struct {
	AverageRent float64 `json:"avg_rent"`
}
