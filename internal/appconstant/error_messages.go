package appconstant

const (
	// auth errors
	MsgAuthUserNotFound       = "user is not found"
	MsgAuthDuplicateUser      = "user with email %s is already registered"
	MsgAuthUnknownCredentials = "unknown credentials, please check your username and password"

	// repository errors
	MsgInsertError = "error inserting data"
	MsgQueryError  = "error querying data"
	MsgUpdateError = "error updating data"
	MsgDeleteError = "error deleting data"
)
