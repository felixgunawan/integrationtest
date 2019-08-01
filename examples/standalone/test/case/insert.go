package testcases

import (
	"log"
	"net/http"

	"github.com/felixgunawan/integrationtest"
)

func Insert() (int, int, []string) {
	testCases := integrationtest.CreateTestCases("Insert", "./test/file/insert/", &integrationtest.DBConfig{
		Host: "localhost",
		Port: "5432",
		User: "user_example",
		Pass: "pass_example",
		Name: "db_example",
	})
	for _, t := range testCases {
		t.ClearDB()
		resp, err := http.Post("http://localhost:55001/add", "application/json", t.ParsePostReq())
		if err != nil {
			log.Fatalln(err)
		}
		t.AssertResp(resp)
		t.AssertDB()
	}
	return integrationtest.PrintTestResult(testCases)
}
