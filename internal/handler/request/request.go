package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

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

// Bind はHTTPリクエストから構造体にデータをバインドします
// JSONボディ、クエリパラメータ、URLパスパラメータをサポートします
// タグ: json, query, path
func Bind(r *http.Request, v interface{}) error {
	// JSONボディのバインド（Content-TypeがapplicationまたはPOST/PUT/PATCHメソッドの場合）
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}
			defer r.Body.Close()

			if len(body) > 0 {
				if err := bindJSON(body, v); err != nil {
					return err
				}
			}
		}
	}

	// クエリパラメータのバインド
	if err := bindQuery(r, v); err != nil {
		return err
	}

	// パスパラメータのバインド（カスタムコンテキストから取得）
	if err := bindPath(r, v); err != nil {
		return err
	}

	return nil
}

// bindJSON はJSONボディを構造体にバインドします（jsonタグを使用）
func bindJSON(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse JSON body: %w", err)
	}
	return nil
}

// bindQuery はクエリパラメータを構造体にバインドします（queryタグを使用）
func bindQuery(r *http.Request, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a struct or pointer to struct")
	}

	typ := val.Type()
	query := r.URL.Query()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// queryタグからフィールド名を取得
		queryTag := field.Tag.Get("query")
		if queryTag == "" || queryTag == "-" {
			continue
		}

		// queryタグからフィールド名を抽出（オプション部分を除去）
		fieldName := strings.Split(queryTag, ",")[0]
		queryValue := query.Get(fieldName)

		if queryValue == "" {
			continue
		}

		// フィールドに値を設定
		if err := setFieldValue(fieldValue, queryValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

// bindPath はパスパラメータを構造体にバインドします（pathタグを使用）
// Go 1.22以降の r.PathValue() メソッドを使用します
func bindPath(r *http.Request, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a struct or pointer to struct")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// pathタグからフィールド名を取得
		pathTag := field.Tag.Get("path")
		if pathTag == "" || pathTag == "-" {
			continue
		}

		fieldName := strings.Split(pathTag, ",")[0]
		pathValue := r.PathValue(fieldName)

		if pathValue == "" {
			continue
		}

		// フィールドに値を設定
		if err := setFieldValue(fieldValue, pathValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

// setFieldValue はリフレクションを使用してフィールドに値を設定します
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse int: %w", err)
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse uint: %w", err)
		}
		field.SetUint(uintValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot parse bool: %w", err)
		}
		field.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot parse float: %w", err)
		}
		field.SetFloat(floatValue)
	case reflect.Pointer:
		elemType := field.Type().Elem()
		newValue := reflect.New(elemType)
		if err := setFieldValue(newValue.Elem(), value); err != nil {
			return err
		}
		field.Set(newValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
