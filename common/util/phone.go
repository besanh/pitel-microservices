package util

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/tel4vn/fins-microservices/common/log"
)

var PrefixMobile = map[string]string{
	"016966": "03966",
	"0169":   "039",
	"0168":   "038",
	"0167":   "037",
	"0166":   "036",
	"0165":   "035",
	"0164":   "034",
	"0163":   "033",
	"0162":   "032",
	"0120":   "070",
	"0121":   "079",
	"0122":   "077",
	"0126":   "076",
	"0128":   "078",
	"0123":   "083",
	"0124":   "084",
	"0125":   "085",
	"0127":   "081",
	"0129":   "082",
	"0199":   "059",
	"0186":   "056",
	"0188":   "058",
	"16966":  "03966",
	"169":    "039",
	"168":    "038",
	"167":    "037",
	"166":    "036",
	"165":    "035",
	"164":    "034",
	"163":    "033",
	"162":    "032",
	"120":    "070",
	"121":    "079",
	"122":    "077",
	"126":    "076",
	"128":    "078",
	"123":    "083",
	"124":    "084",
	"125":    "085",
	"127":    "081",
	"129":    "082",
	"1992":   "059",
	"1993":   "059",
	"1998":   "059",
	"1999":   "059",
	"186":    "056",
	"188":    "058",
}

var PrefixHome = map[string]string{
	"076":  "0296",
	"064":  "0254",
	"0281": "0209",
	"0240": "0204",
	"0781": "0291",
	"0241": "0222",
	"075":  "0275",
	"056":  "0256",
	"0650": "0274",
	"0651": "0271",
	"062":  "0252",
	"0780": "0290",
	"0710": "0292",
	"026":  "0206",
	"0511": "0236",
	"0500": "0262",
	"0501": "0261",
	"0230": "0215",
	"061":  "0251",
	"067":  "0277",
	"059":  "0269",
	"0351": "0226",
	"04":   "024",
	"039":  "0239",
	"0320": "0220",
	"031":  "0225",
	"0711": "0293",
	"08":   "028",
	"0321": "0221",
	"058":  "0258",
	"077":  "0297",
	"060":  "0260",
	"0231": "0213",
	"063":  "0263",
	"025":  "0205",
	"020":  "0214",
	"072":  "0272",
	"0350": "0228",
	"038":  "0238",
	"030":  "0229",
	"068":  "0259",
	"057":  "0257",
	"052":  "0232",
	"0510": "0235",
	"055":  "0255",
	"033":  "0203",
	"053":  "0233",
	"079":  "0299",
	"022":  "0212",
	"066":  "0276",
	"036":  "0227",
	"0280": "0208",
	"037":  "0237",
	"054":  "0234",
	"073":  "0273",
	"074":  "0294",
	"027":  "0207",
	"070":  "0270",
	"029":  "0216",
}

func ParsePhoneNumber(phoneNum string) string {
	// if phone have only 1 character -> fail
	for _, r := range phoneNum {
		if unicode.IsLetter(r) {
			return ""
		}
	}
	r, err := regexp.Compile("^(\\+84|84|084)")
	if err != nil {
		return ""
	}
	phoneNum = r.ReplaceAllString(phoneNum, "0")
	if len(phoneNum) < 10 {
		phoneNum = "0" + phoneNum
	}
	for prefix, v := range PrefixMobile {
		if strings.Contains(phoneNum, prefix) {
			p := `^(` + prefix + `)`
			re, err := regexp.Compile(p)
			if err != nil {
				log.Error(err)
			}
			phoneNum = re.ReplaceAllString(phoneNum, v)
			return phoneNum
		}
	}
	// for prefix, v := range PrefixHome {
	// 	if strings.Contains(phoneNum, prefix) {
	// 		p := `^(` + prefix + `)`
	// 		re, err := regexp.Compile(p)
	// 		if err != nil {
	// 			log.Error(err)
	// 		}
	// 		phoneNum = re.ReplaceAllString(phoneNum, v)
	// 		return phoneNum
	// 	}
	// }
	return phoneNum
}

func ParseTelToTelStr(phoneNumber string) string {
	return strings.Join(strings.Split(phoneNumber, ""), ".")
}

func HandleNetwork(phone string) string {
	regexMobi, _ := regexp.Compile(`^((070|079|077|076|078|090|093|089)\d{7})$`)
	regexViettel, _ := regexp.Compile(`^((03|097|098|086|096)\d{7,8})$`)
	regexVina, _ := regexp.Compile(`^((091|094|081|082|083|084|085|088)\d{7})$`)
	regexOffNet, _ := regexp.Compile(`^((092|052|056|058|059|087|099|095|087|055)\d{7})$`)
	regexTel, _ := regexp.Compile(`^((02)\d{9})$`)
	if regexMobi.MatchString(phone) {
		return "mobi"
	} else if regexViettel.MatchString(phone) {
		return "viettel"
	} else if regexVina.MatchString(phone) {
		return "vina"
	} else if regexOffNet.MatchString(phone) {
		return "offnet"
	} else if regexTel.MatchString(phone) {
		return "tel"
	}
	return "dnc"
}
