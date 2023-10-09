package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"taxApi/internal/logs"
	"taxApi/internal/models"
	"time"
)

type Repository struct {
	Db *gorm.DB
	Lg *logs.Logger
}

func NewRepository(db *gorm.DB, lg *logs.Logger) *Repository {
	return &Repository{Db: db, Lg: lg}
}

func (r *Repository) RegistrationUser(user *models.User, hash []byte) error {
	sqlQuery := `insert into users (full_name,login,password,phone_number,role)
					values (?,?,?,?,?) returning id;`
	sqlQuery1 := r.Db.Raw(sqlQuery, user.FullName, user.Login, hash, user.PhoneNumber, user.Role).Scan(&user)
	err := sqlQuery1.Error
	if err != nil {
		return err
	}
	if user.Role == "driver" {
		err = r.Db.Exec(`insert into cars (title, user_id) values (?,?);`, user.CarTitle, user.Id).Error
		if err != nil {
			return err
		}
	}

	sqlRole := `insert into roles (title, user_id) values (?,?);`
	err = r.Db.Exec(sqlRole, user.Role, user.Id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IsLoginUsed(login string) bool {
	var user *models.User
	sqlQuery := `select login from users where login = ?;`

	sql := r.Db.Raw(sqlQuery, login).Scan(&user)
	if sql.RowsAffected == 0 {
		return false
	}

	err := sql.Error
	if err != nil {
		return false
	}
	return true
}

func (r *Repository) ValidateLogin(login string) (user *models.User, err error) {

	err = r.Db.Where("login = ?", login).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) PutNewToken(ctx context.Context, token string, id int64) error {

	sqlQuery := `insert into users_tokens (user_id,token) values ($1,$2);`
	_, err := r.Db.ConnPool.ExecContext(ctx, sqlQuery, id, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IdByToken(cxt context.Context, token string) (id int, expire time.Time, err error) {
	sqlQuery := `select user_id,expire_at from users_tokens where token = $1;`
	row := r.Db.ConnPool.QueryRowContext(cxt, sqlQuery, token)
	errString := row.Scan(&id, &expire)
	if errString != nil {
		return 0, expire, err
	}
	return id, expire, nil
}

func (r *Repository) PutOrder(order *models.Order) error {
	sqlQuery := `insert into orders (price, source,destination, duration,distance, status,customer_id, fare)
					values (?,?,?,?,?,?,?,?);`

	err := r.Db.Exec(sqlQuery, order.Price, order.Source, order.Destination, order.Duration, order.Distance, order.Status, order.CustomerId, order.Fare).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTravelToCustomer(id int) (*[]models.OrderResponseToCustomer, error) {
	var order *[]models.OrderResponseToCustomer
	sqlQuery := `select orders.id, orders.updated_at travel_date, orders.price, orders.source, orders.destination,u.full_name driver_name, u.phone_number driver_phone, orders.status, orders.fare
		from orders join users u on u.id = orders.driver_id
		where orders.customer_id = ?;`
	rows := r.Db.Raw(sqlQuery, id).Scan(&order)

	if rows.RowsAffected == 0 {
		return &[]models.OrderResponseToCustomer{}, errors.New("заказов на данный момент нет")
	}
	err := rows.Error
	if err != nil {
		return &[]models.OrderResponseToCustomer{}, err
	}

	return order, nil
}

func (r *Repository) GetFareDriver(id int) (*models.Billing, error) {
	var idCar int
	sqlQuery := `select id from cars where user_id = ?;`
	err := r.Db.Raw(sqlQuery, id).Scan(&idCar).Error
	if err != nil {
		return &models.Billing{}, err
	}

	var bill *models.Billing
	sqlQuery2 := `select * from billing where car_id = ?;`
	err = r.Db.Raw(sqlQuery2, idCar).Scan(&bill).Error
	if err != nil {
		return &models.Billing{}, err
	}

	return bill, nil
}

func (r *Repository) ChangeStatus(id, idUser int) error {
	sqlQuery := `update orders set status = 'Заказ принят',driver_id = ?, updated_at = current_timestamp where id = ?;`
	err := r.Db.Exec(sqlQuery, idUser, id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) EndOrder(orderId int) error {
	sqlQuery := `update orders set status = 'Завершен', updated_at = current_timestamp where id = ?`
	err := r.Db.Exec(sqlQuery, orderId).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) EndTravel(idUser, idTravel int) error {
	var order *models.Order
	sqlQuery := `select * from orders where id = ? and driver_id = ? and status =?;`
	statusOrder := "Заказ принят"
	rows := r.Db.Raw(sqlQuery, idTravel, idUser, statusOrder).Scan(&order)
	if rows.RowsAffected == 0 {
		return errors.New("вы пытаетесь завершить неверный заказ")
	}

	err := rows.Error
	if err != nil {
		return err
	}

	sqlQuery2 := `update orders set status = 'Завершен', driver_id =?,updated_at = current_timestamp where id = ?`
	err = r.Db.Exec(sqlQuery2, idUser, idTravel).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTravelToDriver(billing string) (*[]models.ResponseToDriver, error) {
	var order *[]models.ResponseToDriver
	sqlQuery := `select orders.id, source, destination, duration, price, status, u.full_name customer_name, u.phone_number customer_phone, fare
				from orders join users u on u.id = orders.customer_id
				where status = ? and fare =?;`

	rows := r.Db.Raw(sqlQuery, "В ожидании", billing).Scan(&order)
	if rows.RowsAffected == 0 {
		return &[]models.ResponseToDriver{}, errors.New("заказов на данный момент нет")
	}
	err := rows.Error
	if err != nil {
		return &[]models.ResponseToDriver{}, err
	}

	return order, nil
}

func (r *Repository) GetBill(fare string) (*models.Billing, error) {
	var bill *models.Billing
	sqlQuery := `select distinct * from billing where fare = ? limit 1;`
	err := r.Db.Raw(sqlQuery, fare).Scan(&bill).Error
	if err != nil {
		return &models.Billing{}, err
	}
	return bill, nil
}

func (r *Repository) GetDriver(id int) (*models.Car, error) {
	var car *models.Car
	sqlQuery1 := `select * from cars where id = ?;`
	err := r.Db.Raw(sqlQuery1, id).Scan(&car).Error
	if err != nil {
		return &models.Car{}, err
	}

	return car, nil
}

func (r *Repository) ReportDriver(report *models.ReqReport, id int) ([]*models.ResReport, error) {

	var orderResponse []*models.ResReport
	startDate := report.Period[0]
	endDate := report.Period[1]
	sqlQuery := `select price,source,destination,duration, users.full_name customer_name, users.phone_number customer_phone, orders.created_at 
				 from orders join users on users.id = orders.customer_id
				 where orders.driver_id = ? and (orders.created_at >= ? and orders.created_at <= ?) limit 10 offset 0;`
	err := r.Db.Raw(sqlQuery, id, startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")).Scan(&orderResponse).Error
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

func (r *Repository) Report(report *models.ReqReport) ([]*models.RepResponse, error) {

	var orderResponse []*models.RepResponse
	startDate := report.Period[0]
	endDate := report.Period[1]
	sqlQuery := `select o.created_at, o.price, o.source, o.destination, o.duration, Customers.full_name as customer_name, Customers.phone_number as customer_phone, Drivers.full_name as driver_name, Drivers.phone_number as driver_phone
					from orders o
         				left join users as Customers on o.customer_id = Customers.id
         				left join users as Drivers on o.driver_id= Drivers.id
					where (o.created_at>=? and o.created_at <=?)
					order by o.created_at desc
					limit 15 offset 0;`
	err := r.Db.Raw(sqlQuery, startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")).Scan(&orderResponse).Error
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

func (r *Repository) ReadOrderById(id int) ([]*models.ResReport, []*models.ResCustomer, error) {
	var role string
	sqlRow := r.Db.Raw(`select role from users where id =?;`, id).Scan(&role)
	rows := sqlRow.RowsAffected
	if rows == 0 {
		return []*models.ResReport{}, []*models.ResCustomer{}, errors.New("нет такого пользователя")
	}

	err := sqlRow.Error
	if err != nil {
		return []*models.ResReport{}, []*models.ResCustomer{}, err
	}

	if role == "customer" {
		var orderCustomer []*models.ResCustomer
		sqlQuery := `select price,source,destination,duration, users.full_name driver_name, users.phone_number driver_phone, orders.created_at 
				 from orders join users on users.id = orders.driver_id
				 where orders.customer_id = ?;`
		err = r.Db.Raw(sqlQuery, id).Scan(&orderCustomer).Error
		if err != nil {
			return []*models.ResReport{}, []*models.ResCustomer{}, err
		}
		if len(orderCustomer) == 0 {
			return []*models.ResReport{}, []*models.ResCustomer{}, errors.New("У данного клиента нет поездок")
		}

		return []*models.ResReport{}, orderCustomer, nil

	}

	var orderResponse []*models.ResReport
	sqlQuery := `select price,source,destination,duration, users.full_name customer_name, users.phone_number customer_phone, orders.created_at 
				 from orders join users on users.id = orders.customer_id
				 where orders.driver_id = ?;`
	err = r.Db.Raw(sqlQuery, id).Scan(&orderResponse).Error
	if err != nil {
		return []*models.ResReport{}, []*models.ResCustomer{}, err
	}

	if len(orderResponse) == 0 {
		return []*models.ResReport{}, []*models.ResCustomer{}, errors.New("У данного водителя нет поездок")
	}
	return orderResponse, []*models.ResCustomer{}, nil
}
