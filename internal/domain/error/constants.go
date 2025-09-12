package errordom

type ErrorCategory string
type ErrorCategoryCode string

const (
	CATEGORY_USER   ErrorCategory = "USER"
	CATEGORY_DB     ErrorCategory = "DB"
	CATEGORY_SYSTEM ErrorCategory = "SYS"
)

const (
	// DB
	DB_READ_ERROR  ErrorCategoryCode = "01"
	DB_WRITE_ERROR ErrorCategoryCode = "02"

	// User
	INVALID_USER_ROLE ErrorCategoryCode = "01"

	// System
	JSON_DECODE_ERROR ErrorCategoryCode = "01"
	UNKNOWN_ERROR     ErrorCategoryCode = "02"
)

var Errors map[ErrorCategory]map[ErrorCategoryCode]string = map[ErrorCategory]map[ErrorCategoryCode]string{
	CATEGORY_DB: {
		DB_READ_ERROR:  "unable to read from db",
		DB_WRITE_ERROR: "unable to write to db",
	},

	CATEGORY_USER: {
		INVALID_USER_ROLE: "invalid user role",
	},

	CATEGORY_SYSTEM: {
		JSON_DECODE_ERROR: "could not decode json",
		UNKNOWN_ERROR:     "unknown error occured",
	},
}
