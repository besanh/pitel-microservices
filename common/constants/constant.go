package constants

const (
	// Version
	VERSION      = "1.0"
	VERSION_NAME = "Bussiness Support System(BSS)"
)

var (
	// Recipient
	RECIPIENT = []string{"international", "mobifone", "vinaphone", "viettel", "vietnammobile", "itel", "beeline", "reddi", "telnet", "offnet", "other"}

	// Channel
	CHANNEL = []string{"sms", "zns", "email", "autocall", "other"}

	// Role
	ROLE_ABELA = []string{"international"}
	ROLE_INCOM = []string{"mobifone", "vinaphone", "viettel", "vietnammobile", "itel", "reddi", "offnet"}
	ROLE_FPT   = []string{"international", "mobifone", "vinaphone", "viettel", "vietnammobile", "itel", "beeline", "reddi", "offnet"}
)
