package dao

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"roomino/model"
	"roomino/types"
	"strings"
	"time"
)

type TaskDao struct {
	DB *sql.DB
}

func NewTaskDao(ctx context.Context) *TaskDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &TaskDao{NewDBClient(ctx)}
}

func (dao *TaskDao) GetUnitsWithPetPolicy(companyName, buildingName, username string) ([]types.UnitInfoResp, error) {
	queryUnits := "SELECT unit_rent_id, monthly_rent, square_footage, available_date_for_move_in FROM ApartmentUnit WHERE company_name = ? AND building_name = ?"
	unitRows, err := dao.DB.Query(queryUnits, companyName, buildingName)
	if err != nil {
		return nil, err
	}
	defer unitRows.Close()
	var units []model.ApartmentUnit
	for unitRows.Next() {
		var unit model.ApartmentUnit
		if err := unitRows.Scan(&unit.UnitRentID, &unit.MonthlyRent, &unit.SquareFootage, &unit.AvailableDateForMoveIn); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	queryPets := "SELECT pet_name, pet_type, pet_size FROM pets WHERE username = ?"
	petRows, err := dao.DB.Query(queryPets, username)
	if err != nil {
		return nil, err
	}
	defer petRows.Close()

	var userPets []model.Pets
	for petRows.Next() {
		var pet model.Pets
		if err := petRows.Scan(&pet.PetName, &pet.PetType, &pet.PetSize); err != nil {
			return nil, err
		}
		userPets = append(userPets, pet)
	}
	queryPetPolicy := "SELECT pet_type, pet_size, is_allowed FROM PetPolicy WHERE company_name = ? AND building_name = ?"
	petPolicyRows, err := dao.DB.Query(queryPetPolicy, companyName, buildingName)
	if err != nil {
		return nil, err
	}
	defer petPolicyRows.Close()

	var petPolicies []model.PetPolicy
	for petPolicyRows.Next() {
		var policy model.PetPolicy
		if err := petPolicyRows.Scan(&policy.PetType, &policy.PetSize, &policy.IsAllowed); err != nil {
			return nil, err
		}
		petPolicies = append(petPolicies, policy)
	}
	petPolicyMap := make(map[string]bool)
	for _, policy := range petPolicies {
		key := policy.PetType + "-" + policy.PetSize
		petPolicyMap[key] = policy.IsAllowed
	}
	var unitInfos []types.UnitInfoResp
	for _, unit := range units {
		isPetAllowed := true
		for _, pet := range userPets {
			key := pet.PetType + "-" + pet.PetSize
			if allowed, ok := petPolicyMap[key]; !ok || !allowed {
				isPetAllowed = false
				break
			}
		}
		unitInfo := types.UnitInfoResp{
			UnitRentID:             unit.UnitRentID,
			MonthlyRent:            unit.MonthlyRent,
			SquareFootage:          unit.SquareFootage,
			AvailableDateForMoveIn: unit.AvailableDateForMoveIn,
			IsPetAllowed:           isPetAllowed,
		}
		unitInfos = append(unitInfos, unitInfo)
	}
	return unitInfos, nil
}

func (dao *TaskDao) UpdatePet(req *types.UpdatePets, username string) error {
	query := "SELECT pet_name, pet_type, pet_size FROM pets WHERE pet_name = ? AND pet_type = ? AND username = ?"
	row := dao.DB.QueryRow(query, req.CurrentPetName, req.CurrentPetType, username)

	var pet model.Pets
	if err := row.Scan(&pet.PetName, &pet.PetType, &pet.PetSize); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Pet not found")
		}
		return err
	}
	updateQuery := "UPDATE pets SET pet_name = ?, pet_type = ?, pet_size = ? WHERE pet_name = ? AND pet_type = ? AND username = ?"
	_, err := dao.DB.Exec(
		updateQuery,
		req.NewPetName, req.NewPetType, req.NewPetSize,
		req.CurrentPetName, req.CurrentPetType, username,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("DUPLICATE_KEY: pet already exists")
		}
		return err
	}

	return nil
}

func (dao *TaskDao) GetPet(username string) ([]model.Pets, error) {
	query := "SELECT pet_name, pet_type, pet_size FROM pets WHERE username = ?"
	rows, err := dao.DB.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pets []model.Pets
	for rows.Next() {
		var pet model.Pets
		if err := rows.Scan(&pet.PetName, &pet.PetType, &pet.PetSize); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return pets, nil
}

func (dao *TaskDao) CreatePet(req *types.GetPets, username string) error {
	query := "INSERT INTO pets (pet_name, pet_type, pet_size, username) VALUES (?, ?, ?, ?)"
	_, err := dao.DB.Exec(query, req.CurrentPetName, req.CurrentPetType, req.CurrentPetSize, username)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("DUPLICATE_KEY: pet already exists")
		}
		return err
	}
	return nil
}

func (dao *TaskDao) GetInterests(unitRentID int) ([]model.Interests, error) {
	query := "SELECT username, unit_rent_id, roommate_cnt, move_in_date FROM interests WHERE unit_rent_id = ?"
	rows, err := dao.DB.Query(query, unitRentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var interests []model.Interests
	for rows.Next() {
		var interest model.Interests
		if err := rows.Scan(&interest.Username, &interest.UnitRentID, &interest.RoommateCnt, &interest.MoveInDate); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return interests, nil
}

func (dao *TaskDao) CreateInterests(req *types.PostInterestReq, username string) error {
	query := `
		INSERT INTO interests (username, unit_rent_id, roommate_cnt, move_in_date)
		VALUES (?, ?, ?, ?)
	`
	_, err := dao.DB.Exec(
		query,
		username,
		req.UnitRentID,
		req.RoommateCnt,
		req.MoveInDate,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("DUPLICATE_KEY")
		}
		return err
	}
	return nil
}

func (dao *TaskDao) GetApartmentUnitByUnitRentID(unitRentID int) (*model.ApartmentUnit, error) {
	query := "SELECT unit_rent_id, company_name, building_name, unit_number, monthly_rent, square_footage, available_date_for_move_in FROM ApartmentUnit WHERE unit_rent_id = ?"
	row := dao.DB.QueryRow(query, unitRentID)
	var unit model.ApartmentUnit
	if err := row.Scan(&unit.UnitRentID, &unit.CompanyName, &unit.BuildingName, &unit.UnitNumber, &unit.MonthlyRent, &unit.SquareFootage, &unit.AvailableDateForMoveIn); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

func (dao *TaskDao) GetApartmentBuildingByUnitRentID(unitRentID int) (*model.ApartmentBuilding, error) {
	queryUnit := "SELECT company_name, building_name FROM ApartmentUnit WHERE unit_rent_id = ?"
	row := dao.DB.QueryRow(queryUnit, unitRentID)

	var unit model.ApartmentUnit
	if err := row.Scan(&unit.CompanyName, &unit.BuildingName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	queryBuilding := "SELECT company_name, building_name, addr_num, addr_street, addr_city, addr_state, addr_zip_code, year_built FROM ApartmentBuilding WHERE company_name = ? AND building_name = ?"
	rowBuilding := dao.DB.QueryRow(queryBuilding, unit.CompanyName, unit.BuildingName)
	var building model.ApartmentBuilding
	if err := rowBuilding.Scan(&building.CompanyName, &building.BuildingName, &building.AddrNum, &building.AddrStreet, &building.AddrCity, &building.AddrState, &building.AddrZipCode, &building.YearBuilt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &building, nil
}
func (dao *TaskDao) GetProvidesByUnitRentID(unitRentID int) ([]model.Provides, error) {
	queryUnit := "SELECT company_name, building_name FROM ApartmentUnit WHERE unit_rent_id = ?"
	rowUnit := dao.DB.QueryRow(queryUnit, unitRentID)
	var unit model.ApartmentUnit
	if err := rowUnit.Scan(&unit.CompanyName, &unit.BuildingName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	queryProvides := "SELECT a_type, company_name, building_name, fee, waiting_list FROM provides WHERE company_name = ? AND building_name = ?"
	rows, err := dao.DB.Query(queryProvides, unit.CompanyName, unit.BuildingName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var provides []model.Provides
	for rows.Next() {
		var provide model.Provides
		if err := rows.Scan(&provide.AType, &provide.CompanyName, &provide.BuildingName, &provide.Fee, &provide.WaitingList); err != nil {
			return nil, err
		}
		provides = append(provides, provide)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return provides, nil
}
func (dao *TaskDao) GetAmenitiesInByUnitRentID(unitRentID int) ([]model.AmenitiesIn, error) {
	query := "SELECT a_type, unit_rent_id FROM AmenitiesIn WHERE unit_rent_id = ?"
	rows, err := dao.DB.Query(query, unitRentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var amenities []model.AmenitiesIn
	for rows.Next() {
		var amenity model.AmenitiesIn
		if err := rows.Scan(&amenity.AType, &amenity.UnitRentID); err != nil {
			return nil, err
		}
		amenities = append(amenities, amenity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return amenities, nil
}

func (dao *TaskDao) CountAvailableUnitsByUnitRentID(unitRentID int) (int, error) {
	queryUnit := "SELECT company_name, building_name FROM ApartmentUnit WHERE unit_rent_id = ?"
	rowUnit := dao.DB.QueryRow(queryUnit, unitRentID)
	var unit model.ApartmentUnit
	if err := rowUnit.Scan(&unit.CompanyName, &unit.BuildingName); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	queryCount := `
		SELECT COUNT(*) FROM ApartmentUnit
		WHERE company_name = ? AND building_name = ?
		AND available_date_for_move_in IS NOT NULL
	`
	rowCount := dao.DB.QueryRow(queryCount, unit.CompanyName, unit.BuildingName)
	var count int
	if err := rowCount.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (dao *TaskDao) GetRoomCountsByUnitRentID(unitRentID int) (int, int, int, error) {
	query := "SELECT name FROM rooms WHERE unit_rent_id = ?"
	rows, err := dao.DB.Query(query, unitRentID)
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()
	var bedroomCount, bathroomCount, livingRoomCount int
	bedroomRegex := regexp.MustCompile(`(?i)bedroom\d*`)
	bathroomRegex := regexp.MustCompile(`(?i)bathroom\d*`)
	livingRoomRegex := regexp.MustCompile(`(?i)livingroom\d*`)
	for rows.Next() {
		var roomName string
		if err := rows.Scan(&roomName); err != nil {
			return 0, 0, 0, err
		}
		if bedroomRegex.MatchString(roomName) {
			bedroomCount++
		} else if bathroomRegex.MatchString(roomName) {
			bathroomCount++
		} else if livingRoomRegex.MatchString(roomName) {
			livingRoomCount++
		}
	}
	return bedroomCount, bathroomCount, livingRoomCount, nil
}

func (dao *TaskDao) SearchInterestswithcond(unitRentID int, moveInDate *types.CustomTime, roommateCnt *uint8) ([]model.Interests, error) {
	dbLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}
	var moveInDateValue *time.Time
	if moveInDate != nil {
		dateStr := moveInDate.Time.Format("2006-01-02")
		parsedTime, err := time.ParseInLocation("2006-01-02", dateStr, dbLocation)
		if err != nil {
			return nil, err
		}
		moveInDateValue = &parsedTime
	}
	query := "SELECT username, unit_rent_id, roommate_cnt, move_in_date FROM interests WHERE unit_rent_id = ?"
	args := []interface{}{unitRentID}
	if moveInDateValue != nil {
		query += " AND move_in_date = ?"
		args = append(args, *moveInDateValue)
	}
	if roommateCnt != nil {
		query += " AND roommate_cnt = ?"
		args = append(args, *roommateCnt)
	}
	rows, err := dao.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []model.Interests
	for rows.Next() {
		var interest model.Interests
		if err := rows.Scan(&interest.Username, &interest.UnitRentID, &interest.RoommateCnt, &interest.MoveInDate); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return interests, nil
}
