package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"
	"github.com/valyala/fastjson"
	"sync"
	"bufio"
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
	var p fastjson.Parser
	var i int = 0

    scanner := bufio.NewScanner(r)

    for scanner.Scan() {
		v, err2 := p.Parse(string(scanner.Text()))
		if err2 != nil {
				return
		}

		user := userPool.Get().(*User)
		user.ID = v.GetInt("Id")
		user.Name = string(v.GetStringBytes("Name"))
		user.Username = string(v.GetStringBytes("Username"))
		user.Email = string(v.GetStringBytes("Email"))
		user.Phone = string(v.GetStringBytes("Phone"))
		user.Password = string(v.GetStringBytes("Password"))
		user.Address = string(v.GetStringBytes("Address"))

		result[i] = *user
		i++
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
