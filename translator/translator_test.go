package translator_test

import (
	"fmt"
	"testing"

	"github.com/godarkthrt/filenameconverter/translator"
)

func TestTranslate(t *testing.T) {
	expected := "नमस्ते !! आप कैसे हैं"
	got, _ := translator.Translate("Hello !! how are you", "en", "hi")

	fmt.Println("Translated result : ", got)
	if expected != got {
		t.Errorf("error while translating, expected :'%s' , got : '%s'", expected, got)
	}
}
