package roles

//goland:noinspection ALL
const (
	SUPER_ADMIN = 1
)

func LangRole(role int) string {
	switch role {
	case SUPER_ADMIN:
		return "super_admin"
	default:
		return "undefined"
	}
}
