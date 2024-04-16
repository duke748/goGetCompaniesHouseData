package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// CompaniesHouseResponseAdvanced represents the response structure for advanced company information from Companies House.
type CompaniesHouseResponseAdvanced struct {
	Items []struct {
		CompanyName             string `json:"company_name"`
		RegisteredOfficeAddress struct {
			AddressLine1 string `json:"address_line_1"`
			Locality     string `json:"locality"`
			PostalCode   string `json:"postal_code"`
			Country      string `json:"country"`
		} `json:"registered_office_address"`
	} `json:"items"`
}

// CompaniesHouseResponse represents the response structure returned by the Companies House API.
type CompaniesHouseResponse struct {
	Items []struct {
		Title   string `json:"title"`
		Address struct {
			AddressLine1 string `json:"address_line_1"`
			PostalCode   string `json:"postal_code"`
			Number       string `json:"premises"`
			Town         string `json:"locality"`
			County       string `json:"region"`
		} `json:"address"`
	} `json:"items"`
}

func getCompanyAddressAdvanced(companyName string) (string, error) {
	baseUrl := "https://api.company-information.service.gov.uk/advanced-search/companies"
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return "", err
	}

	q := req.URL.Query()
	q.Add("company_name_includes", companyName)
	q.Add("company_status", "active")
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(os.Getenv("BEARER_TOKEN"), "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error reading response body: %v", err)
		return "", err
	}

	var data CompaniesHouseResponseAdvanced
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Errorf("Error unmarshalling response body: %v", err)
		return "", err
	}
	// Uncomment below to print full JSON response
	// fmt.Println("Data: ", data)

	// Work out distance from your Lattitude & Longditude co-ordinates
	lat1, err := strconv.ParseFloat(os.Getenv("lat1"), 64)
	if err != nil {
		fmt.Println("Error parsing lat1: ", err)
		return "", err
	}
	lon1, err := strconv.ParseFloat(os.Getenv("lon1"), 64)
	if err != nil {
		fmt.Println("Error parsing lon1: ", err)
		return "", err
	}

	if len(data.Items) > 0 {
		details := ""

		for _, item := range data.Items {
			lat2, lon2, err := getLatLong(item.RegisteredOfficeAddress.PostalCode)
			if err != nil {
				fmt.Printf("Error getting lat/long for postcode %s: %v\n", item.RegisteredOfficeAddress.PostalCode, err)
				return "", err
			}

			// calculate distance
			distance := haversine(lat1, lon1, lat2, lon2)
			details += fmt.Sprintf("%s, %s, %s, %s, %f\n", item.CompanyName, item.RegisteredOfficeAddress.AddressLine1, item.RegisteredOfficeAddress.Locality, item.RegisteredOfficeAddress.PostalCode, distance)
		}
		return details, nil
		// return fmt.Sprintf("%s, %s, %s, %s, %s, %f", data.Items[0].Address.Number, data.Items[0].Address.AddressLine1, data.Items[0].Address.Town, data.Items[0].Address.County, data.Items[0].Address.PostalCode, distance), nil
	}

	return "Address not found", nil
}

func getCompanyByAddress(companyName string, companyAddress string) (string, error) {
	baseUrl := "https://api.company-information.service.gov.uk/advanced-search/companies"
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return "", err
	}

	q := req.URL.Query()
	q.Add("company_name_includes", companyName)
	q.Add("company_status", "active")
	q.Add("location", companyAddress)
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(os.Getenv("BEARER_TOKEN"), "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error reading response body: %v", err)
		return "", err
	}

	var data CompaniesHouseResponseAdvanced
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Errorf("Error unmarshalling response body: %v", err)
		return "", err
	}
	// Uncomment below to print full JSON response
	// fmt.Println("Data: ", data)

	// Work out distance from your Lattitude & Longditude co-ordinates
	lat1, err := strconv.ParseFloat(os.Getenv("lat1"), 64)
	if err != nil {
		fmt.Println("Error parsing lat1: ", err)
		return "", err
	}
	lon1, err := strconv.ParseFloat(os.Getenv("lon1"), 64)
	if err != nil {
		fmt.Println("Error parsing lon1: ", err)
		return "", err
	}

	if len(data.Items) > 0 {
		details := ""

		for _, item := range data.Items {
			lat2, lon2, err := getLatLong(item.RegisteredOfficeAddress.PostalCode)
			if err != nil {
				fmt.Printf("Error getting lat/long for postcode %s: %v\n", item.RegisteredOfficeAddress.PostalCode, err)
				return "", err
			}

			// calculate distance
			distance := haversine(lat1, lon1, lat2, lon2)
			details += fmt.Sprintf("%s, %s, %s, %s, %f\n", item.CompanyName, item.RegisteredOfficeAddress.AddressLine1, item.RegisteredOfficeAddress.Locality, item.RegisteredOfficeAddress.PostalCode, distance)
		}
		return details, nil
		// return fmt.Sprintf("%s, %s, %s, %s, %s, %f", data.Items[0].Address.Number, data.Items[0].Address.AddressLine1, data.Items[0].Address.Town, data.Items[0].Address.County, data.Items[0].Address.PostalCode, distance), nil
	}

	return "Address not found", nil
}

func getCompanyAddress(companyName string) (string, error) {
	baseUrl := "https://api.companieshouse.gov.uk/search/companies"
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("q", companyName)
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(os.Getenv("BEARER_TOKEN"), "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Error reading response body: %v", err)
		return "", err
	}

	var data CompaniesHouseResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Errorf("Error unmarshalling response body: %v", err)
		return "", err
	}

	// Work out distance from your Lattitude & Longditude co-ordinates
	lat1, err := strconv.ParseFloat(os.Getenv("lat1"), 64)
	if err != nil {
		fmt.Println("Error parsing lat1: ", err)
		return "", err
	}
	lon1, err := strconv.ParseFloat(os.Getenv("lon1"), 64)
	if err != nil {
		fmt.Println("Error parsing lon1: ", err)
		return "", err
	}

	lat2, lon2, err := getLatLong(data.Items[0].Address.PostalCode)
	if err != nil {
		fmt.Printf("Error getting lat/long for postcode %s: %v\n", data.Items[0].Address.PostalCode, err)
		return "", err
	}

	// calculate distance
	distance := haversine(lat1, lon1, lat2, lon2)

	if len(data.Items) > 0 {
		details := ""

		for _, item := range data.Items {
			details += fmt.Sprintf("%s, %s, %s, %s, %s, %s, %f\n", item.Title, item.Address.Number, item.Address.AddressLine1, item.Address.Town, item.Address.County, item.Address.PostalCode, distance)
		}
		return details, nil
		// return fmt.Sprintf("%s, %s, %s, %s, %s, %f", data.Items[0].Address.Number, data.Items[0].Address.AddressLine1, data.Items[0].Address.Town, data.Items[0].Address.County, data.Items[0].Address.PostalCode, distance), nil
	}

	return "Address not found", nil
}
