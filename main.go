package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Color constants
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var tokenCache = struct {
	Token  string
	Expiry time.Time
}{}

type TokenResponse struct {
	Token string `json:"token"`
}

type OrderResponse struct {
	Data OrderData `json:"data"`
	Err  bool      `json:"err"`
}

type OrderData struct {
	ID      int         `json:"id"`
	Date    string      `json:"date"`
	Status  Status      `json:"status"`
	Address Address     `json:"address"`
	Items   []OrderItem `json:"items"`
	Total   float64     `json:"total"`
	VAT     float64     `json:"vat"`
}

type Status struct {
	Name string `json:"name"`
}

type Address struct {
	HouseNo string `json:"house_no"`
	Street  string `json:"street"`
	City    string `json:"city"`
}

type OrderItem struct {
	Item  Item    `json:"item"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

type Item struct {
	Name string `json:"name"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func getToken() (string, error) {
	if tokenCache.Token != "" && time.Now().Before(tokenCache.Expiry) {
		return tokenCache.Token, nil
	}

	data := map[string]string{
		"email": os.Getenv("API_USERNAME"),
		"key":   os.Getenv("API_KEY"),
	}
	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", os.Getenv("TOKEN_URL"), bytes.NewBuffer(jsonData))
	req.Header.Set(os.Getenv("HEADER_KEY"), os.Getenv("HEADER_VALUE"))
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	tokenCache.Token = tokenResp.Token
	tokenCache.Expiry = time.Now().Add(24 * time.Hour)
	return tokenCache.Token, nil
}

func fetchOrder(token, orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf(os.Getenv("ORDER_URL"), orderID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set(os.Getenv("HEADER_KEY"), os.Getenv("HEADER_VALUE"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, err
	}

	return &orderResp, nil
}

func formatMoney(amount float64) string {
	// Convert the float to a string with 2 decimal places
	amountStr := fmt.Sprintf("%.2f", amount)

	parts := strings.Split(amountStr, ".")

	intPart := parts[0]
	decimalPart := parts[1]

	// Add commas to the integer part
	var result strings.Builder
	for i, n := range reverse(intPart) {
		if i != 0 && i%3 == 0 {
			result.WriteString(",")
		}
		result.WriteByte(byte(n))
	}

	formattedAmount := reverse(result.String()) + "." + decimalPart
	return formattedAmount
}

// Helper function to reverse a string
func reverse(s string) string {
	var result strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		result.WriteByte(s[i])
	}
	return result.String()
}

func loadTemplate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func printReceipt(orderData *OrderResponse) {
	if orderData == nil || orderData.Err {
		fmt.Println(Red + "Error: Invalid order data" + Reset)
		return
	}

	data := orderData.Data
	storeName := os.Getenv("STORE_NAME")
	storeAddress := os.Getenv("STORE_ADDRESS")
	receipt := fmt.Sprintf(
		"\x1b\x61\x01"+ // Center alignment
			"\x1b\x45\x01%s\x1b\x45\x00"+
			"%s\n"+
			"======================\n\n"+
			"\x1b\x61\x00"+ // Left alignment
			"Order ID: #%d\n"+
			"Date: %s\n"+
			"Status: %s\n\n"+
			"Shipping Address:\n"+
			"%s %s\n"+
			"%s\n\n"+
			"Item                          Qty  Amount\n"+
			"---------------------------------------------\n",
		storeName, storeAddress, data.ID, data.Date, data.Status.Name, data.Address.HouseNo, data.Address.Street, data.Address.City,
	)

	for _, item := range data.Items {
		name := item.Item.Name
		if len(name) > 30 {
			name = name[:30]
		}
		receipt += fmt.Sprintf("%-30s %-3d %10s\n", name, item.Qty, formatMoney(item.Price*float64(item.Qty)))
	}

	footerTemplate, err := loadTemplate("footer_template.txt")
	if err != nil {
		fmt.Println(Red+"Error loading footer template:"+Reset, err)
		return
	}

	footerTemplate = strings.ReplaceAll(footerTemplate, "{{LINK}}", os.Getenv("STORE_LINK"))

	receipt += fmt.Sprintf(
		"---------------------------------------------\n"+
			"%-30s %10s\n"+
			"\x1b\x61\x01"+ // Center alignment for footer
			footerTemplate+
			"\x1b\x61\x00"+ // Left alignment
			"\x1d\x56\x42\x00", // Cut command
		"Total:", formatMoney(data.Total),
	)

	printToDefaultPrinter(receipt)
	fmt.Println(Green + receipt + Reset)
}

func main() {
	for {
		fmt.Print(Cyan + "Enter the order ID: " + Reset)
		var orderID string
		fmt.Scan(&orderID)

		token, err := getToken()
		if err != nil {
			fmt.Println(Red+"Error getting token:"+Reset, err)
			continue
		}

		orderData, err := fetchOrder(token, orderID)
		if err != nil {
			fmt.Println(Red+"Error fetching order:"+Reset, err)
			continue
		}

		printReceipt(orderData)
	}
}
