package util

import (
	"strconv"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGetTimeNowMetaFormat(t *testing.T) {
	// case 1
	result := GetTimeNowMetaFormat()

	year := result[0:4]
	if _, err := strconv.Atoi(year); err != nil {
		t.Errorf("%s is not intString", year)
	}
	month := result[5:7]
	if _, err := strconv.Atoi(month); err != nil {
		t.Errorf("%s is not intString", month)
	}
	day := result[8:10]
	if _, err := strconv.Atoi(day); err != nil {
		t.Errorf("%s is not intString", day)
	}
	isOkDateSep := slices.Compare([]string{"-", "-"}, []string{result[4:5], result[7:8]})
	if isOkDateSep != 0 {
		t.Error("isOkDateSep is not -", result[4:5], result[7:8])
	}

	hours := result[11:13]
	if _, err := strconv.Atoi(hours); err != nil {
		t.Errorf("%s is not intString", hours)
	}
	minutes := result[14:16]
	if _, err := strconv.Atoi(minutes); err != nil {
		t.Errorf("%s is not intString", minutes)
	}
	seconds := result[17:19]
	if _, err := strconv.Atoi(seconds); err != nil {
		t.Errorf("%s is not intString", seconds)
	}
	isOkTimeSep := slices.Compare([]string{":", ":"}, []string{result[13:14], result[16:17]})
	if isOkTimeSep != 0 {
		t.Error("isOkTimeSep is not -", result[13:14], result[16:17])
	}

	isSign := slices.Contains([]string{"+", "-"}, result[20:21])
	if !isSign {
		t.Errorf("%s is not sign", result[20:21])
	}
	timeDiff := result[21:]
	if _, err := strconv.Atoi(timeDiff); err != nil {
		t.Errorf("%s is not intString", timeDiff)
	}
}
