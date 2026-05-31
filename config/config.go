package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	AppName string
	AppEnv  string
