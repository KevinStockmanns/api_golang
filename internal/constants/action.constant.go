package constants

type Action string

const (
	Create Action = "CREATE"
	Update Action = "UPDATE"
	Delete Action = "DELETE"
)

var Actions = map[Action]bool{
	Create: true,
	Update: true,
	Delete: true,
}
