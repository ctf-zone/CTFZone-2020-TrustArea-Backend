package config

import "time"

const BuildVersion = "1.0.0"

var Config struct {
	App struct {
		Host string `json:"host"`
		Port int `json:"port"`
	} `json:"app"`
	DB struct {
		Host string `json:"host"`
		Port int `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
		DBName string `json:"db_name"`
	} `json:"db"`
	Redis struct {
		Host string `json:"host"`
		Port int `json:"port"`
		SessionTTL time.Duration `json:"session_ttl"`
	} `json:"redis"`
}