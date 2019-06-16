package instrument

const (
	marketCloseTime      = "%s 15:30:00"
	marketActualOpenTime = "%s 09:13:00 MST"
	tstringFormat        = "2006-01-02 15:04:05"
	layOut               = "2006-01-02 15:04:05"
	influxLayout         = "2006-01-02T15:04:05Z"
	accessToken          = "TtSuNlGuMI4E44zNm0gYW3Q31RgNZ7Uj"
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
	{"INTERGLOBE AVIATION", "INDIGO", "2865921", "NSE", "5m"},
}
