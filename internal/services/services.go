package services

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	adapters "taxApi/internal/adapter"
	"taxApi/internal/logs"
	"taxApi/internal/models"
	"taxApi/internal/repositories"
)

type Service struct {
	Adapter    *adapters.Adapter
	Repository *repositories.Repository
	Lg         *logs.Logger
}

func NewService(adapter *adapters.Adapter, repository *repositories.Repository, lg *logs.Logger) *Service {
	return &Service{Adapter: adapter, Repository: repository, Lg: lg}
}

func (s *Service) IsLoginUsed(login string) error {
	if s.Repository.IsLoginUsed(login) {
		return errors.New("already registered with this login")
	}

	return nil
}

func (s *Service) RegistrationUser(user *models.User) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.Repository.RegistrationUser(user, hash)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ValidateUser(user *models.User) error {
	if len(user.FullName) > 20 || len(user.FullName) < 3 {
		return errors.New("invalid data in FIO")
	}
	if len(user.Login) > 20 || len(user.Login) < 3 {
		return errors.New("invalid data in Login")
	}
	if len(user.Password) > 20 || len(user.Password) <= 7 {
		return errors.New("invalid data in Password")
	}
	return nil
}

func (s *Service) ValidateLoginAndPass(login, pass string) (*models.User, error) {
	user, err := s.Repository.ValidateLogin(login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetRoute(reqBody *models.RequestBody, id int) (models.ResponseBody, float64, error) {
	getRoute, err := s.Adapter.GetRoute(reqBody)
	if err != nil {
		return models.ResponseBody{}, 0, err
	}
	carWithBill, err := s.Repository.GetBill(reqBody.Fare)
	if err != nil {
		return models.ResponseBody{}, 0, err
	}

	lenRoute := getRoute.Distances[0][0]
	priceOfOrder := s.CostPrice(lenRoute, float64(carWithBill.MinPrice))

	durationRoute := getRoute.Durations[0][0]
	distanceRoute := getRoute.Distances[0][0]
	sourceOfRouter := getRoute.Sources[0].Name
	destinationOfRoute := getRoute.Destinations[0].Name

	var order = models.Order{
		Price:       priceOfOrder,
		Source:      sourceOfRouter,
		Destination: destinationOfRoute,
		Duration:    durationRoute,
		Distance:    distanceRoute,
		CustomerId:  id,
		Status:      "В ожидании",
		Fare:        reqBody.Fare,
		Active:      true,
	}

	err = s.Repository.PutOrder(&order)
	if err != nil {
		return models.ResponseBody{}, 0, err
	}

	return getRoute, priceOfOrder, nil
}

func (s *Service) GetTravelToCustomer(id int) (*[]models.OrderResponseToCustomer, error) {
	return s.Repository.GetTravelToCustomer(id)
}

func (s *Service) GetTravelToDriver(id int) ([]models.ResponseToDriver, error) {
	bill, err := s.Repository.GetFareDriver(id)
	if err != nil {
		return []models.ResponseToDriver{}, err
	}
	travelToDriver, err := s.Repository.GetTravelToDriver(bill.Fare)
	if err != nil {
		return []models.ResponseToDriver{}, err
	}
	var orderToDriver []models.ResponseToDriver
	for _, order := range *travelToDriver {
		if order.Fare != bill.Fare {
			continue
		}
		orderToDriver = append(orderToDriver, order)
	}

	return orderToDriver, nil
}

func (s *Service) GetTravelById(orderToDriver []models.ResponseToDriver, id, idUser int) (models.ResponseToDriver, error) {
	var nilOrder models.ResponseToDriver
	for _, toDriver := range orderToDriver {
		if toDriver.Id == id {
			nilOrder = toDriver
		}
		continue
	}
	err := s.Repository.ChangeStatus(nilOrder.Id, idUser)
	if err != nil {
		return models.ResponseToDriver{}, err
	}

	//ранее была горутина, которая по истечению времени поездки  - меняла статус на завершен
	/*duration, err := time.ParseDuration(fmt.Sprintf("%f", nilOrder.Duration) + "s")
	if err != nil {
		return models.ResponseToDriver{}, err
	}
	go func() {
		time.Sleep(duration)
		s.Repository.EndOrder(nilOrder.Id)
	}()*/

	return nilOrder, nil
}

func (s *Service) EndTravelById(idUser, idTravel int) error {
	return s.Repository.EndTravel(idUser, idTravel)
}

func (s *Service) CostPrice(lenRoute float64, price float64) float64 {

	if lenRoute <= 2500 {
		return price
	} else {
		price = (float64(lenRoute-2500)*2.5)/1000 + price
		return price
	}

}

func (s *Service) ReportDriver(report *models.ReqReport, id int) (error, *xlsx.File) {

	reportDriver, err := s.Repository.ReportDriver(report, id)
	if err != nil {
		return err, &xlsx.File{}
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return err, &xlsx.File{}
	}

	colNames := []string{"№", "Дата поездки", "Цена", "Пункт отправления", "Пункт прибытия", "Время поездки (с)", "ФИО клиента", "Номер клиента"}
	row := sheet.AddRow()
	for _, colName := range colNames {
		cell := row.AddCell()
		cell.Value = colName
	}
	i := 1
	totalSum := 0.0
	for _, record := range reportDriver {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = strconv.Itoa(i)
		i++
		cell = row.AddCell()
		timeFormat := record.CreatedAt.Format("2006-01-02 15:04:05")
		cell.Value = timeFormat
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", record.Price)
		cell = row.AddCell()
		totalSum = totalSum + record.Price
		cell.Value = record.Source
		cell = row.AddCell()
		cell.Value = record.Destination
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", record.Duration)
		cell = row.AddCell()
		cell.Value = record.CustomerName
		cell = row.AddCell()
		cell.Value = record.CustomerPhone
	}

	row = sheet.AddRow()
	cell := row.AddCell()
	cell = row.AddCell()
	cell.Value = "ИТОГО"
	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", totalSum)

	//startReport := report.Period[0].Format("2006-01-02")
	//endReport := report.Period[1].Format("2006-01-02")
	//err = file.Save("reports/Отчет за период с " + startReport + " по " + endReport + ".xlsx")
	//err = file.Save("report.xlsx")
	//if err != nil {
	//	return err, &xlsx.File{}
	//}

	return nil, file
}

func (s *Service) Report(report *models.ReqReport) error {

	reports, err := s.Repository.Report(report)
	if err != nil {
		return err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	colNames := []string{"№", "Дата поездки", "Цена", "Пункт отправления", "Пункт прибытия", "Расстояние (м)", "ФИО клиента", "Номер клиента", "ФИО водителя", "Номер водителя"}
	row := sheet.AddRow()
	for _, colName := range colNames {
		cell := row.AddCell()
		cell.Value = colName
	}
	i := 1
	totalSum := 0.0
	for _, record := range reports {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = strconv.Itoa(i)
		i++
		cell = row.AddCell()
		timeFormat := record.CreatedAt.Format("2006-01-02 15:04:05")
		cell.Value = timeFormat
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", record.Price)
		cell = row.AddCell()
		totalSum = totalSum + record.Price
		cell.Value = record.Source
		cell = row.AddCell()
		cell.Value = record.Destination
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%.2f", record.Duration)
		cell = row.AddCell()
		cell.Value = record.CustomerName
		cell = row.AddCell()
		cell.Value = record.CustomerPhone
		cell = row.AddCell()
		cell.Value = record.DriverName
		cell = row.AddCell()
		cell.Value = record.DriverPhone
	}

	row = sheet.AddRow()
	cell := row.AddCell()
	cell = row.AddCell()
	cell.Value = "Итог"
	cell = row.AddCell()
	cell.Value = fmt.Sprintf("%.2f", totalSum)

	startReport := report.Period[0].Format("2006-01-02")
	endReport := report.Period[1].Format("2006-01-02")
	err = file.Save("reports/common/Отчет за период с " + startReport + " по " + endReport + ".xlsx")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ReadOrderById(id int) ([]*models.ResReport, []*models.ResCustomer, error) {
	reportsDriver, reportsCustomer, err := s.Repository.ReadOrderById(id)
	if err != nil {
		return []*models.ResReport{}, []*models.ResCustomer{}, err
	}

	return reportsDriver, reportsCustomer, nil
}
