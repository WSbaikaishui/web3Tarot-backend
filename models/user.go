package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v4"
	"time"
)

const userTable = "users"

type User struct {
	Model
	ID          uint        `gorm:"column:id;type:int(11) unsigned NOT NULL AUTO_INCREMENT;PRIMARY_KEY"`
	Name        null.String `gorm:"column:name;type:VARCHAR(255)"`
	Address     string      `gorm:"column:address;type:VARCHAR(255);NOT NULL;DEFAULT:'';unique"`
	SeedMessage string      `gorm:"column:seed_message;type:VARCHAR(255);NOT NULL;DEFAULT:''"` // will store secret code only to save storage
	PublicKey   null.String `gorm:"column:public_key;type:VARCHAR(255)"`
	KeyStore    null.String `gorm:"column:key_store;type:VARCHAR(1023)"`
	IsOnline    bool        `gorm:"column:is_online;type:bool;NOT NULL;DEFAULT:false"`

	CreatedAt time.Time `gorm:"column:created_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func CreateUser(address string, seedMessage string, Name string) error {
	isSuccess := db.Create(&User{
		Address:     address,
		SeedMessage: seedMessage,
		Name:        null.StringFrom(Name),
	}).Error
	if isSuccess != nil {
		return isSuccess
	}
	return nil
}

func GetUser(address string) (*User, bool, error) {
	user := new(User)
	if err := db.Where("address = ?", address).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}

func GetUsers(addresses []string) ([]*User, bool, error) {
	var users []*User
	if err := db.Where("address IN ?", addresses).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return users, true, nil
}

func SetKeyInfo(address string, publicKey string, keyStore string) error {
	isSuccess := db.Model(&User{}).Where("address = ?", address).Updates(map[string]interface{}{
		"public_key": publicKey,
		"key_store":  keyStore,
	}).Error
	if isSuccess != nil {
		return isSuccess
	}
	return nil
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}