package models

type Customer struct {
	Id        uint
	FirstName string
	LastName  string
	Contact   int
	Type      string
}

type MenuItem struct {
	Name  string
	Id    uint
	Price float64
}

type Order struct {
	Id         uint
	CustomerId uint
	Amount     float64
	Status     string
}

type Bill struct {
	Id      uint
	OrderId uint
	Amount  float64
}
type HourDetails struct {
	Hour  int
	Table []Table
}
type Table struct {
	TableId int
	IsEmpty bool
}

type User struct {
	Id       uint
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Islogged bool
}

type UserRole int

const (
	Admin UserRole = iota
	Customers
	Guest
)

var RoleType = map[UserRole]string{
	Admin:     "Admin",
	Customers: "Customer",
	Guest:     "Guest",
}

type WeekDay int

const (
	Monday WeekDay = iota
	Tuesday
	Wednesday
	Thurday
	Friday
	Saturday
	Sunday
)

func (d WeekDay) String() string {
	return [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}[d]
}

type Book struct {
	WeekDay WeekDay `json:"weekday"`
	Time    int     `json:"time"`
}

type TableCount struct {
	Count int `json:"count"`
}
