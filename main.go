package main

import (
	"flag"
	"fmt"

	"github.com/redtoad/exceldb"
)

func main() {
	importFile := flag.String("import", "./import.xlsx", "输入路径")
	dateFormat := flag.String("date", "01/02/06", "日期格式")
	db, err := exceldb.LoadFromExcel(*importFile, exceldb.InMemoryDb,
		exceldb.DateColum("Date", *dateFormat))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
	strftime("%m-%Y", "Date") as 'month-year',
		Employee,
		SUM("Hours worked")
	FROM data
	WHERE "Status" != "non billable"
	GROUP BY strftime("%m-%Y", "Date"), Employee;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var monthYear string
		var resource string
		var sum float64
		if err := rows.Scan(&monthYear, &resource, &sum); err != nil {
			panic(err)
		}
		fmt.Printf("%s -> %s -> %f\n", monthYear, resource, sum)
	}

}
