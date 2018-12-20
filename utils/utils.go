package utils

import (
	"regexp"
	"strconv"
)

type Color struct {
	Name string
	Hex string
}

var Colors = []Color{
	Color{"black",   "#777777"},
	Color{"red",     "#F65A5A"},
	Color{"green",   "#00F100"},
	Color{"yellow",  "#EFEE04"},
	Color{"blue",    "#0F04E8"},
	Color{"cyan",    "#04E8E7"},
	Color{"magenta", "#F75AF6"},
	Color{"white",   "#FFFFFF"},
	Color{"gray",    "#7E7E7E"},
	Color{"brown",   "#6E3C3C"},
}

var matchAllColorCodes = regexp.MustCompile(`(\^\d)`)

var Ranks = map[int]string {
	 1: "Private First Class",
	 2: "Private First Class I",
	 3: "Private First Class II",
	 4: "Lance Corporal",
	 5: "Lance Corporal I",
	 6: "Lance Corporal II",
	 7: "Corporal",
	 8: "Corporal I",
	 9: "Corporal II",
	10: "Sergeant",
	11: "Sergeant I",
	12: "Sergeant II",
	13: "Staff Sergeant",
	14: "Staff Sergeant I",
	15: "Staff Sergeant II",
	16: "Gunnery Sergeant",
	17: "Gunnery Sergeant I",
	18: "Gunnery Sergeant II",
	19: "Master Sergeant",
	20: "Master Sergeant I",
	21: "Master Sergeant II",
	22: "Master Gunnery Sergeant",
	23: "Master Gunnery Sergeant I",
	24: "Master Gunnery Sergeant II",
	25: "Second Lieutenant",
	26: "Second Lieutenant I",
	27: "Second Lieutenant II",
	28: "First Lieutenant",
	29: "First Lieutenant I",
	30: "First Lieutenant II",
	31: "Captain",
	32: "Captain I",
	33: "Captain II",
	34: "Major",
	35: "Major I",
	36: "Major II",
	37: "Lieutenant Colonel",
	38: "Lieutenant Colonel I",
	39: "Lieutenant Colonel II",
	40: "Colonel",
	41: "Colonel I",
	42: "Colonel II",
	43: "Brigadier General",
	44: "Brigadier General I",
	45: "Brigadier General II",
	46: "Major General",
	47: "Major General I",
	48: "Major General II",
	49: "Lieutenant General",
	50: "Lieutenant General I",
	51: "Lieutenant General II",
	52: "General",
	53: "General I",
	54: "General II",
	55: "Commander",
}

func getRegexpsFor(i int) (*regexp.Regexp, *regexp.Regexp) {
	index := strconv.Itoa(i)

	replaceColorCode := `\^` + index + `(.*?)(\^[^`+ index + `]|$)`
	cleanup := `\^` + index

	return regexp.MustCompile(replaceColorCode), regexp.MustCompile(cleanup)
}

func Colorize(s string) string {
	result := s;

	for i, color := range Colors {
		colorRegexp, cleanupRegexp := getRegexpsFor(i)
		result = colorRegexp.ReplaceAllString(result, "<span style='color:" + color.Hex +";'>$1</span>$2")
		result = cleanupRegexp.ReplaceAllString(result, "")
	}

	return result;
}

func StripFormat(s string) string {
	result := matchAllColorCodes.ReplaceAllString(s, "")

	return result
}

func GetRankText(rank int) string {
	if rankText, ok := Ranks[rank]; ok {
		return rankText
	} else {
		return "Unknown"
	}
}