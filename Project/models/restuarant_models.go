package models

// Customer represents a customer in the restaurant system.
type Customer struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Contact   int    `json:"contact"`
	Type      string `json:"type"`
}

// MenuItem represents a menu item in the restaurant.
type MenuItem struct {
	Name  string  `json:"name"`
	Id    uint    `json:"id"`
	Price float64 `json:"price"`
}

// Order represents a customer's order in the restaurant.
type Order struct {
	Id         uint    `json:"id"`
	CustomerId uint    `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

// Bill represents the bill generated for an order.
type Bill struct {
	Id       uint     `json:"id"`
	OrderId  uint     `json:"order_id"`
	Amount   float64  `json:"amount"`
	Customer Customer `json:"customer"`
}

// HourDetails represents the details of each hour for table availability.
type HourDetails struct {
	Hour  int     `json:"hour"`
	Table []Table `json:"table"`
}

// Table represents a restaurant table.
type Table struct {
	ID       string `json:"id"`
	Capacity int    `json:"capacity"`
	IsEmpty  bool   `json:"is_empty"`
}

// Booking represents a booking made by a customer for a table.
type Booking struct {
	ID            string   `json:"id"`
	Customer      Customer `json:"customer"`
	WeekDay       WeekDay  `json:"weekday"`
	Time          int      `json:"time"`
	BookNow       bool     `json:"booknow"`
	Table         Table    `json:"table"`
	BookingStatus string   `json:"status"`
}

// User represents a system user (Admin, Customer, Guest).
type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Islogged bool   `json:"islogged"`
}

// UserRole defines roles for the users in the system.
type UserRole int

const (
	Admin UserRole = iota
	Customers
	Guest
)

// RoleType maps UserRole to string.
var RoleType = map[UserRole]string{
	Admin:     "Admin",
	Customers: "Customer",
	Guest:     "Guest",
}

// WeekDay represents the days of the week.
type WeekDay int

const (
	Monday WeekDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// WeekDayToString map to convert WeekDay to its string representation
var WeekDayToString = map[WeekDay]string{
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
	Sunday:    "Sunday",
}
var StringToWeekDay = map[string]WeekDay{
	"Monday":    Monday,
	"Tuesday":   Tuesday,
	"Wednesday": Wednesday,
	"Thursday":  Thursday,
	"Friday":    Friday,
	"Saturday":  Saturday,
	"Sunday":    Sunday,
}

type WeekdayTime struct {
	WeekDay WeekDay
	Time    int
}

// String converts WeekDay to its string representation.
func (d WeekDay) String() string {
	return [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}[d]
}

// TableCount represents the count of tables in the system.
type TableCount struct {
	Count int `json:"count"`
}

// QueueEntry represents an entry in the queue for booking.
type QueueEntry struct {
	CustomerName string
	WeekDay      WeekDay
	Time         int
}
