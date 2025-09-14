package errordom

type ErrorCategory string
type ErrorCategoryCode string

const (
	CATEGORY_USER     ErrorCategory = "USER"
	CATEGORY_DB       ErrorCategory = "DB"
	CATEGORY_SYSTEM   ErrorCategory = "SYS"
	CATEGORY_EVENT    ErrorCategory = "EVENT"
	CATEGORY_TICKET   ErrorCategory = "TICKET"
	CATEGORY_WAITLIST ErrorCategory = "WAITLIST"
)

const (
	// DB
	DB_READ_ERROR  ErrorCategoryCode = "01"
	DB_WRITE_ERROR ErrorCategoryCode = "02"
	DB_TX_ERROR    ErrorCategoryCode = "03"

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
	INVALID_UUID      ErrorCategoryCode = "03"

	// Event
	INVALID_SEAT_COUNT     ErrorCategoryCode = "01"
	INVALID_EVENT_ID       ErrorCategoryCode = "02"
	NO_EVENT_FOUND         ErrorCategoryCode = "03"
	INVALID_NEW_EVENT      ErrorCategoryCode = "04"
	CANT_REDUCE_SEAT_COUNT ErrorCategoryCode = "05"

	// Ticket
	TOO_FEW_TICKETS     ErrorCategoryCode = "01"
	TICKET_NOT_RESERVED ErrorCategoryCode = "02"
	NOT_YOUR_TICKET     ErrorCategoryCode = "03"

	// waitlist
	ALREADY_IN_WAILIST     ErrorCategoryCode = "01"
	INVALID_WAITLIST_COUNT ErrorCategoryCode = "02"
)
