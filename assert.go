package integrationtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/tokopedia/orderapp/acceptancetest/acceptance/config"
)

//TestCase holds test case object
type TestCase struct {
	DBConfig       *DBConfig
	Name           string
	Path           string
	FileName       string
	Pass           int
	Fail           int
	Total          int
	MessageSuccess []string
	MessageFail    []string
}

type DBConfig struct {
	Name string
	User string
	Pass string
	Host string
	Port string
}

//CreateTestCases with path string
func CreateTestCases(name string, path string, dbCfg *DBConfig) []*TestCase {
	testCases := make([]*TestCase, 0)
	reqFiles, err := ioutil.ReadDir(path + "req/")
	if err != nil {
		log.Fatal(err)
	}
	for _, reqFile := range reqFiles {
		testCases = append(testCases, &TestCase{
			Name:     name,
			Path:     path,
			FileName: reqFile.Name(),
		})
	}
	return testCases
}

//ParsePostReq for testcase
func (t *TestCase) ParsePostReq() interface{} {
	return parseJSONPost(t.Path + "req/" + t.FileName)
}

//ParseGetReq for testcase
func (t *TestCase) ParseGetReq() string {
	result := "?"
	req := parseJSONGet(t.Path + "req/" + t.FileName)
	if len(req) == 0 {
		return ""
	}
	for k, v := range req {
		result += k + "=" + v + "&"
	}
	return result[:len(result)-1]
}

//ClearDB wipeout existing rows for tables that will be asserted
func (t *TestCase) ClearDB() {
	err := clearDB(t.Path+"db/"+t.FileName, t.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

//SeedDB will seed for each test files if exist
func (t *TestCase) SeedDB() {
	err := seedDB(t.Path+"seed/"+strings.Replace(t.FileName, "json", "sql", -1), t.DBConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

//AssertDB asserts expected db result
func (t *TestCase) AssertDB() {
	pass := false
	tableFail := ""
	var err error
	for n := 0; n <= config.MaxDBRetry; n++ {
		pass, tableFail, err = assertDB(t.Path+"db/"+t.FileName, t.DBConfig)
		if err != nil {
			log.Fatalln(err)
		}
		if pass {
			break
		}
	}
	t.Total++
	if pass {
		t.Pass++
		t.MessageSuccess = append(t.MessageSuccess, fmt.Sprintf("[%s] Assert DB Passed : %s", t.Name, t.FileName))
	} else {
		t.Fail++
		t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Assert DB Failed : %s", t.Name, t.FileName))
		msgs, err := getDBErrorMessage(t, tableFail)
		if err != nil {
			log.Fatalln(err)
		}
		for _, msg := range msgs {
			t.MessageFail = append(t.MessageFail, msg)
		}
	}
}

//AssertResp asserts expected api response
func (t *TestCase) AssertResp(resp *http.Response, onlyCompareThisField ...string) {
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	respJSON := parseJSONPost(t.Path + "resp/" + t.FileName)
	for _, field := range onlyCompareThisField {
		pass := reflect.DeepEqual(result[field], respJSON[field])
		t.Total++
		if pass {
			t.Pass++
			t.MessageSuccess = append(t.MessageSuccess, fmt.Sprintf("[%s] Assert resp Passed : %s", t.Name, t.FileName))
		} else {
			t.Fail++
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Assert resp Failed : %s", t.Name, t.FileName))
			resultMarshal, _ := json.Marshal(result[field])
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Actual : %v", t.Name, string(resultMarshal)))
			respMarshal, _ := json.Marshal(respJSON[field])
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Expected : %v", t.Name, string(respMarshal)))
		}
	}
	if len(onlyCompareThisField) == 0 {
		pass := reflect.DeepEqual(result, respJSON)
		t.Total++
		if pass {
			t.Pass++
			t.MessageSuccess = append(t.MessageSuccess, fmt.Sprintf("[%s] Assert resp Passed : %s", t.Name, t.FileName))
		} else {
			t.Fail++
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Assert resp Failed : %s", t.Name, t.FileName))
			resultMarshal, _ := json.Marshal(result)
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Actual : %v", t.Name, string(resultMarshal)))
			respMarshal, _ := json.Marshal(respJSON)
			t.MessageFail = append(t.MessageFail, fmt.Sprintf("[%s] Expected : %v", t.Name, string(respMarshal)))
		}
	}
}

//PrintTestResult and calculate sum pass and fail
func PrintTestResult(testCases []*TestCase) (int, int, []string) {
	var sumPass, sumFail int
	var msgFail []string
	for _, t := range testCases {
		sumPass += t.Pass
		sumFail += t.Fail
		for _, msg := range t.MessageSuccess {
			log.Println(msg)
		}
		for _, msg := range t.MessageFail {
			log.Println(msg)
			msgFail = append(msgFail, msg)
		}
	}
	return sumPass, sumFail, msgFail
}
