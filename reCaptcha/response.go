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

type ErrorCaptcha struct {
	code string
	description string
}

func (this *ErrorCaptcha) Error() string {
	return "code : "+this.code+". "+this.description
}

func (this ErrorCaptcha) Code() string {
	return this.code
}

func (this ErrorCaptcha) Description() string {
	return this.description
}

/*
	send the captcha js response to the google's api
	format with the above struct the response and send
	eventual errors
	@captcha : js response
	@secretKey : the necessary secret key send by google for the website	
*/
func Get(captcha, secretKey, apiURL string) (res *Response, err error) {
	if apiURL == "" {
		apiURL = APIURL
	}
	req := httplib.Post(apiURL)
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
	var errCodes *ErrorCaptcha
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
	        			
        			default: 
	        			return nil, errors.New("google a renvoyé une value string non supportée.")
	        	}
	    	case bool:
	    		res.Success = v.(bool)

	    	case []interface{}:
		       for _, u := range vv {
		            switch u {
		            	case "missing-input-secret":
		            		errCodes = &ErrorCaptcha{"missing-input-secret", "The secret parameter is missing."}
		            		
		            	case "invalid-input-secret":
		            		errCodes = &ErrorCaptcha{"invalid-input-secret", "The secret parameter is invalid or malformed."}
		            		
		            	case "missing-input-response":
		            		errCodes = &ErrorCaptcha{"missing-input-response", "The response parameter is missing."} 
		            	
		            	case "invalid-input-response":
		            		errCodes = &ErrorCaptcha{"invalid-input-response", "The response parameter is invalid or malformed."}
	            		
	            		default:		       
		            		return nil, errors.New("google a retourné un code erreur non supporté : " + u.(string))
		            } 
		        }
		       err = errCodes
	       
			case nil:
				//do nothing in this case. errors can be null in json				      
	       
	       default:
		       return nil, errors.New("google a retourné un type de value non supporté.")
    	}
	}
	return res, err
}


