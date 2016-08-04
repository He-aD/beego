package user

import (
	"github.com/astaxie/beego/validation"
)

type adresse struct {
	Nom string 			`valid:"Required;MinSize(2);MaxSize(30)"`
	Prenom string		`valid:"MinSize(2);MaxSize(30)"`
	NumeroRue string	`valid:"Required;Numeric;MaxSize(5);"`
	Voie string			`valid:"Required;MinSize(4);MaxSize(100);"`
	CodePostal string	`valid:"Required;AlphaDash;MinSize(1);MaxSize(30);"`
	Ville string		`valid:"Required;MinSize(2);MaxSize(30);"`
}

type Adresse struct {
	adresse
}

func NewAdresse(no, p, nu, vo, c, vi string) (*Adresse) {
	return &Adresse{adresse{no, p, nu, vo, c, vi}}
}

func (this Adresse) Valid(v *validation.Validation) (bool, error){
	return v.Valid(this.adresse)
}

func (this Adresse) Validate (bool, error){
	
}

type contact struct {
	Nom string 			`valid:"Required;MinSize(2);MaxSize(30)"`
	Prenom string 		`valid:"MinSize(2);MaxSize(30)"`
	Email string 		`valid:"Required;Email;MinSize(7);MaxSize(60)"`
	Telephone string 	`valid:"Phone"`
}

type Contact struct {
	contact
}

func NewContact(n, p, e, t string) (*Contact) {
	return &Contact{contact{n, p, e, t}}
}

func (this Contact) Valid(v *validation.Validation) (bool, error) {
	return v.Valid(this.contact)
}
