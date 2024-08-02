package core

type RedisCommand struct {
	Command string   // Redis command
	Args    []string // Array of arguments that comes with it
}
