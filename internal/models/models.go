package models

import (
	"time"
)

type RequestBody struct {
	LongitudeSource      string `json:"longitude_source"`
	LatitudeSource       string `json:"latitude_source"`
	LongitudeDestination string `json:"longitude_destination"`
	LatitudeDestination  string `json:"latitude_destination"`
	Fare                 string `json:"fare"`
}

type ReqReport struct {
	Period [2]time.Time `json:"period"`
}

type Billing struct {
	Id        int       `json:"id"`
	Fare      string    `json:"fare"`
	MinPrice  int       `json:"min_price"`
	CarId     int       `json:"car_id"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type ResponseBody struct {
	Code         string      `json:"code"`
	Distances    [][]float64 `json:"distances"`
	Destinations []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"destinations"`
	Durations [][]float64 `json:"durations"`
	Sources   []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"sources"`
}

type User struct {
	Id          int       `json:"id"`
	FullName    string    `json:"full_name"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CarTitle    string    `json:"car_title"`
	CreatedAt   time.Time `json:"created_at"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type Order struct {
	Id          int       `json:"id"`
	Price       float64   `json:"price"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Duration    float64   `json:"duration"`
	Distance    float64   `json:"distance"`
	DriverId    int       `json:"driver_id"`
	CustomerId  int       `json:"customer_id"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
	Fare        string    `json:"fare"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type ResReport struct {
	Price         float64   `json:"price"`
	Source        string    `json:"source"`
	Destination   string    `json:"destination"`
	Duration      float64   `json:"duration"`
	CustomerName  string    `json:"customer_name"`
	CustomerPhone string    `json:"customer_phone"`
	CreatedAt     time.Time `json:"created_at"`
}

type ResCustomer struct {
	Price       float64   `json:"price"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Duration    float64   `json:"duration"`
	DriverName  string    `json:"driver_name"`
	DriverPhone string    `json:"driver_phone"`
	CreatedAt   time.Time `json:"created_at"`
}

type RepResponse struct {
	CreatedAt     time.Time `json:"created_at"`
	Price         float64   `json:"price"`
	Source        string    `json:"source"`
	Destination   string    `json:"destination"`
	Duration      float64   `json:"duration"`
	CustomerName  string    `json:"customer_name"`
	CustomerPhone string    `json:"customer_phone"`
	DriverName    string    `json:"driver_name"`
	DriverPhone   string    `json:"driver_phone"`
}

type ResponseToDriver struct {
	Id            int     `json:"id"`
	Source        string  `json:"source"`
	Destination   string  `json:"destination"`
	Duration      float64 `json:"duration"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	CustomerName  string  `json:"customer_name"`
	CustomerPhone string  `json:"customer_phone"`
	Fare          string  `json:"fare"`
}

type Car struct {
	Id        int       `json:"id"`
	Title     string    `json:"full_name"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type OrderResponseToDriver struct {
	Message       string  `json:"message"`
	CustomerName  string  `json:"customer_name"`
	CustomerPhone string  `json:"customer_phone"`
	Price         float64 `json:"price"`
}

type OrderResponseToCustomer struct {
	Id          int       `json:"id"`
	TravelDate  time.Time `json:"travel_date"`
	Price       float64   `json:"price"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	DriverName  string    `json:"driver_name"`
	DriverPhone string    `json:"driver_phone"`
	Status      string    `json:"status"`
	Fare        string    `json:"fare"`
}

type Config struct {
	Server  Server   `json:"server"`
	Db      ConfigDb `json:"db"`
	Adapter Adapter  `json:"adapter"`
}

type Adapter struct {
	Url     string `json:"url"`
	Timeout int    `json:"timeout"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ConfigDb struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type Token struct {
	TokenText string `json:"token_text"`
}

type MessageToUser struct {
	Message string `json:"message"`
}

type MessageToCustomer struct {
	Message  string `json:"message"`
	Duration string `json:"duration"`
}

var ErrInvalidData = "Введенные данные недействительны"
var ErrWrongNumberChar = "Неверное количество символов в логине или пароле"
var ErrLoginUsed = "Пользователь с таким логином существует"
var ErrInvalidLoginOrPassword = "Неверный логин или пароль"
var ErrTokenExpired = "Срок Вашего Токена истек, авторизуйтесь заново"
var SuccessRegistration = "Регистрация успешно завершена"

var ErrInternal = "Технические проблемы, попробуйте чуть позже"
var ErrNotOrders = "На данный момент у вас заказов нет"
var ErrNotValidateOrder = "По данному Id у вас нет заказов"

var OkStatus = "Ваша поездка будет составлят "
var ReportStatus = "Отчет за заданный период находится в папке reports"
var ErrNotFound = "Такого пользователя нет в системе, или у данного пользователя нет поездок"
var OrderSuccess = "Заказ принят, номер заказа "
var SuccessOrder = "Заказ успешно выполнен"
var ErrNotFoundOrder = "Вы пытаетесь завершить не тот заказ"
