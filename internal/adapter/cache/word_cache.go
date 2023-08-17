package cache

import "github.com/bytedance/sonic"

const (
	//domain
	domainPrefix = ":word"
	//key
	wordCloudPrefix   = ":wordcloud"
	filterCloudPrefix = ":filtercloud"
	stringRank        = ":stringrank"
	amountRank        = ":amountrank"
)

// WordCloud
func (c *Cache) SetWordCloudCache(UUID string, result *map[string]int) error {
	UUIDKey := UUID + domainPrefix + wordCloudPrefix
	value, err := sonic.Marshal(result)
	if err != nil {
		return err
	}
	err = c.cache.Set(UUIDKey, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) GetWordCloudCache(UUID string, result *map[string]int) error {
	UUIDKey := UUID + domainPrefix + wordCloudPrefix
	CacheValue, err := c.cache.Get(UUIDKey)
	if err != nil {
		return err
	}
	err = sonic.Unmarshal(CacheValue, &result)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) DeleteWordCloudCache(UUID string) error {
	UUIDKey := UUID + domainPrefix + wordCloudPrefix
	err := c.cache.Delete(UUIDKey)
	if err != nil {
		return err
	}
	return nil
}

// FilterCloud
func (c *Cache) SetFilterCloudCache(UUID string, result *map[string]int) error {
	UUIDKey := UUID + domainPrefix + filterCloudPrefix
	value, err := sonic.Marshal(result)
	if err != nil {
		return err
	}
	err = c.cache.Set(UUIDKey, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) GetFilterCloudCache(UUID string, result *map[string]int) error {
	UUIDKey := UUID + domainPrefix + filterCloudPrefix
	CacheValue, err := c.cache.Get(UUIDKey)
	if err != nil {
		return err
	}
	err = sonic.Unmarshal(CacheValue, &result)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) DeleteFilterCloudCache(UUID string) error {
	UUIDKey := UUID + domainPrefix + filterCloudPrefix
	err := c.cache.Delete(UUIDKey)
	if err != nil {
		return err
	}
	return nil
}

// stringRank
func (c *Cache) SetStringRankCache(UUID string, result *[]string) error {
	UUIDKey := UUID + domainPrefix + stringRank
	value, err := sonic.Marshal(result)
	if err != nil {
		return err
	}
	err = c.cache.Set(UUIDKey, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) GetStringRankCache(UUID string, result *[]string) error {
	UUIDKey := UUID + domainPrefix + stringRank
	CacheValue, err := c.cache.Get(UUIDKey)
	if err != nil {
		return err
	}
	err = sonic.Unmarshal(CacheValue, &result)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) DeleteStringRankCache(UUID string) error {
	UUIDKey := UUID + domainPrefix + stringRank
	err := c.cache.Delete(UUIDKey)
	if err != nil {
		return err
	}
	return nil
}

// amount rank
func (c *Cache) SetAmountRankCache(UUID string, result *[]int) error {
	UUIDKey := UUID + domainPrefix + stringRank
	value, err := sonic.Marshal(result)
	if err != nil {
		return err
	}
	err = c.cache.Set(UUIDKey, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) GetAmountRankCache(UUID string, result *[]int) error {
	UUIDKey := UUID + domainPrefix + stringRank
	CacheValue, err := c.cache.Get(UUIDKey)
	if err != nil {
		return err
	}
	err = sonic.Unmarshal(CacheValue, &result)
	if err != nil {
		return err
	}
	return nil
}
func (c *Cache) DeleteAmountRankCache(UUID string) error {
	UUIDKey := UUID + domainPrefix + stringRank
	err := c.cache.Delete(UUIDKey)
	if err != nil {
		return err
	}
	return nil
}
