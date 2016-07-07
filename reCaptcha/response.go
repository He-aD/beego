package reCaptcha

import (
	"github.com/astaxie/beego/httplib"
	"encoding/json"
	"errors"
)

const NBERRORS = 4 //flag the number of google's error codes
const APIURL = "https://www.google.com/recaptcha/api/siteverify"


type Response struct {
	success bool
	challenge_ts string
	hostname string
}

func (this Response) Success() bool {
	return this.success
}

func (this Response) Challenge_ts() string {
	return this.challenge_ts
}

func (this Response) Hostname() string {
	return this.hostname
}

/*
	send the captcha js response to the google's api
	format with the above struct the response and send
	eventual errors
	@captcha : js response
	@secretKey : the necessary secret key send by google for the website	
*/
func Get(captcha, secretKey string) (res *Response, err error) {
	req := httplib.Post(APIURL)
	req.Param("secret", secretKey)
	req.Param("response", captcha)
	rJson, err := req.String()
	if err != nil {
	    
	}

	var i interface{}
	b := []byte(rJson)
	err = json.Unmarshal(b, &i) //can't parse directly with the right struct because can contain attribute "error-code" "-" pb with golang attr
	if err != nil {
		
	}
	verify := i.(map[string]interface{})
	var errCodes string
	var resp Response
	res = &resp
	for k, v := range verify {
    	switch vv := v.(type) {
	    	case string:
	        	switch k {
	        		case "challenge_ts":
	        			res.challenge_ts = v.(string)
	        		
	        		case "hostname":
	        			res.hostname = v.(string)
	        	}
	    	case bool:
	    		res.success = v.(bool)

	    	case []interface{}:
		       for _, u := range vv {
		            switch u {
		            	case "missing-input-secret":
		            		errCodes += " The secret parameter is missing."
		            		
		            	case "invalid-input-secret":
		            		errCodes += " The secret parameter is invalid or malformed."
		            		
		            	case "missing-input-response":
		            		errCodes += " The response parameter is missing."
		            	
		            	case "invalid-input-response":
		            		errCodes += " The response parameter is invalid or malformed."
		            } 
		        }
		       err = errors.New(errCodes)
    	}
	}
	return res, err
}


