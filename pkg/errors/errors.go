package errors

import (
	"context"
	"errors"

	"github.com/o-ga09/web-ya-hime/pkg/logger"
)

const (
	callStack = "callStack"
)

type ErrType string

type ErrCode string

var (
	ErrCodeUnAuthorized    ErrCode = "unauthorized"     // 401
	ErrCodeUnAuthorization ErrCode = "unauthorization"  // 403
	ErrCodeInValidArgument ErrCode = "invalid argument" // 400
	ErrCodeBussiness       ErrCode = "business error"   // 400
	ErrCodeConflict        ErrCode = "conflict"         // 409
	ErrCodeNotFound        ErrCode = "not found"        // 404
	ErrCodeCritical        ErrCode = "critical error"   // 500
)

var (
	ErrTypeUnAuthorized    ErrType = "unauthorized"
	ErrTypeUnAuthorization ErrType = "unauthorization"
	ErrTypeBussiness       ErrType = "business error"
	ErrTypeConflict        ErrType = "conflict"
	ErrTypeNotFound        ErrType = "not found"
	ErrTypeCritical        ErrType = "critical error"
)

var (
	// ドメインエラー
	ErrInvalidFirebaseID  = errors.New("不正なFirebaseIDです。")
	ErrInvalidUserID      = errors.New("不正なUserIDです。")
	ErrInvalidName        = errors.New("不正なユーザー名です。")
	ErrInvalidDisplayName = errors.New("不正な表示名です。")
	ErrInvalidGroupID     = errors.New("不正なグループIDです。")
	ErrInvalidRelationID  = errors.New("不正なリレーションIDです。")
	ErrInvalidTwitterID   = errors.New("不正なTwitterIDです。")
	ErrInvalidGender      = errors.New("性別の値の範囲が不正です。")
	ErrInvalidDateTime    = errors.New("日付のフォーマットが不正です。")
	ErrInvalidProfileURL  = errors.New("不正なプロフィールURLです。")
	ErrInvalidUserType    = errors.New("無効なユーザータイプフォーマットです。")
	ErrFollowed           = errors.New("すでにフォロー済みです。")
	ErrFollowSelf         = errors.New("自分自身をフォローすることはできません。")
	ErrRequestNotNil      = errors.New("リクエストが正しくありません。")

	// ulidエラー
	ErrEmptyULID   = errors.New("empty ulid")
	ErrInvalidULID = errors.New("invalid ulid")

	// データベースエラー
	ErrRecordNotFound         = errors.New("record not found")
	ErrConflict               = errors.New("conflict")
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict")
	ErrForeignKeyConstraint   = errors.New("foreign key constraint error")
	ErrUniqueConstraint       = errors.New("unique constraint error")

	// 画像エラー
	ErrInvalidImageType  = errors.New("ファイルの種類が不正です。")
	ErrFailedImageName   = errors.New("ファイル名の生成に失敗しました。")
	ErrFailedDecodeImage = errors.New("画像のデコードに失敗しました。")
	ErrNotFoundImage     = errors.New("画像が見つかりません。")

	// リクエストエラー
	ErrRequestBodyNil = errors.New("リクエストボディが空です。")

	// その他エラー
	ErrSystem           = errors.New("システムエラーが発生しました。")
	ErrAuthorized       = errors.New("認証に失敗しました。")
	ErrUnauthorized     = errors.New("認可に失敗しました。")
	ErrInvalidArgument  = errors.New("バリデーションエラーが発生しました。")
	ErrInvalidOperation = errors.New("無効な操作です。")
	ErrNotFound         = errors.New("指定されたデータが見つかりません。")
)

// ginのcontextに認証エラーをセットして、ログ出力する
func MakeAuthorizationError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
}

// ginのcontextに認可エラーをセットして、ログ出力する
func MakeAuthorizedError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
}

// ginのcontextにシステムエラーをセットして、ログ出力する
func MakeSystemError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Error(ctx, errMessage, callStack, stack)
}

func MakeBusinessError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
}

func MakeConflictError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
}

func MakeNotFoundError(ctx context.Context, msg string) {
	var wrapped error
	stack := getCallstack(wrapped)
	errMessage := GetMessage(wrapped)
	logger.Warn(ctx, errMessage, callStack, stack)
}

// エラーをラップして、返す
// もしエラーがラップされていない場合は、システムエラーでラップして返す
func Wrap(ctx context.Context, err error) error {
	return err
}

func New(ctx context.Context, err string) error {
	return errors.New(err)
}

func IsWrapped(err error) bool {
	return errors.Is(err, ErrAuthorized)
}

func Is(err error, target error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, target) {
		return true
	}
	return false
}

func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	if IsWrapped(err) {
		return errors.Unwrap(err).Error()
	}
	return err.Error()
}

func GetCode(err error) ErrCode {
	if err == nil {
		return ""
	}

	return ""
}

func getCallstack(err error) string {
	if err == nil {
		return ""
	}

	return ""
}
