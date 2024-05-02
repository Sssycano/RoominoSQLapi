package dao

import (
	"context"
	"database/sql"
	"errors"
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

	query := `SELECT au.UnitRentID, au.MonthlyRent, au.squareFootage, au.AvailableDateForMoveIn,
					CASE 
						WHEN COUNT(IF(pp.IsAllowed = 0, 1, NULL)) > 0 THEN 0 
						ELSE 1 
					END AS IsPetAllowed
				FROM ApartmentUnit au
				LEFT JOIN Pets p ON p.username = ?
				LEFT JOIN PetPolicy pp ON pp.CompanyName = au.CompanyName AND pp.BuildingName = au.BuildingName 
						AND pp.PetType = p.PetType AND pp.PetSize = p.PetSize
				WHERE 
					au.CompanyName = ? AND au.BuildingName = ? 
				GROUP BY 
					au.UnitRentID
			`

	rows, err := dao.DB.Query(query, username, companyName, buildingName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var unitInfos []types.UnitInfoResp
	for rows.Next() {
		var unitInfo types.UnitInfoResp
		if err := rows.Scan(
			&unitInfo.UnitRentID,
			&unitInfo.MonthlyRent,
			&unitInfo.SquareFootage,
			&unitInfo.AvailableDateForMoveIn,
			&unitInfo.IsPetAllowed,
		); err != nil {
			return nil, err
		}
		unitInfos = append(unitInfos, unitInfo)
	}
	return unitInfos, nil
}

func (dao *TaskDao) GetPetPolicies(companyName, buildingName, username string) ([]types.PetPolicy, error) {
	query := `
		SELECT 
			pp.PetType, 
			pp.PetSize, 
			pp.IsAllowed,
			pp.RegistrationFee,
			pp.MonthlyFee
		FROM 
			PetPolicy pp
		JOIN 
			Pets p ON p.PetType = pp.PetType 
				AND p.PetSize = pp.PetSize 
				AND p.username = ?
		WHERE 
			pp.CompanyName = ? 
			AND pp.BuildingName = ? 
	`

	rows, err := dao.DB.Query(query, username, companyName, buildingName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var petPolicies []types.PetPolicy
	for rows.Next() {
		var petPolicy types.PetPolicy
		if err := rows.Scan(
			&petPolicy.PetType,
			&petPolicy.PetSize,
			&petPolicy.IsAllowed,
			&petPolicy.RegistrationFee,
			&petPolicy.MonthlyFee,
		); err != nil {
			return nil, err
		}

		petPolicies = append(petPolicies, petPolicy)
	}

	return petPolicies, nil
}

func (dao *TaskDao) UpdatePet(req *types.UpdatePets, username string) error {
	query := "SELECT PetName, PetType, PetSize FROM pets WHERE PetName = ? AND PetType = ? AND username = ?"
	row := dao.DB.QueryRow(query, req.CurrentPetName, req.CurrentPetType, username)

	var pet model.Pets
	if err := row.Scan(&pet.PetName, &pet.PetType, &pet.PetSize); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("PET not found")
		}
		return err
	}
	updateQuery := "UPDATE pets SET PetName = ?, PetType = ?, PetSize = ? WHERE PetName = ? AND PetType = ? AND username = ?"
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
	query := "SELECT PetName, PetType, PetSize FROM pets WHERE username = ?"
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
	query := "INSERT INTO pets (PetName, PetType, PetSize, username) VALUES (?, ?, ?, ?)"
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
	query := "SELECT username, UnitRentID, RoommateCnt, MoveInDate FROM interests WHERE UnitRentID= ?"
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
	query := "INSERT INTO interests (username, UnitRentID, RoommateCnt, MoveInDate) VALUES (?, ?, ?, ?)"
	dbLocation, _ := time.LoadLocation("America/New_York")

	var moveInDateValue *time.Time

	dateStr := req.MoveInDate.Time.Format("2006-01-02")
	parsedTime, _ := time.ParseInLocation("2006-01-02", dateStr, dbLocation)

	moveInDateValue = &parsedTime

	_, err := dao.DB.Exec(
		query,
		username,
		req.UnitRentID,
		moveInDateValue,
		req.MoveInDate,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("DUPLICATE_KEY: Interests already exists")
		}
		return err
	}
	return nil
}

func (dao *TaskDao) GetApartmentUnitByUnitRentID(unitRentID int) (*model.ApartmentUnit, error) {
	query := "SELECT UnitRentID, CompanyName, BuildingName, unitNumber, MonthlyRent, squareFootage, AvailableDateForMoveIn FROM ApartmentUnit WHERE UnitRentID= ?"
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
	queryUnit := "SELECT CompanyName, BuildingName FROM ApartmentUnit WHERE UnitRentID= ?"
	row := dao.DB.QueryRow(queryUnit, unitRentID)

	var unit model.ApartmentUnit
	if err := row.Scan(&unit.CompanyName, &unit.BuildingName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	queryBuilding := "SELECT CompanyName, BuildingName, AddrNum, AddrStreet, AddrCity, AddrState, AddrZipCode, YearBuilt FROM ApartmentBuilding WHERE CompanyName = ? AND BuildingName = ?"
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
	queryUnit := "SELECT CompanyName, BuildingName FROM ApartmentUnit WHERE UnitRentID= ?"
	rowUnit := dao.DB.QueryRow(queryUnit, unitRentID)
	var unit model.ApartmentUnit
	if err := rowUnit.Scan(&unit.CompanyName, &unit.BuildingName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	queryProvides := "SELECT aType, CompanyName, BuildingName, fee, waitingList FROM provides WHERE CompanyName = ? AND BuildingName = ?"
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
	query := "SELECT aType, UnitRentID FROM AmenitiesIn WHERE UnitRentID= ?"
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
	query := `SELECT COUNT(*) 
  				FROM ApartmentUnit au
			WHERE au.CompanyName = (SELECT CompanyName FROM ApartmentUnit WHERE UnitRentID = ?)
				AND au.BuildingName = (SELECT BuildingName FROM ApartmentUnit WHERE UnitRentID = ?)
				AND au.AvailableDateForMoveIn IS NOT NULL
			`

	row := dao.DB.QueryRow(query, unitRentID, unitRentID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *TaskDao) GetRoomCountsByUnitRentID(unitRentID int) (int, int, int, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN name REGEXP 'bedroom' THEN 1 ELSE NULL END) AS BedroomCount,
			COUNT(CASE WHEN name REGEXP 'bathroom' THEN 1 ELSE NULL END) AS BathroomCount,
			COUNT(CASE WHEN name REGEXP 'livingroom' THEN 1 ELSE NULL END) AS LivingRoomCount
		FROM Rooms r
		WHERE UnitRentID = ?
	`

	row := dao.DB.QueryRow(query, unitRentID)

	var bedroomCount, bathroomCount, livingRoomCount int
	if err := row.Scan(&bedroomCount, &bathroomCount, &livingRoomCount); err != nil {
		return 0, 0, 0, err
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
	query := "SELECT username, UnitRentID, RoommateCnt, MoveInDate FROM interests WHERE UnitRentID= ?"
	args := []interface{}{unitRentID}
	if moveInDateValue != nil {
		query += " AND MoveInDate = ?"
		args = append(args, *moveInDateValue)
	}
	if roommateCnt != nil {
		query += " AND RoommateCnt = ?"
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

func (dao *TaskDao) GetAverageRentByZipAndRoom(addrZipCode string, bedroomNum int, bathroomNum int) (float64, error) {
	query := `
		SELECT AVG(au.MonthlyRent)
		FROM ApartmentUnit au
		JOIN ApartmentBuilding ab ON au.CompanyName = ab.CompanyName AND au.BuildingName = ab.BuildingName
		JOIN Rooms r ON r.UnitRentID = au.UnitRentID
		WHERE ab.AddrZipCode = ?
		AND (
			SELECT COUNT(*)
			FROM Rooms
			WHERE UnitRentID = au.UnitRentID AND name REGEXP ?
		) = ?
		AND (
			SELECT COUNT(*)
			FROM Rooms
			WHERE UnitRentID = au.UnitRentID AND name REGEXP ?
		) = ?
	`
	bedroomPattern := "bedroom.*"
	bathroomPattern := "bathroom.*"
	row := dao.DB.QueryRow(query, addrZipCode, bedroomPattern, bedroomNum, bathroomPattern, bathroomNum)

	var averageRent float64
	if err := row.Scan(&averageRent); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
		return 0, nil
	}
	return averageRent, nil
}
