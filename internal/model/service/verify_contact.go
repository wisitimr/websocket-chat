package service

type VerifyContactService interface {
	VerifyContact(username string) error
}
