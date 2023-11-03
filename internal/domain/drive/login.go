package drive

import "net/url"

//add ?openExternalBrowser=1
func AppendOpenExternalBrowserParam(originalURL string) (string, error) {
	// 將提供的原始URL字串解析成*url.URL
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	// 取得URL的查詢參數
	query := parsedURL.Query()
	// 將 "openExternalBrowser" 參數添加到查詢參數中，值為 "1"
	query.Add("openExternalBrowser", "1")
	// 將更新後的查詢參數重新設置到parsedURL中
	parsedURL.RawQuery = query.Encode()
	// 返回修改後的URL
	return parsedURL.String(), nil
}
