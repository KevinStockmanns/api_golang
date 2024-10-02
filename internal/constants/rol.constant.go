package constants

type UserRole string

const (
	SuperAdmin UserRole = "SUPERADMIN"
	Admin      UserRole = "ADMIN"
	User       UserRole = "USER"
)

var UserRoles = map[UserRole]bool{
	SuperAdmin: true,
	Admin:      true,
	User:       true,
}
