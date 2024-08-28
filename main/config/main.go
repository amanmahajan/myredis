package config

var Host string = "0.0.0.0"

var Port int = 7370

var EvictionStrategy string = "allkeys-random"

var AOFFile string = "./persistence.aof"

var KeyLimits int = 100

var EvictionRatio float64 = 0.4
