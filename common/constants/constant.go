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

	STANDARD_CODE = map[string]string{
		"1": "1", // success
		"2": "2", // error internal
		"3": "3", // fail
		"4": "4", // wrong phone
	}

	STANDARD_CODE_INCOM_TO_TEL4VN = map[string]string{
		"1":  "5",  // success
		"-1": "6",  // PhoneNumber Wrong Format, incom
		"-2": "7",  // wrong parameter,incom
		"-3": "8",  // out of quota,incom
		"-6": "9",  // wrong template code, incom
		"-8": "10", // wrong route rule, incom
		"-9": "11", // channel route rule wrong, incom
	}

	STANDARD_CODE_ABENLA_TO_TEL4VN = map[string]string{
		"100": "12", // other, abenla
		"101": "13", // user not exist, abenla
		"102": "14", // wron password, abenla
		"103": "15", // account deactivated, abenla
		"104": "16", // can not access, abenla
		"105": "17", // account is zero, abenla
		"106": "18", // success, abenla
		"107": "19", // wrong sign, abenla
		"108": "20", // wrong brandname, abenla
		"109": "21", // exceed sms, abenla
		"110": "22", // send sms fail, abenla
		"111": "23", // wrong service type, abenla
		"112": "24", // expire request, abenla
		"113": "25", // not has aes key, abenla
		"114": "26", // process fail, abenla
		"115": "27", // record not exist, abenla
		"116": "28", // record existed, abenla
		"117": "29", // blacklist keyword, abenla
		"118": "30", // wrong template, abenla
		"119": "31", // wrong caller or scriptId, abenla
		"120": "32", // caller and script id blocked, abenal
		"201": "33", // pending, abenla
		"203": "34", // sent success, abenla
		"204": "35", // sent fail, abenla
		"212": "36", // wait result, abenla
	}

	STANDARD_CODE_FPT_TO_TEL4VN = map[string]string{
		"1001": "37", // request khong hop le, fpt
		"1002": "38", // client khong duoc phep truy cap, fpt
		"1005": "39", // scope khong hop le, fpt
		"1006": "40", // loi server, fpt
		"1007": "41", // server tam thoi khong xl cac request tu client, fpt
		"1008": "42", // sai client_id hoac secret, fpt
		"1010": "43", // scope khong du quyen truy cap, fpt
		"1011": "44", // access token khong hop le, fpt
		"1013": "45", // access token het han, fpt
		"1014": "46", // cac tham so truyen vao khong hop le, fpt
		"1015": "47", // khong ho tro loai hinh cap quyen, fpt
		"1016": "48", // so tin nhan vuo qua han muc, fpt
		"2501": "49", // tn trung trong 5p, fpt
		"1":    "50", // tn trung trong 5p, fpt
		"2502": "51", // het han muc tn, fpt
		"2503": "52", // chua cau hinh han muc tn, fpt
		"2504": "53", // brandname chua kich hoat hoac bi khoa, fpt
		"54":   "54", // brandname chua kich hoat hoac bi khoa, fpt
		"2505": "55", // sdt bi chan, fpt
		"-11":  "56", // sdt bi chan, fpt
		"2506": "57", // loi khong xac dinh, fpt
		"2507": "58", // loi service, fpt
		"2":    "59", // brandname chua dk vs nha mang, fpt
		"-8":   "60", // brandname chua dk vs nha mang, fpt
		"3":    "61", // loi service nha mang, fpt
		"4":    "62", // do dai tn vuot qua QD cua nha mang, fpt
		"-14":  "63", // do dai tn vuot qua QD cua nha mang, fpt
		"901":  "64", // do dai tn vuot qua QD cua nha mang, fpt
		"5":    "65", // noi dung tn sai template hoac template chua dk, fpt
		"-20":  "66", // noi dung tn sai template hoac template chua dk, fpt
		"55":   "67", // noi dung tn sai template hoac template chua dk, fpt
		"6":    "68", // noi dung tn co keyword bi chan, fpt
		"-18":  "69", // noi dung tn co keyword bi chan, fpt
		"7":    "70", // noi dung tn co Unicode khi ma hoa(viettel), fpt
		"8":    "71", // khong the giai ma, tn khong duoc ma hoa(viettel), fpt
		"53":   "72", // sai sdt, fpt
		"-10":  "73", // sai sdt, fpt
		"902":  "74", // sai sdt, fpt
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
	MESSAGE_TEL4VN = map[string]string{
		"success":    "Success",
		"fail":       "Fail",
		"processing": "Processing",
	}

	ABENLA_CODE = map[string]string{
		"3": "3", // success
		"4": "4",
	}
	STATUS_ABENLA = map[string]string{
		"sent_fail":          "sent_fail",
		"account_expired":    "account_expired",
		"wrong_phone_number": "wrong_phone_number",
		"amount_zero":        "amount_zero",
		"not_price":          "not_price",
		"can_not_sent":       "can_not_sent",
		"deny_phone_number":  "deny_phone_number",
		"wrong_sender_name":  "wrong_sender_name",
	}

	// Message incom
	MESSAGE_INCOM = map[int]string{
		1:  "Success",
		-1: "Phone Number Wrong Format",
		-2: "Wrong Format Parameter",
		-3: "Out of message trial",
		-6: "Can't find template with template code",
		-8: "Wrong routerule",
		-9: "Channel XXXXX in routerule was wrong",
	}

	// Message abenla
	MESSAGE_ABENLA = map[int]string{
		100: "Other",
		101: "User not exist",
		102: "Wrong password",
		103: "Account deactivated",
		104: "Can not access",
		105: "Account is zero",
		106: "Success",
		107: "Wrong sign",
		108: "Wrong brandname",
		109: "Exceed sms",
		110: "Send sms fail",
		111: "Wrong service type",
		112: "Expire request",
		113: "Not has aes key",
		114: "Process fail",
		115: "Record not exist",
		116: "Record existed",
		117: "Blacklist keyword",
		118: "Wrong template",
		119: "Wrong caller or scriptId",
		120: "Caller and script id blocked",
		201: "Pending",
		203: "Sent success",
		204: "Sent fail",
		212: "Wait result",
	}

	// Message fpt
	MESSAGE_FPT = map[int]string{
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
