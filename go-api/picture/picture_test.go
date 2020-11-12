package picture_test

import (
	"testing"

	"github.com/peashoot/peapi/picture"
)

func TestFillWordsInto(t *testing.T) {
	err := picture.FillWordsIntoPic("/root/code/Peapi/go-api/picture/scripts/source/timg2.jpg", "/root/code/Peapi/go-api/picture/scripts/source/output.jpg", "我的中国♡", "/root/code/Peapi/go-api/picture/scripts/source/SimSun.ttf", 2, 4)
	if err != nil {
		t.Error(err.Error())
	}
}
