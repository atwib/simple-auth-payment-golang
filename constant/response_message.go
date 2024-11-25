package constant

const (

	// Authentication
	InvalidRequest                   = "Invalid request"
	ErrAuthFailedHashPassword        = "failed to hash password"
	ErrAuthFailedRegisterUser        = "Registration failed"
	ErrAuthFailedLogin               = "Login failed"
	ErrAuthCannotBeEmpty             = "username and password cannot be empty"
	ErrAuthInvalidRefreshToken       = "invalid refresh token"
	ErrAuthCouldNotParseToken        = "could not parse token claims"
	ErrAuthInvalidTokenClaims        = "invalid token claims"
	ErrAuthFailedGenerateTokens      = "failed to generate tokens"
	ErrAuthFailedUnauthorized        = "Authorization header missing"
	ErrAuthFailedUnauthorizedFormat  = "Invalid authorization format"
	ErrAuthInvalidOrExpiredToken     = "Invalid or expired token"
	ErrAuthInvalidUsernameOrPassword = "invalid username or password"
	ErrAuthInvalidOldPassword        = "invalid old password"
	ErrAuthFailedResetPassword       = "failed to reset password"
	ErrAuthFailedSelf                = "failed to get user profile"
	ErrAuthFailedBlacklistToken      = "failed to blacklist token"
	ErrAuthTokenIsBlacklisted        = "token is blacklisted"

	SucAuthRegister      = "Successfully registered"
	SucAuthLogin         = "Login successfully"
	SucAuthRefreshToken  = "Refresh token successfully"
	SucAuthSelf          = "Get user profile successfully"
	SucAuthResetPassword = "Successfully reset password"
	SucAuthLogout        = "Logout successfully"

	// Wallet

	ErrWalletMinTopUp            = "Top up amount must be greater than _MINTOPUP_"
	ErrWalletInsufficientBalance = "insufficient balance"
	ErrWalletAlreadyExists       = "wallet already exists, you can't add more wallet"
	ErrWalletPinInvalid          = "pin invalid"
	ErrWalletPinMustBe4Digit     = "pin must be 4 digit"
	ErrWalletPinWrong            = "pin wrong"
	ErrWalletFailedHashPin       = "failed to hash pin"
	ErrWalletType                = "Type must be Customer or Merchant"
	ErrWalletPinConfirmNotMatch  = "pin confirm not match"
	ErrWalletFailedRegister      = "Registration failed"

	SucWalletRegister        = "Successfully registered"
	SucWalletTopUp           = "Successfully top up"
	SucWalletTransfer        = "Successfully transfer"
	SucWalletTransferConfirm = "Are you sure to transfer this amount?"
	SucWalletWithdraw        = "Successfully withdraw"
)
