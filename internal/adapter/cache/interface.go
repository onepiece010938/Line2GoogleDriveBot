package cache

type CacheI interface {
	GetMessageCache(input string) string

	SetWordCloudCache(UUID string, result *map[string]int) error
	SetFilterCloudCache(UUID string, result *map[string]int) error
	SetStringRankCache(UUID string, result *[]string) error
	SetAmountRankCache(UUID string, result *[]int) error

	GetWordCloudCache(UUID string, result *map[string]int) error
	GetFilterCloudCache(UUID string, result *map[string]int) error
	GetStringRankCache(UUID string, result *[]string) error
	GetAmountRankCache(UUID string, result *[]int) error

	DeleteWordCloudCache(UUID string) error
	DeleteFilterCloudCache(UUID string) error
	DeleteStringRankCache(UUID string) error
	DeleteAmountRankCache(UUID string) error
}
