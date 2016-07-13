package reCaptcha

import (
	"github.com/astaxie/beego/httplib"
	"encoding/json"
	"errors"
)

const NBERRORS = 4 //flag the number of google's error codes
const APIURL = "https://www.google.com/recaptcha/api/siteverify"


type Response struct {
	Success bool
	Challenge_ts string
	Hostname string
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
	    return nil, err
	}

	var i interface{}
	b := []byte(rJson)
	err = json.Unmarshal(b, &i) //can't parse directly with the right struct because can contain attribute "error-code" "-" pb with golang attr
	if err != nil {
		return nil, err
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
	        			res.Challenge_ts = v.(string)
	        		
	        		case "hostname":
	        			res.Hostname = v.(string)
	        	}
	    	case bool:
	    		res.Success = v.(bool)

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


