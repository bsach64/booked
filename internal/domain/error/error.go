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
		"Category=%v, CategoryCode=%v, msg=%v, err=%v",
		ae.Category, ae.CategoryCode, ae.ErrorToWrap.Error(), ae.Msg,
	)
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
