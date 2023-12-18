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

	EXPORT_KEY = "export_"
	EXPORT_DIR = "/root/go/src/exported/"

	PLUGIN = map[string]string{
		"internal": "INERTAL",
		"external": "EXTERNAL",
		"abenla":   "ABENLA",
		"incom":    "INCOM",
	}

	PLUGIN_OFFICAL = []string{
		"internal",
		"external",
		"abenla",
		"incom",
	}

	// Incom
	ROUTERULE = map[string]string{
		"1": "zns",
		"2": "autocall",
		"3": "sms",
	}

	// Network incom
	MAP_NETWORK = map[string]string{
		"0":  "0", // foreign
		"1":  "1", // mobifone: mobifone
		"2":  "2", // vinaphone: vinaphone
		"3":  "3", // viettel: viettel
		"11": "4", // gtel: gtel
		"12": "5", // vietnamobile: vietnamobile
		"14": "6", // i-telecom: i-telecom
	}

	NETWORKS = map[string]string{
		"0": "foreign",
		"1": "mobifone",
		"2": "vinaphone",
		"3": "viettel",
		"4": "gtel",
		"5": "vietnamobile",
		"6": "i-telecom",
	}

	// Status
	STATUS = map[string]string{
		// incom
		"success":    "success",
		"fail":       "fail",
		"processing": "processing",

		// abenla
		"sent_fail":          "sent_fail",
		"account_expired":    "account_expired",
		"wrong_phone_number": "wrong_phone_number",
		"amount_zero":        "amount_zero",
		"not_price":          "not_price",
		"can_not_sent":       "can_not_sent",
		"deny_phone_number":  "deny_phone_number",
		"wrong_sender_name":  "wrong_sender_name",
	}

	// External plugin connect
	EXTERNAL_PLUGIN_CONNECT_TYPE = []string{"incom", "abenla", "fpt"}
	SCOPE_FPT                    = []string{"send_brandname_otp", "send_brandname"}
)
