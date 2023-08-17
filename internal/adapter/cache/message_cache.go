package cache

const message_prefix = "message:"

func (c *Cache) GetMessageCache(input string) string {
	// c.cache.Set()
	return message_prefix + "123"
}
