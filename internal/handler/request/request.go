package request

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SaveSummaryRequest は保存リクエストの構造体
type SaveSummaryRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
	Content     string `json:"content" validate:"required"`
	UserID      string `json:"user_id"`
}

// ListSummaryRequest はリスト取得リクエストの構造体
type ListSummaryRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// DetailSummaryRequest は詳細取得リクエストの構造体
type DetailSummaryRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteSummaryRequest は削除リクエストの構造体
type DeleteSummaryRequest struct {
	ID string `json:"id" validate:"required"`
}

// SaveUserRequest は保存リクエストの構造体
type SaveUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,max=255"`
	UserType string `json:"user_type" validate:"required"`
}

// ListUserRequest はリスト取得リクエストの構造体
type ListUserRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// DetailUserRequest は詳細取得リクエストの構造体
type DetailUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteUserRequest は削除リクエストの構造体
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}

// Validate はリフレクションを使用してvalidateタグに基づいたバリデーションを行います
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		validateTag := field.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			if err := validateField(field.Name, value, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, rule string) error {
	switch {
	case rule == "required":
		if value.Kind() == reflect.String && strings.TrimSpace(value.String()) == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
		if value.Kind() == reflect.Int && value.Int() == 0 {
			return fmt.Errorf("%s is required", fieldName)
		}
	case strings.HasPrefix(rule, "max="):
		maxStr := strings.TrimPrefix(rule, "max=")
		max, err := strconv.Atoi(maxStr)
		if err != nil {
			return fmt.Errorf("invalid max value for %s", fieldName)
		}
		if value.Kind() == reflect.String && len(value.String()) > max {
			return fmt.Errorf("%s must be less than or equal to %d characters", fieldName, max)
		}
		if value.Kind() == reflect.Int && value.Int() > int64(max) {
			return fmt.Errorf("%s must be less than or equal to %d", fieldName, max)
		}
	case strings.HasPrefix(rule, "min="):
		minStr := strings.TrimPrefix(rule, "min=")
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return fmt.Errorf("invalid min value for %s", fieldName)
		}
		if value.Kind() == reflect.String && len(value.String()) < min {
			return fmt.Errorf("%s must be at least %d characters", fieldName, min)
		}
		if value.Kind() == reflect.Int && value.Int() > 0 && value.Int() < int64(min) {
			return fmt.Errorf("%s must be at least %d", fieldName, min)
		}
	}

	return nil
}
