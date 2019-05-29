package model

import (
	"time"
)

//リポジトリ
var Personrepo *PersonRepo

type Person struct {
	Id         int       `xorm:"id"`
	Company    string    `xorm:"company"`
	Email      string    `xorm:"email"`
	Name       string    `xorm:"name"`
	Honorific  string    `xorm:"honorific"`
	Post_h     string    `xorm:"post_h"`
	Post_l     string    `xorm:"post_l"`
	Prefecture string    `xorm:"prefecture"`
	Address_h  string    `xorm:"address_h"`
	Address_l  string    `xorm:"address_l"`
	Jinto      string    `xorm:"jinto"`
	Saibaru    string    `xorm:"saibaru"`
	Created    time.Time `xorm:"created"`
}

type PersonRepo struct {
}

func init() {
	Personrepo = &PersonRepo{}
}

func (p PersonRepo) Update(person *Person) (err error) {
	_, err = Engine.Where("id = ?", person.Id).Update(person)
	if err != nil {
		return err
	}

	return nil
}

func (p PersonRepo) FindByJinto(jinto string) *[]Person {
	var persons []Person

	_ = Engine.Where("jinto = ?", jinto).
		Desc("created").
		Find(&persons)

	return &persons
}
