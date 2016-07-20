package reCaptcha

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "encoding/json"
)


type mockGGResponse struct {
	Success bool			`json:"success"`
	Challenge_ts string		`json:"challenge_ts"`
	Hostname string			`json:"hostname"`
	Err []string			`json:"error-codes"`
}

func (this mockGGResponse) testNResponse(r *Response, e error, t *testing.T) {
	if this.Success != r.Success {
		t.Error("Entrée : %v.l'attribut success n'est pas cohérent %v au lieu de %v", this, r.Success, this.Success)
	}
	if this.Challenge_ts != r.Challenge_ts {
		t.Error("Entrée : %vl'attribut challenge_ts n'est pas cohérent %v au lieu de %v", r.Challenge_ts, this.Challenge_ts)
	}
	if this.Hostname != r.Hostname {
		t.Error("Entrée : %vl'attribut hostname n'est pas cohérent %v au lieu de %v", r.Hostname, this.Hostname)
	}
	
	if this.Err == nil && e != nil || this.Err != nil && e == nil {
		t.Error("Entrée : %v. la fonction a retourné une captcha erreur qui n'existe pas")
	} else if this.Err != nil && e != nil{
		find := false
		er := e.(*ErrorCaptcha)
		for _, v := range this.Err {
			if v == er.Code() {
				find = true
				break
			}
		}
		if !find {
			t.Error("Entrée : %v. la fonction a retourné une mauvaise captcha erreur %v au lieu de %v", e, this.Err)
		}
	}
}

type wrongResponse struct {
	V1 int							
	V2 map[string]int
}

func TestGet(t *testing.T) {
	captcha, secret := "titi", "toto"
	var resp mockGGResponse
	var result *Response
	var err error
	var b []byte
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //créer le serveur qui simule les réponses du captcha de google
		r.ParseForm()
		if r.FormValue("response") != captcha {
			t.Error("la fonction n'envoie pas le captcha passé en paramètre au service de google.")
		}
		if r.FormValue("secret") != secret {
			t.Error("la fonction n'envoie pas la secret key passé en paramètre au service de google.")
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}))
	defer ts.Close()
	
	//test les cas nominaux qui ne doient pas retourner d'erreur
	sl := []mockGGResponse{
		mockGGResponse{true, "titi", "toto", nil},
		mockGGResponse{false, "titi", "toto", nil},
		mockGGResponse{false, "", "toto", nil},
		mockGGResponse{true, "titi", "", nil},
		mockGGResponse{true, "titi", "toto", []string{"missing-input-secret", "invalid-input-secret", "missing-input-response", "invalid-input-response"}},
		mockGGResponse{false, "titi", "toto", []string{"invalid-input-secret"}},
		mockGGResponse{false, "", "toto", []string{"missing-input-response"}},
		mockGGResponse{true, "titi", "", []string{"invalid-input-response"}},
	}
	for _, v := range sl { 
		resp = v
		b, err = json.Marshal(resp)
		if err != nil {
			t.Fatal(err.Error())
		}	
		result, err = Get(captcha, secret, ts.URL)
		if _, ok := err.(*ErrorCaptcha); err != nil && !ok {
			t.Fatal("Erreur retournée par la fonction testée. " + err.Error())
		}
		resp.testNResponse(result, err, t)
	}
	
	//test les cas non nominaux qui doient retourné une erreur à chaque fois
	resp = mockGGResponse{true, "titi", "toto", []string{"pr", "invalid-input-secret", "missing-input-response", "invalid-input-response"}}
	b, err = json.Marshal(resp)
	if err != nil {
		t.Fatal(err.Error())
	}	
	result, err = Get(captcha, secret, ts.URL)
	if err == nil {
		t.Fatal("La fonction n'a pas retourné d'erreur ! contexte : entrée : %v, réponse : %v.", resp, result)
	}
	wresp := wrongResponse{1, map[string]int{"tt": 1, "ll": 2}}
	b, err = json.Marshal(wresp)
	if err != nil {
		t.Fatal(err.Error())
	}	
	result, err = Get(captcha, secret, ts.URL)
	if err == nil {
		t.Fatal("La fonction n'a pas retourné d'erreur ! contexte : entrée : %v, réponse : %v.", wresp, result)
	}
}