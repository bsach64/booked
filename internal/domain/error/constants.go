package errordom

type ErrorCategory string
type ErrorCategoryCode string

const (
	CATEGORY_USER   ErrorCategory = "USER"
	CATEGORY_DB     ErrorCategory = "DB"
	CATEGORY_SYSTEM ErrorCategory = "SYS"
	CATEGORY_EVENT  ErrorCategory = "EVENT"
)

const (
	// DB
	DB_READ_ERROR  ErrorCategoryCode = "01"
	DB_WRITE_ERROR ErrorCategoryCode = "02"

	// User
	INVALID_USER_ROLE ErrorCategoryCode = "01"
	EMTPY_PASSWORD    ErrorCategoryCode = "02"
	USER_NOT_FOUND    ErrorCategoryCode = "03"
	INVALID_PASSWORD  ErrorCategoryCode = "04"
	JWT_SIGN_FAILURE  ErrorCategoryCode = "05"
	INVALID_TOKEN     ErrorCategoryCode = "06"
	EMPTY_EMAIL       ErrorCategoryCode = "07"
	EMPTY_NAME        ErrorCategoryCode = "08"

	// System
	JSON_DECODE_ERROR ErrorCategoryCode = "01"
	UNKNOWN_ERROR     ErrorCategoryCode = "02"

	// Event
	INVALID_SEAT_COUNT ErrorCategoryCode = "01"
)

var Errors map[ErrorCategory]map[ErrorCategoryCode]string = map[ErrorCategory]map[ErrorCategoryCode]string{
	CATEGORY_DB: {
		DB_READ_ERROR:  "unable to read from db",
		DB_WRITE_ERROR: "unable to write to db",
	},

	CATEGORY_USER: {
		INVALID_USER_ROLE: "invalid user role",
		EMTPY_PASSWORD:    "empty password",
		USER_NOT_FOUND:    "user not found",
		INVALID_PASSWORD:  "invalid password",
		JWT_SIGN_FAILURE:  "could not sign token",
		INVALID_TOKEN:     "invalid token",
		EMPTY_EMAIL:       "empty email",
		EMPTY_NAME:        "empty name",
	},

	CATEGORY_SYSTEM: {
		JSON_DECODE_ERROR: "could not decode json",
		UNKNOWN_ERROR:     "unknown error occured",
	},

	CATEGORY_EVENT: {
		INVALID_SEAT_COUNT: "invalid seat count",
	},
}
