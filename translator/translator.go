package translator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const googleTranslateBaseUrl = "https://translate.googleapis.com/translate_a/single"

// Translate is used to convert textToConvert from language passed in sourceLang to targetLang
func Translate(textToConvert, sourceLang, targetLang string) (string, error) {

	googleTranslateFinalUrl := generateCompleteGoogleTranslateUrl(sourceLang, targetLang, textToConvert)

	resp, err := http.Get(googleTranslateFinalUrl)

	if err != nil {
		return "", fmt.Errorf("error while getting data from google translate : %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("error while calling google translate , error code : %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("error while reading response body : %s ", err)
	}

	var result []interface{}
	err = json.Unmarshal(respBody, &result)

	if err != nil {
		return "", fmt.Errorf("error unmarshalling google translate data : %s", err)
	}

	if len(result) > 0 {
		var translatedText strings.Builder
		innerResult := result[0]
		for _, iir := range innerResult.([]interface{}) {
			for _, tText := range iir.([]interface{}) {
				if s, ok := tText.(string); ok {
					translatedText.WriteString(s)
				}
				break
			}
		}
		return translatedText.String(), nil
	} else {
		return "", fmt.Errorf("no translated data in google translate response")
	}

}

func generateCompleteGoogleTranslateUrl(sourceLang, targetLang, textToConvert string) string {
	uriEncodedTextToConvert := convertTextToUriEncoded(textToConvert)
	return fmt.Sprintf("%s?client=gtx&sl=%s&tl=%s&dt=t&q=%s", googleTranslateBaseUrl, sourceLang, targetLang, uriEncodedTextToConvert)
}

func convertTextToUriEncoded(text string) string {
	return url.QueryEscape(text)
}
