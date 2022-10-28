package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"
	"github.com/valyala/fastjson"
	"bytes"
	"sync"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

var userPool = sync.Pool{
	New: func() interface{} {
		user := User{}
        return &user
	},
}

func getUsers(r io.Reader) (result users, err error) {
	var buf  bytes.Buffer
	_, err = io.Copy(&buf, r)

	if err != nil {
		return
	}

	lines := strings.Split(buf.String(), "\n")

	var p fastjson.Parser
	//var user User

	// pool.Put(user)

	for i, line := range lines {
		v, err2 := p.Parse(line)
		if err2 != nil {
				return
		}

		user := userPool.Get().(*User)

		// user = User{
		// 	v.GetInt("Id"),
		// 	string(v.GetStringBytes("Name")),
		// 	string(v.GetStringBytes("Username")),
		// 	string(v.GetStringBytes("Email")),
		// 	string(v.GetStringBytes("Phone")),
		// 	string(v.GetStringBytes("Password")),
		// 	string(v.GetStringBytes("Address")),
		// }
		user.ID = v.GetInt("Id")
		user.Name = string(v.GetStringBytes("Name"))
		user.Username = string(v.GetStringBytes("Username"))
		user.Email = string(v.GetStringBytes("Email"))
		user.Phone = string(v.GetStringBytes("Phone"))
		user.Password = string(v.GetStringBytes("Password"))
		user.Address = string(v.GetStringBytes("Address"))

		result[i] = *user
		userPool.Put(user)
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if user.Email == "" {
			continue
		}

		matched := strings.HasSuffix(user.Email, "." + domain)

		if matched {
			position := strings.LastIndex(user.Email, "@")
			value := strings.ToLower(user.Email[position+1:])

			result[value] += 1
		}
	}
	return result, nil
}
