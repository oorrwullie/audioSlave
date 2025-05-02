package credentials

import (
	"encoding/base64"
	"fmt"

	gokc "github.com/keybase/go-keychain"
)

type Credentials struct {
	ServiceName string
	username    string
	password    string
}

func New(serviceName string) (*Credentials, error) {
	u, p, err := getKeychainItem(serviceName)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		ServiceName: serviceName,
		username:    u,
		password:    p,
	}, nil
}

func (c *Credentials) SetUsername(u string) {
	c.username = u
}

func (c *Credentials) SetPassword(p string) {
	c.password = p
}

func (c *Credentials) Save() error {
	if err := c.validate(); err != nil {
		return err
	}

	// Try update first
	query := gokc.NewItem()
	query.SetSecClass(gokc.SecClassGenericPassword)
	query.SetService(c.ServiceName)
	query.SetAccount(c.username)

	// Data to update
	attrs := gokc.NewItem()
	attrs.SetSecClass(gokc.SecClassGenericPassword)
	attrs.SetService(c.ServiceName)
	attrs.SetAccount(c.username)
	attrs.SetLabel("Homebridge UI credentials for AudioSlave")
	attrs.SetData([]byte(c.password))
	attrs.SetAccessible(gokc.AccessibleWhenUnlocked)
	attrs.SetSynchronizable(gokc.SynchronizableNo)

	if err := gokc.UpdateItem(query, attrs); err == nil {
		return nil
	}

	// If update fails (not found), fallback to add
	return gokc.AddItem(attrs)
}

func (c *Credentials) GetUsername() string {
	return c.username
}

func (c *Credentials) GetPassword() string {
	return c.password
}

func (c *Credentials) validate() error {
	if c.username == "" || c.password == "" {
		return fmt.Errorf("username and password must be set")
	}
	return nil
}

func getKeychainItem(serviceName string) (username, password string, err error) {
	query := gokc.NewItem()
	query.SetSecClass(gokc.SecClassGenericPassword)
	query.SetService(serviceName)
	query.SetMatchLimit(gokc.MatchLimitOne)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)

	results, err := gokc.QueryItem(query)
	if err != nil {
		return "", "", err
	}
	if len(results) == 0 {
		return "", "", fmt.Errorf("no credentials found for service: %s", serviceName)
	}

	item := results[0]
	return item.Account, string(item.Data), nil
}

func GenerateRandomPassword(length int) (string, error) {
	bytes, err := gokc.RandBytes(length)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
