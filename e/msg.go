package e

var MsgFlags = map[int]string{
	SUCCESS:       "Operation successful",
	ERROR:         "Operation failed",
	InvalidParams: "Invalid request parameters",

	ErrorExistUser:    "User already exists",
	ErrorNotExistUser: "User does not exist",

	ErrorAuthCheckTokenFail:    "Token authentication failed",
	ErrorAuthCheckTokenTimeout: "Token has expired",
	ErrorAuthToken:             "Token generation failed",
	ErrorAuth:                  "Token error",
	ErrorNotCompare:            "Does not match",
	ErrorDatabase:              "Database operation error, please try again",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
