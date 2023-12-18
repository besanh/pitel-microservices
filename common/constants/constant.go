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

	// Match with MAP_NETWORK
	MAP_NETWORK_FPT = map[string]string{
		"viettel": "3",
		"vina":    "2",
		"mobi":    "1",
		"htc":     "5",
		"beeline": "4",
		"itel":    "6",
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

	// Status fpt
	STATUS_FPT = map[int]string{
		1001: "Request không hợp lệ",
		1002: "Client không được phép truy cập",
		1005: "Các scope không hợp lệ",
		1006: "Lỗi server",
		1007: "Server tạm thời không thể xử lý các request từ client",
		1008: "Thông tin client không đúng (sai client_id hoặc client_secret)",
		1010: "Scope không đủ quyền truy cập",
		1011: "Access token không hợp lệ",
		1013: "Access token hết hạn",
		1014: "Các tham số truyền vào bị lỗi",
		1015: "Không hỗ trợ kiểu loại hình cấp quyền này",
		1016: "Số lượng tin nhắn gửi đã vượt hạn mức",
		2501: "Tin nhắn trùng trong 5p",
		1:    "Tin nhắn trùng trong 5p",
		2502: "Hết hạn mức gửi tin",
		2503: "Chưa cấu hình hạn mức gửi tin",
		2504: "Brandname chưa kích hoạt hoặc bị khóa",
		54:   "Brandname chưa kích hoạt hoặc bị khóa",
		2505: "Số điện thoại bị chặn",
		-11:  "Số điện thoại bị chặn",
		2506: "Lỗi service",
		2507: "Lỗi không xác định",
		2:    "Brandname chưa được đăng kí với nhà mạng",
		-8:   "Brandname chưa được đăng kí với nhà mạng",
		3:    "Lỗi service của nhà mạng",
		4:    "Độ dài tin nhắn vượt quá qui định của nhà mạng",
		-14:  "Độ dài tin nhắn vượt quá qui định của nhà mạng",
		901:  "Độ dài tin nhắn vượt quá qui định của nhà mạng",
		5:    "Nội dung tin nhắn (template) chưa được đăng kí hoặc gửi sai so với template đã đăng kí",
		-20:  "Nội dung tin nhắn (template) chưa được đăng kí hoặc gửi sai so với template đã đăng kí",
		55:   "Nội dung tin nhắn (template) chưa được đăng kí hoặc gửi sai so với template đã đăng kí",
		6:    "Nội dung gửi có chứa từ khóa bị chặn",
		-18:  "Nội dung gửi có chứa từ khóa bị chặn",
		7:    "Nội dung chứ kí tự tiếng việt (Unicode) khi mã hóa (Hướng Viettel Bank)",
		8:    "Không thể giải mã, tin gửi không được mã hóa .. (Hướng Viettel Bank)",
		53:   "Sai số điện thoại",
		-10:  "Sai số điện thoại",
		902:  "Sai số điện thoại",
	}

	// External plugin connect
	EXTERNAL_PLUGIN_CONNECT_TYPE = []string{"incom", "abenla", "fpt"}
	SCOPE_FPT                    = []string{"send_brandname_otp", "send_brandname"}
)
