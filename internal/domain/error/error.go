package errordom

import "fmt"

type AppError struct {
	Category     ErrorCategory
	CategoryCode ErrorCategoryCode
	ErrorToWrap  error
	Msg          string
}

func (ae *AppError) Error() string {
	return fmt.Sprintf(
		"category=%v category_code=%v err=%v msg=%v",
		ae.Category,
		ae.CategoryCode,
		ae.ErrorToWrap.Error(),
		ae.Msg,
	)
}

func GetEventError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_EVENT,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}

func GetUserError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_USER,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}

func GetDBError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_DB,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}

func GetWaitlistError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_WAITLIST,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}

func GetTicketError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_TICKET,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}

func GetSystemError(categoryCode ErrorCategoryCode, msg string, err error) error {
	return &AppError{
		Category:     CATEGORY_SYSTEM,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
	}
}
