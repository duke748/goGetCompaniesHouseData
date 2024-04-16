package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Company struct {
	Name    string
	Address string
}

// getCompaniesInfo retrieves information about companies and their addresses.
func getCompaniesInfo() {
	companies := GetCompanies()

	for i, company := range companies {
		address, err := getCompanyAddressAdvanced(company.Name)
		if err != nil {
			fmt.Printf("Error getting address for company %s: %v\n", company.Name, err)
			continue
		}

		companies[i].Address = address

		fmt.Printf("%s, %s\n", company.Name, address)
	}
}

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Uncomment the following line to get companies info from a list of companies you suitably define in the companiesData.go file
	// getCompaniesInfo()

	// create a list of company types, and for each type, get the address for the company
	companyTypes := []Company{
		{Name: "Architect", Address: ""},
		{Name: "Architects", Address: ""},
		{Name: "Architecture", Address: ""},
	}

	for _, company := range companyTypes {

		address, err := getCompanyByAddress(company.Name, os.Getenv("AddressFilter"))
		if err != nil {
			fmt.Printf("Error getting address for company %s: %v\n", company.Name, err)
			continue
		}
		fmt.Printf("%s, %s\n", company.Name, address)
	}
}
