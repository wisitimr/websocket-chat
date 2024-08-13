package service

type Service struct {
	WebSocket     WsService
	Register      RegisterService
	Login         LoginService
	VerifyContact VerifyContactService
	ChatHistory   ChatHistoryService
	ContactList   ContactListService
}
