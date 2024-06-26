package mqueue

import (
	"github.com/streadway/amqp"
	"strings"
)

type Options map[string]interface{}

var opts = map[string]interface{}{
	"contentType":  "application/json",
	"deliveryMode": 1,
	"durable":      true,
	"autoDelete":   false,
	"internal":     false,
	"exclusive":    false,
	"noWait":       false,
	"noAck":        false,
	"noLocal":      false,
}

func Ucfirst(str string) string {
	bytes := []byte(str)
	bytes[0] &= 0xdf
	return string(bytes)
}

func CamelStrings(s string) string {
	var result string
	if re := strings.Replace(s, " ", "", -1); re != "" {
		for i, v := range strings.Split(re, "_") {
			if i == 0 {
				result = v
			} else {
				result += Ucfirst(v)
			}
		}
	}
	return result
}

func (options Options) normalizeKeys() {
	for name, value := range options {
		delete(options, name)
		correctName := CamelStrings(name)
		options[correctName] = value
	}
}

func (options Options) buildArgs() {
	if _, ok := options["args"]; !ok {
		return
	}

	args := options.convertArgsToAMQPTable(options["args"])
	args = options.fixArgsValuesTypes(args)
	options["args"] = args
}

func (options Options) convertArgsToAMQPTable(args interface{}) amqp.Table {
	var table amqp.Table

	switch arguments := args.(type) {
	case map[string]interface{}:
		table = amqp.Table(arguments)
	case map[interface{}]interface{}:
		table = make(amqp.Table, len(arguments))

		for k, v := range arguments {
			table[k.(string)] = v
		}
	}

	return table
}

func (options Options) fixArgsValuesTypes(args amqp.Table) amqp.Table {
	for k, v := range args {
		switch v2 := v.(type) {
		case int:
			args[k] = int32(v2)
		case float64:
			if k == "x-max-priority" {
				args[k] = int64(v2)
			}
		default:
			args[k] = v
		}
	}
	return args
}
