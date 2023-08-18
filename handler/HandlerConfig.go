package handler

var userHandler HandlerI

func SetHandler(i HandlerI) {
	userHandler = i
}
func GetHandler() HandlerI {
	if userHandler == nil {
		return &DefaultHandler{}
	}
	return userHandler
}
