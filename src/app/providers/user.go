package providers

import(
	"github.com/astaxie/beego"
	"github.com/ausrasul/redisorm"
	"time"
)

var usertable = "users"

// User object
type User struct {
	Email string
	Name string
	Expires string
}

// Check if the user exists in the database and it is not expired.
func AuthenticateUser(email string) bool{
	//users := make(map[string]User)
	beego.Debug("Getting list of users")
	users, err := GetUsers()
	if err != nil {
		beego.Debug("Cannot get users from DB")
		return false
	}
	beego.Debug("Looking for existing user")
	user, ok := users[email]
	if !ok {
		beego.Debug("User not exist >" + email + "<, creating it...")
		users[email] = User{Email: email, Expires: "2016-12-31"}
		beego.Debug("Saving the user to db")
		err = SaveUsers(users)
		if err != nil{
			beego.Debug(err)
			return false
		}
	}
	beego.Debug("Reading the user we  just added")
	user, ok = users[email]
	if !ok {
		beego.Debug("User not exist >" + email + "<")
		return false
	}
	beego.Debug("checking the user account expiry")
	expT, err := time.Parse("2006-01-02", user.Expires)
	if err != nil || expT.Before(time.Now()) {
		beego.Debug("User expired ", err)
		return false
	}
	beego.Debug("All good, user authenticated.")
	return true
}


// Get list of users from database
func GetUsers() (map[string]User, error){
	users := make(map[string]User)
	err := redisorm.Get(usertable, &users)
	if err != nil {
		beego.Debug("Cannot find users in db ", usertable, " ", err)
		return users, err
	}
	return users, err
}

func SaveUsers(users map[string]User) error{
	err := redisorm.Set(usertable, users)
	if err != nil {
		beego.Debug("Cannot save users into db" , usertable, " ", err)
		return err
	}
	return err
}
