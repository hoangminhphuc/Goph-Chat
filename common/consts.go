package common

import "time"

/*
	! Plugin prefix name consts
*/
const (
	PluginDBMain = "mysql"
	PluginWSMain = "websocket"
	PluginHTTPMain = "gin"
	PluginRedisMain = "redis"
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

/* 
	! TTL consts
*/

const (
	// Session consts 
	SessionTTL = 2 * time.Hour

	// User
	UserProfileTTL = 4 * time.Hour

	// Room 
	RoomMetadataTTL = 30 * time.Minute
	RoomListPageTTL = 10 * time.Minute
)