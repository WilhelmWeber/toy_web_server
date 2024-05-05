package toy

import "fmt"

type ToyReponse struct {
	HttpStatus  string
	ReponseBody string
}

func (r *ToyReponse) GetResponse(http_ver string, res_header string) (string, string) {
	res := fmt.Sprintf("%s %s\r\n%s\r\n%s", http_ver, r.HttpStatus, res_header, r.ReponseBody)
	return res, r.HttpStatus
}
