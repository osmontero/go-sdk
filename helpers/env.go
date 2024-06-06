package helpers

import (
	"os"
	"strconv"
	"strings"

	"github.com/threatwinds/logger"
)

type Env struct {
	ClusterPort     int
	RestPort        int
	GrpcPort        int
	Workdir         string
	RulesRepository string
	SearchNodes     []string
	LogLevel        int
}

func getEnvStr(name, def string, required bool) (string, *logger.Error) {
	val := os.Getenv(name)

	if val == "" {
		if required {
			return "", Logger().ErrorF("configuration required: %s", name)
		} else {
			return def, nil
		}
	}

	return val, nil
}

func getEnvInt(name string, def string, required bool) (int, *logger.Error) {
	str, e := getEnvStr(name, def, required)
	if e != nil {
		return 0, e
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, Logger().ErrorF(err.Error())
	}

	return val, nil
}

func getEnvStrSlice(name, def string, required bool) ([]string, *logger.Error) {
	str, e := getEnvStr(name, def, required)
	if e != nil {
		return nil, e
	}

	var items = make([]string, 0, 1)
	for _, item := range strings.Split(str, ",") {
		items = append(items, strings.TrimSpace(item))
	}

	return items, nil
}

func getEnv() Env {
	var env Env
	var e *logger.Error
	
	env.ClusterPort, e = getEnvInt("CLUSTER_PORT", "8082", false)
	if e != nil {
		panic(e.Message)
	}

	env.RestPort, e = getEnvInt("REST_PORT", "8080", false)
	if e != nil {
		panic(e.Message)
	}

	env.GrpcPort, e = getEnvInt("GRPC_PORT", "8081", false)
	if e != nil {
		panic(e.Message)
	}

	env.Workdir, e = getEnvStr("WORK_DIR", "", true)
	if e != nil {
		panic(e.Message)
	}

	env.SearchNodes, e = getEnvStrSlice("SEARCH_NODES", "", true)
	if e != nil {
		panic(e.Message)
	}

	env.RulesRepository, e = getEnvStr("RULES_REPOSITORY", "", true)
	if e != nil {
		panic(e.Message)
	}

	env.LogLevel, e = getEnvInt("LOG_LEVEL", "200", false)
	if e != nil {
		panic(e.Message)
	}

	return env
}