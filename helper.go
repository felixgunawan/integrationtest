package integrationtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //using postgres
)

const (
	NoRowFlag = "NO_ROW"
)

func parseJSONPost(filePathAndName string) map[string]interface{} {
	jsonFile, err := os.Open(filePathAndName)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

func parseJSONGet(filePathAndName string) map[string]string {
	jsonFile, err := os.Open(filePathAndName)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)
	return result
}

type cronArgs struct {
	Args []string `json:"args"`
}

func parseCronArgsReq(cronName string, filePathAndName string) []string {
	result := make([]string, 0)
	jsonFile, err := os.Open(filePathAndName)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var argsCron cronArgs
	json.Unmarshal([]byte(byteValue), &argsCron)
	result = append(result, cronName)
	for _, args := range argsCron.Args {
		result = append(result, args)
	}
	return result
}

type AssertTable map[string]Table

type Table struct {
	Flags   []string                 `json:"flags"`
	Columns []map[string]interface{} `json:"columns"`
}

func assertDB(filePath string, dbCfg *DBConfig) (bool, string, error) {
	var assertTable AssertTable
	db, err := connectDb(dbCfg)
	if err != nil {
		return false, "", err
	}
	defer db.Close()
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, "", err
	}
	err = json.Unmarshal(jsonFile, &assertTable)
	if err != nil {
		return false, "", err
	}

	for tableName, table := range assertTable {
		noRow := IsInArrayString(table.Flags, NoRowFlag)
		for _, column := range table.Columns {
			query := generateQuery(tableName, column)
			var count int
			err = db.Get(&count, query)
			if err != nil {
				return false, tableName, err
			}
			if !noRow && count != 1 {
				return false, tableName, nil
			}
			if noRow && count > 0 {
				return false, tableName, nil
			}
		}
	}
	return true, "", nil
}

func getDBErrorMessage(t *TestCase, tableFail string) ([]string, error) {
	var assertTable AssertTable
	result := make([]string, 0)
	db, err := connectDb(t.DBConfig)
	if err != nil {
		return result, err
	}
	defer db.Close()
	jsonFile, err := ioutil.ReadFile(t.Path + "db/" + t.FileName)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(jsonFile, &assertTable)
	if err != nil {
		return result, err
	}

	for tableName, table := range assertTable {
		if tableName == tableFail {
			query := fmt.Sprintf("SELECT * FROM %s", tableName)
			rows, _ := db.Query(query)
			cols, _ := rows.Columns()
			dbResult := fmt.Sprintf("[%s] (%s) Actual : ", t.Name, tableName)
			rowCount := 0
			for rows.Next() {
				columns := make([]string, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i := range columns {
					columnPointers[i] = &columns[i]
				}

				rows.Scan(columnPointers...)

				for i, colName := range cols {
					dbResult += fmt.Sprintf("%s:%s ", colName, columns[i])
				}
				result = append(result, dbResult)
				rowCount++
			}
			if rowCount == 0 {
				dbResult += fmt.Sprintf("no data found")
				result = append(result, dbResult)
			}
			for _, col := range table.Columns {
				expected := fmt.Sprintf("[%s] (%s) Expected : ", t.Name, tableName)
				for k, v := range col {
					expected += fmt.Sprintf("%s:%s ", k, v)
				}
				result = append(result, expected)
			}
		}
	}
	return result, nil
}

func clearDB(filePath string, dbCfg *DBConfig) error {
	var x AssertTable
	db, err := connectDb(dbCfg)
	if err != nil {
		return err
	}
	defer db.Close()
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, &x)
	if err != nil {
		return err
	}

	for tableName := range x {
		query := generateDeleteQuery(tableName)
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedDB(filePath string, dbCfg *DBConfig) error {
	db, err := connectDb(dbCfg)
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Println("seedDB()", filePath)
	c, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		fmt.Println("seedDB(): IsNotExist")
		return nil
	}
	if err != nil {
		return err
	}
	sql := string(c)
	fmt.Println("seedDB(): execute: ", sql)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func generateDeleteQuery(tableName string) string {
	return fmt.Sprintf("DELETE FROM %s", tableName)
}

func generateQuery(tableName string, column map[string]interface{}) string {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE", tableName)
	i := 0
	for key, val := range column {
		colDef, ok := val.(map[string]interface{})
		if ok {
			query += fmt.Sprintf(` %s %s %s `, key, colDef["operator"], colDef["value"])
		} else {
			query += fmt.Sprintf(` %s = '%s' `, key, val)
		}

		if i < len(column)-1 {
			query += " AND "
		}
		i++
	}
	return query
}

func connectDb(dbCfg *DBConfig) (*sqlx.DB, error) {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Pass, dbCfg.Name)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
