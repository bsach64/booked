package errordom

import "fmt"

type AppError struct {
	Category     ErrorCategory     `json:"category"`
	CategoryCode ErrorCategoryCode `json:"category_code"`
	ErrorToWrap  error             `json:"-"`
	ErrorString  string            `json:"error"`
	Msg          string            `json:"msg"`
}

func (ae *AppError) Error() string {
	return fmt.Sprintf(
		"Category=%v, CategoryCode=%v, CodeMsg=%v, msg=%v, err=%v",
		ae.Category, ae.CategoryCode, Errors[ae.Category][ae.CategoryCode], ae.ErrorToWrap.Error(), ae.Msg,
	)
}

func GetEventError(categoryCode ErrorCategoryCode, msg string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return &AppError{
		Category:     CATEGORY_EVENT,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
		ErrorString:  errStr,
	}
}

func GetUserError(categoryCode ErrorCategoryCode, msg string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return &AppError{
		Category:     CATEGORY_USER,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
		ErrorString:  errStr,
	}
}

func GetDBError(categoryCode ErrorCategoryCode, msg string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return &AppError{
		Category:     CATEGORY_DB,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
		ErrorString:  errStr,
	}
}

func GetSystemError(categoryCode ErrorCategoryCode, msg string, err error) error {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return &AppError{
		Category:     CATEGORY_SYSTEM,
		CategoryCode: categoryCode,
		Msg:          msg,
		ErrorToWrap:  err,
		ErrorString:  errStr,
	}
}
