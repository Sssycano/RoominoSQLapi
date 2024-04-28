package service

import (
	"context"
	"errors"
	"fmt"
	"roomino/ctl"
	"roomino/dao"
	"roomino/types"
	"sync"
	"time"
)

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})
	return TaskSrvIns
}

func (s *TaskSrv) GetAvailableUnitsWithPetPolicy(ctx context.Context, req *types.UnitInforeq) (resp interface{}, err error) {
	taskDao := dao.NewTaskDao(ctx)
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	units, err := taskDao.GetUnitsWithPetPolicy(req.CompanyName, req.BuildingName, u.UserName)
	if err != nil {
		return nil, errors.New("failed to retrieve units")
	}
	var unitResp []types.UnitInfoResp
	for _, unit := range units {
		unitResp = append(unitResp, types.UnitInfoResp{
			UnitRentID:             unit.UnitRentID,
			MonthlyRent:            unit.MonthlyRent,
			SquareFootage:          unit.SquareFootage,
			AvailableDateForMoveIn: unit.AvailableDateForMoveIn,
			IsPetAllowed:           unit.IsPetAllowed,
		})
	}

	return ctl.RespSuccessWithData(unitResp), nil
}

func (s *TaskSrv) UpdatePetInfo(ctx context.Context, req *types.UpdatePets) (resp interface{}, err error) {
	taskDao := dao.NewTaskDao(ctx)
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}
	if err := taskDao.UpdatePet(req, u.UserName); err != nil {
		return nil, err
	}

	return ctl.RespSuccess(), nil
}

func (s *TaskSrv) GetPet(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, errors.New("failed to get user info")
	}

	taskDao := dao.NewTaskDao(ctx)
	pets, err := taskDao.GetPet(u.UserName)
	if err != nil {
		return nil, errors.New("failed to retrieve pets")
	}
	var petsResp []types.GetPets
	for _, pet := range pets {
		petsResp = append(petsResp, types.GetPets{
			CurrentPetName: pet.PetName,
			CurrentPetType: pet.PetType,
			CurrentPetSize: pet.PetSize,
		})
	}
	return ctl.RespSuccessWithData(petsResp), nil
}

func (s *TaskSrv) CreatePet(ctx context.Context, req *types.GetPets) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	taskDao := dao.NewTaskDao(ctx)
	if err := taskDao.CreatePet(req, u.UserName); err != nil {
		return nil, err
	}
	return ctl.RespSuccess(), nil
}

func (s *TaskSrv) GetInterests(ctx context.Context, req *types.UnitRentIDReq) (interface{}, error) {
	taskDao := dao.NewTaskDao(ctx)
	unitRentID := req.UnitRentID
	interests, err := taskDao.GetInterests(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve interests")
	}
	var interestsResp []types.InterestResp
	for _, interest := range interests {
		interestsResp = append(interestsResp, types.InterestResp{
			Username:    interest.Username,
			UnitRentID:  interest.UnitRentID,
			RoommateCnt: interest.RoommateCnt,
			MoveInDate:  interest.MoveInDate,
		})
	}
	return ctl.RespSuccessWithData(interestsResp), nil
}
func (s *TaskSrv) CreateInterests(ctx context.Context, req *types.PostInterestReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, errors.New("failed to get user info")
	}
	taskDao := dao.NewTaskDao(ctx)
	if err := taskDao.CreateInterests(req, u.UserName); err != nil {
		if err.Error() == "DUPLICATE_KEY: Interest already exists" {
			return nil, errors.New("duplicate interest")
		}

		return nil, errors.New("failed to create interest")
	}
	return ctl.RespSuccess(), nil
}

func (s *TaskSrv) GetComplexUnitinfo(ctx context.Context, req *types.UnitRentIDReq) (interface{}, error) {
	taskDao := dao.NewTaskDao(ctx)
	unitRentID := req.UnitRentID

	unit, err := taskDao.GetApartmentUnitByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve apartment unit")
	}

	building, err := taskDao.GetApartmentBuildingByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve apartment building")
	}

	moveInDate, err := time.Parse("2006-01-02T15:04:05-07:00", unit.AvailableDateForMoveIn)
	if err != nil {
		return nil, errors.New("failed to parse AvailableDateForMoveIn")
	}

	provides, err := taskDao.GetProvidesByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve provides")
	}
	amenitiesinbuilding := make([]string, len(provides))
	for i, p := range provides {
		amenitiesinbuilding[i] = p.AType
	}
	availableUnitsCount, err := taskDao.CountAvailableUnitsByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve available units count")
	}

	amenitiesIn, err := taskDao.GetAmenitiesInByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve amenities in")
	}
	amenitiesinunit := make([]string, len(amenitiesIn))
	for i, a := range amenitiesIn {
		amenitiesinunit[i] = a.AType
	}
	bedroomCount, bathroomCount, livingRoomCount, err := taskDao.GetRoomCountsByUnitRentID(unitRentID)
	if err != nil {
		return nil, errors.New("failed to retrieve room counts")
	}

	complexInfo := types.ComplexUnitinfo{
		CompanyName:            building.CompanyName,
		BuildingName:           building.BuildingName,
		Address:                fmt.Sprintf("%d %s, %s, %s %s", building.AddrNum, building.AddrStreet, building.AddrCity, building.AddrState, building.AddrZipCode),
		YearBuilt:              building.YearBuilt,
		Amenitiesinbuding:      amenitiesinbuilding,
		Amenitiesinunit:        amenitiesinunit,
		AvailableUnits:         availableUnitsCount,
		UnitRentID:             unit.UnitRentID,
		MonthlyRent:            unit.MonthlyRent,
		SquareFootage:          unit.SquareFootage,
		AvailableDateForMoveIn: moveInDate,
		BedroomNum:             bedroomCount,
		BathroomNum:            bathroomCount,
		LivingRoomNum:          livingRoomCount,
	}

	return ctl.RespSuccessWithData(complexInfo), nil
}

func (s *TaskSrv) SearchInterestswithcond(ctx context.Context, req *types.InteresCondReq) (interface{}, error) {
	taskDao := dao.NewTaskDao(ctx)
	var moveInDate *types.CustomTime
	if !req.MoveInDate.IsZero() {
		moveInDate = &req.MoveInDate
	}

	var roommateCnt *uint8
	if req.RoommateCnt != 0 {
		roommateCnt = &req.RoommateCnt
	}

	interests, err := taskDao.SearchInterestswithcond(req.UnitRentID, moveInDate, roommateCnt)
	if err != nil {
		return nil, errors.New("failed to search interests")
	}
	var interestsResp []types.InterestResp
	for _, interest := range interests {
		interestsResp = append(interestsResp, types.InterestResp{
			Username:    interest.Username,
			UnitRentID:  interest.UnitRentID,
			RoommateCnt: interest.RoommateCnt,
			MoveInDate:  interest.MoveInDate,
		})
	}
	return ctl.RespSuccessWithData(interestsResp), nil
}
