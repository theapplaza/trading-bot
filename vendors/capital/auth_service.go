package capital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type session struct {
	securitytoken    string
	cst              string
	StreamingHost    string
	CurrentAccountId string
}

var (
    activeSession *session
    authOnce      sync.Once
    authErr       error
)

func authenticate() error {
    authOnce.Do(func() {
        authErr = doAuthenticate()
    })
    return authErr
}

func doAuthenticate() (err error) {

	req, err := _setupReq()
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		err = fmt.Errorf("cannot authenticate %s %s", err, response.Status)
		return
	}
	defer response.Body.Close()

	activeSession, err = _readBody(response)

	return
}

func _setupReq() (req *http.Request, err error) {
	reqBody, _ := json.Marshal(map[string]string{
		"identifier": activeConfig.ApiKeyUser,
		"password":   activeConfig.ApiKeyPassword,
	})

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/session", activeConfig.ApiBaseUrl), bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-CAP-API-KEY", activeConfig.ApiKey)
	return
}

func _readBody(resp *http.Response) (sess *session, err error) {
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &sess)

	sess.securitytoken = resp.Header.Get("X-SECURITY-TOKEN")
	sess.cst = resp.Header.Get("CST")

	return

}
