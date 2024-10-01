package constants

type UserRole string

const (
	SuperAdmin UserRole = "SUPERADMIN"
	Admin      UserRole = "ADMIN"
	User       UserRole = "USER"
)
