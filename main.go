package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	routerURL = "http://192.168.1.1" //change
	rebootURL = "/newgui_adv_home.cgi?id=1422364133"
	authToken = "Basic YWRtaW46cGFzc3dvcmQ=" //admin:password
	debug = true
)

func main() {
	sendPayload()
	time.Sleep(60 * time.Second)
}

func sendPayload() {
	fmt.Println("Forming reboot request...")
	
	//add post data
	form := url.Values{}
	form.Add("buttonSelect", "2")
	form.Add("wantype", "dhcp")
	form.Add("enable_apmode", "0")

	req, err := http.NewRequest(http.MethodPost, routerURL+rebootURL, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("Error forming request:", err)
	}

	//add required headers
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	fmt.Println("Sending request...")
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println("Error in response:", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if strings.Contains(string(body), "Rebooting") != true {
		if debug == true {
			fmt.Printf("Got %d as status code\n", resp.StatusCode)
			fmt.Println("Response:", string(body))
			time.Sleep(3 * time.Second)
			sendPayload()
		} else {
			fmt.Println("Didn't get expected response, retrying in 3 seconds...")
			time.Sleep(3 * time.Second)
			sendPayload()
		}
	} else {
		fmt.Println("Reboot Success!")
	}
}
