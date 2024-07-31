package constant

type FileAction string

const (
	FileActionNoChange FileAction = "nochange"
	FileActionUpdate   FileAction = "update"
	FileActionDelete   FileAction = "delete"
)
