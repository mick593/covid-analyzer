package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/covid/summary", covidSummary)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type caseObject struct {
	ConfirmDate    string
	No             string
	Age            int
	Gender         string
	GenderEn       string
	Nation         string
	NationEn       string
	Province       string
	ProvinceId     int
	District       string
	ProvinceEn     string
	StatQuarantine int
}

func covidSummary(c *gin.Context) {
	filepath := "./covid-cases.json"
	fmt.Printf("// reading file %s\n", filepath)
	file, err1 := ioutil.ReadFile(filepath)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", filepath)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}
	var cases []caseObject

	err2 := json.Unmarshal(file, &cases)
	if err2 != nil {
		fmt.Println("error:", err2)

		os.Exit(1)
	}

	provinceList := make(map[string]int)
	ageGroup := map[string]int{"0-30": 0, "31-60": 0, "61+": 0, "N/A": 0}

	for k := range cases {
		// increment the number of case in the current case's province
		if cases[k].Province != "" {
			provinceList[cases[k].Province]++

		}

		//increment age in the correct group
		age := cases[k].Age
		if age >= 0 && age <= 30 {
			ageGroup["0-30"]++
		} else if age >= 31 && age <= 60 {
			ageGroup["31-60"]++
		} else if age >= 61 {
			ageGroup["61+"]++
		} else {
			ageGroup["N/A"]++
		}
	}
	fmt.Println(len(cases))
	type responseObj struct {
		Province map[string]int
		AgeGroup map[string]int
	}
	response := responseObj{Province: provinceList, AgeGroup: ageGroup}
	c.JSON(200, response)
}
