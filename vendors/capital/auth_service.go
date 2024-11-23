package capital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
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

func authenticate(force bool) error {

	if force {
		return doAuthenticate()		
	}
	
    authOnce.Do(func() {
        authErr = doAuthenticate()
    })
    return authErr
}

func doAuthenticate() (err error) {

	req := _setupReq()
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("cannot authenticate %s %s", err, err)
		
		//check if it is timeout
		if strings.Contains(err.Error(), "timeout") {
			time.Sleep(30 * time.Second)
			return doAuthenticate()
		}
		return
	}

	if response == nil {
		err = fmt.Errorf("cannot authenticate %s %s", err, err)
		return
	}

	if response.StatusCode != http.StatusOK {
		
	}
	defer response.Body.Close()

	activeSession, err = _readBody(response)

	return
}

func _setupReq() (req *http.Request) {
	reqBody, _ := json.Marshal(map[string]string{
		"identifier": activeConfig.ApiKeyUser,
		"password":   activeConfig.ApiKeyPassword,
	})

	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/session", activeConfig.ApiBaseUrl), bytes.NewBuffer(reqBody))
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
