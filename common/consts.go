package common

/* 
	! Plugin prefix name consts
*/
const (
	PluginDBMain = "mysql"
)

/* 
	! Security consts
*/
const (
	// 32 bytes
	DefaultSaltLength = 32
	// Absolute maximum salt length to prevent resource exhaustion
	MaxSaltLength = 1024

	BCRYPT_COST = 12
)

/* 
	! Current user consts
*/
const (
	CurrentUser = "current_user"
)