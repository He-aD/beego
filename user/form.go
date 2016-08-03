package user

import (
	
)

type Adresse struct {
	Nom string 			`valid:"Required;MinSize(2);MaxSize(30)"`
	Prenom string		`valid:"MinSize(2);MaxSize(30)"`
	NumeroRue string	`valid:"Required;Numeric;MaxSize(5);"`
	Voie string			`valid:"Required;MinSize(4);MaxSize(100);"`
	CodePostal string	`valid:"Required;AlphaDash;MinSize(1);MaxSize(30);"`
	Ville string		`valid:"Required;MinSize(2);MaxSize(30);"`
}

type Contact struct {
	Nom string 			`valid:"Required;MinSize(2);MaxSize(30)"`
	Prenom string 		`valid:"MinSize(2);MaxSize(30)"`
	Email string 		`valid:"Required;Email;MinSize(7);MaxSize(60)"`
	Telephone string 	`valid:"Phone"`
}
