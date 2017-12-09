package twilio

import "net/http"

type Response struct {
	*http.Response
}

func (r *Response) IsOK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

/*
var data map[string]interface{}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return err
	}
 */
