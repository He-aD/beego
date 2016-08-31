package conf

import (
	"os"
	"strings"
)

//create a conf.js file with the same content as the app.conf in an array
func ToJS() (err error) {
	var content, s, sCom string // sCom => string comon conf
	var t, sN, sJs []string //sN split with '\n', sJs => key authorized for conf.js file
	var conf, file *os.File
	arrConf := make(map[string]string)
	b := make([]byte, 2000, 2000)
	
	defer func () {
		file.Close()
		conf.Close()
	} ()
	
	if conf, err = os.Open("conf/app.conf"); err != nil {
		return err
	}
	
	if _, err = conf.Read(b); err != nil {
		return err
	}	
	
	if _, err := os.Open("static/source/js/conf.js"); err != nil { //error in opening file
		if b := os.IsExist(err); b { //the file exist 
			if err = os.Remove("static/source/js/conf.js"); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	
	if file, err = os.Create("static/source/js/conf.js"); err != nil {
		return err
	}
	
	s = string(b)
	s = strings.Replace(s, " ", "", -1) //remove all white spaces
	sJs = strings.Split(s, "#js")
	sJs = strings.Split(sJs[1], "\n")
	t = strings.Split(s, "[dev]")
	sCom = t[0]
	sN = strings.Split(sCom, "\n") //comon key, value
	getKeyValue(sN, arrConf)
	
	switch arrConf["runmode"] {
		case "dev":
			sN = strings.Split(strings.Split(t[1], "[prod]")[0], "\n")
			break;
		
		case "prod":
			sN = strings.Split(strings.Split(t[1], "[prod]")[1], "\n")
			break;
			
		case "test":
			sN = strings.Split(strings.Split(t[1], "[test]")[1], "\n")
			break;
	}
	
	getKeyValue(sN, arrConf)
	content += "/*\n"
	content += " * main conf file for javascript frontend's files\n"
	content += " */\n\n"
	content += "conf = new Array();\n"
	sJs[len(sJs)-1] = strings.Replace(sJs[len(sJs)-1], "\x00", "", -1) //delete all end file char
	
	for k, v := range arrConf {
		for _, va := range sJs {
			sp := strings.Split(va, "/") 
			if strings.Replace(sp[0], "#", "", -1) == k { //this key is authorized for conf.js file
				if len(sp) < 2 {
					content += "conf[\"" + k + "\"] = \"" + v + "\";\n"
				} else { //if there is a slash the value has not to be initialized as a string
					content += "conf[\"" + k + "\"] = " + v + ";\n"
				}	
			}
		}
		
	}
	
	_, err = file.WriteString(content)
	
	return err
}

//conf is a conf part Split by '\n'
func getKeyValue (conf []string, array map[string]string) {
	for _, line := range conf {
		sE := strings.Split(line, "=")
		if len(sE) > 1 {
			for k, v := range sE { //conf key value 
				if k % 2 == 0 {
					array[v] = sE[k+1]
				}
			}
		}
		
	}
}