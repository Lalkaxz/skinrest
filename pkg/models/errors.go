package models

type AppError struct {
	Name    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUserNotFound         = &AppError{"UserNotFound", "This user does not exist"}
	ErrSkinNotFound         = &AppError{"SkinNotFound", "This skin does not exists"}
	ErrAlrRegistered        = &AppError{"AlrRegistered", "This user is already registered"}
	ErrInvalidToken         = &AppError{"InvalidToken", "Invalid token"}
	ErrInvalidTokenFormat   = &AppError{"InvalidTokenFormat", "Invalid token format"}
	ErrInvalidTokenClaims   = &AppError{"InvalidTokenClaims", "Invalid token claims"}
	ErrTokenExpired         = &AppError{"TokenExpired", "Token has expired"}
	ErrInvalidSigningMethod = &AppError{"InvalidSigningMethod", "Unexpected signing method"}
	ErrTokenNotProvided     = &AppError{"TokenNotProvided", "Token is not provided"}
	ErrInvalidSkinType      = &AppError{"InvalidSkinType", "Invalid skin type"}
	ErrInvalidIdFormat      = &AppError{"InvalidIdFormat", "Invalid ID format"}
)
