package testcases

import (
	"log"
	"net/http"

	"github.com/felixgunawan/integrationtest"
)

func Index() (int, int, []string) {
	testCases := integrationtest.CreateTestCases("Index", "./test/file/index/", nil)
	for _, t := range testCases {
		resp, err := http.Get("http://localhost:55001" + t.ParseGetReq())
		if err != nil {
			log.Fatalln(err)
		}
		t.AssertResp(resp)
	}
	return integrationtest.PrintTestResult(testCases)
}
