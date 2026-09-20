package switchbotapi

import (
	"net/http"
	"strconv"
)

func headerReqSet(req *http.Request, token, sign, nonce string, time int64) {
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("sign", sign)
	req.Header.Set("t", strconv.FormatInt(time, 10))
	req.Header.Set("nonce", nonce)
}
