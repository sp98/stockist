package instrument

import "os"

const (
	marketCloseTime      = "%s 15:30:00"
	marketActualOpenTime = "%s 09:13:00 MST"
	tstringFormat        = "2006-01-02 15:04:05"
	layOut               = "2006-01-02 15:04:05"
	influxLayout         = "2006-01-02T15:04:05Z"
)

var (
	accessToken = os.Getenv("ACCESSTOKEN")
)

var data = [][]string{
	//Instrument Name, Sybmol, Token, Exchange, Interval
	{"YES BANK", "YESBANK", "3050241", "NSE", "5m"},
	{"INDIABULLS HSG FIN", "IBULHSGFIN", "7712001", "NSE", "5m"},
	{"ZEE ENTERTAINMENT ENT", "ZEEL", "975873", "NSE", "5m"},
	{"BANK OF BARODA", "BANKBARODA", "1195009", "NSE", "5m"},
	{"SHRIRAM TRANSPORT FIN CO.", "SRTRANSFIN", "1102337", "NSE", "5m"},
	{"MOTHERSON SUMI SYSTEMS LT", "MOTHERSUMI", "1076225", "NSE", "5m"},
	{"CADILA HEALTHCARE", "CADILAHC", "2029825", "NSE", "5m"},
	{"ICICI PRU LIFE INS CO", "ICICIPRULI", "4774913", "NSE", "5m"},
	{"ICICI LOMBARD GIC", "ICICIGI", "5573121", "NSE", "5m"},
	{"ADANI PORT & SEZ", "ADANIPORTS", "3861249", "NSE", "5m"},
	{"JSW STEEL", "JSWSTEEL", "3001089", "NSE", "5m"},
	{"VEDANTA", "VEDL", "784129", "NSE", "5m"},
	{"RELIANCE POWER", "RPOWER", "3906305", "NSE", "5m"},
	{"DISH TV INDIA", "DISHTV", "3721473", "NSE", "5m"},
	{"SUZLON ENERGY", "SUZLON", "3076609", "NSE", "5m"},
	{"BHARTI INFRATEL", "INFRATEL", "7458561", "NSE", "5m"},
	{"DLF", "DLF", "3771393", "NSE", "5m"},
	{"NTPC", "NTPC", "2977281", "NSE", "5m"},
	{"BHARAT PETROLEUM CORP  LT", "BPCL", "134657", "NSE", "5m"},
	{"TATA STEEL", "TATASTEEL", "895745", "NSE", "5m"},
	{"INDUSIND BANK", "INDUSINDBK", "1346049", "NSE", "5m"},
	{"GAIL (INDIA)", "GAIL", "1207553", "NSE", "5m"},
	{"STEEL AUTHORITY OF INDIA", "SAIL", "758529", "NSE", "5m"},
	{"L&T FINANCE HOLDINGS", "L&TFH", "6386689", "NSE", "5m"},
	{"HINDUSTAN PETROLEUM CORP", "HINDPETRO", "359937", "NSE", "5m"},
	{"AUROBINDO PHARMA", "AUROPHARMA", "70401", "NSE", "5m"},
	{"BANDHANBNK", "BANDHANBNK", "579329", "NSE", "5m"},
	{"DIVI'S LABORATORIES", "DIVISLAB", "2800641", "NSE", "5m"},
	{"BIOCON", "BIOCON", "2911489", "NSE", "5m"},
	{"INTERGLOBE AVIATION", "INDIGO", "2865921", "NSE", "5m"},
	{"UNITED SPIRITS", "MCDOWELL-N", "2674433", "NSE", "5m"},
	{"LUPIN", "LUPIN", "2672641", "NSE", "5m"},
	{"ICICI BANK", "ICICIBANK", "1270529", "NSE", "5m"},
	{"AXIS BANK", "AXISBANK", "1510401", "NSE", "5m"},
	{"STATE BANK OF INDIA", "SBIN", "779521", "NSE", "5m"},
	{"RELIANCE INDUSTRIES", "RELIANCE", "738561", "NSE", "5m"},
	{"SBI LIFE INSURANCE CO", "SBILIFE", "5582849", "NSE", "5m"},
	{"DABUR INDIA", "DABUR", "197633", "NSE", "5m"},
	{"INDIAN OIL CORP", "IOC", "415745", "NSE", "5m"},
	{"DEWAN HOUSING FIN CORP LT", "DHFL", "215553", "NSE", "5m"},
	{"OIL AND NATURAL GAS CORP.", "ONGC", "633601", "NSE", "5m"},
	{"BHARTI AIRTEL", "BHARTIARTL", "2714625", "NSE", "5m"},
	{"HDFC STAND LIFE IN CO", "HDFCLIFE", "119553", "NSE", "5m"},
}
