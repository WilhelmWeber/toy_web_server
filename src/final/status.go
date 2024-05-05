package final

import "net/http"

var STATUSES = map[int]string{
	http.StatusContinue:             "100 Continue",
	http.StatusSwitchingProtocols:   "101 Switching Protocols",
	http.StatusOK:                   "200 OK",
	http.StatusCreated:              "201 Created",
	http.StatusAccepted:             "202 Accepted",
	http.StatusNonAuthoritativeInfo: "203 Non-Authoritative Information",
	http.StatusNoContent:            "204 No Content",
	http.StatusResetContent:         "205 Reset Content",
	http.StatusPartialContent:       "206 Partial Content",
	http.StatusMultiStatus:          "207 Multi-Status",
	http.StatusAlreadyReported:      "208 Already Reported",
	http.StatusIMUsed:               "226 IM Used",
	http.StatusMultipleChoices:      "300 Multiple Choices",
	http.StatusMovedPermanently:     "301 Moved Permanently",
	http.StatusFound:                "302 Found",
	http.StatusSeeOther:             "303 See Other",
	http.StatusNotModified:          "305 Not Modified",
	http.StatusTemporaryRedirect:    "307 Temporary Redirect",
	http.StatusPermanentRedirect:    "308 Permanent Redirect",
	http.StatusBadRequest:           "400 Bad Request",
	http.StatusUnauthorized:         "401 Unauthorized",
	http.StatusPaymentRequired:      "402 Payment Required",
	http.StatusForbidden:            "403 Forbidden",
	http.StatusNotFound:             "404 not Found",
	http.StatusMethodNotAllowed:     "405 Method Not Allowed",
	http.StatusNotAcceptable:        "406 Not Acceptable",
}
