package main

import (
	"log"

	"github.com/felixgunawan/integrationtest"
	testcase "github.com/felixgunawan/integrationtest/examples/standalone/test/case"
)

func main() {
	var sumPass, sumFail int
	var msgFail []string
	listTestCase := []func() (int, int, []string){
		testcase.Index,
	}

	for _, tc := range listTestCase {
		pass, fail, msg := tc()
		sumPass += pass
		sumFail += fail
		msgFail = integrationtest.CombineArrayString(msgFail, msg)
	}

	//Calculation
	log.Println(" ")
	log.Println("=== TOTAL === ")
	log.Printf("Pass : %d", sumPass)
	log.Printf("Fail : %d", sumFail)
	if len(msgFail) > 0 {
		log.Println("=== FAILED TESTS === ")
		for _, msg := range msgFail {
			log.Println(msg)
		}
	}
	if sumFail > 0 {
		panic("Test failed!")
	}
}
