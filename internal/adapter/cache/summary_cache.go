package cache

import "fmt"

const summary_prefix = "summary:"

func (c *Cache) GetUserNames(root_prefix string) (usernames []string, err error) {
	prefix := root_prefix + summary_prefix
	entry, _ := c.cache.Get(prefix + "usernames")
	fmt.Println(string(entry))
	// c.cache.Set()
	return usernames, err
}
