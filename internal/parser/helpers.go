package parser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func convertToNum(element string) int {
	var number []string
	for _, char := range element {
		if char >= '0' && char <= '9' {
			number = append(number, string(char))
		}
	}
	result := strings.Join(number, "")
	num, _ := strconv.Atoi(result)
	log.
		WithFields(log.Fields{
			"Raw":     element,
			"_Output": num,
			"_Type":   reflect.TypeOf(num).String(),
		}).
		Debug("Strip number and convert to int")
	return num
}

func findMinimal(data []string, sep string) int {
	var (
		price  int
		prices []int
	)

	for _, value := range data {
		if strings.Contains(value, sep) && sep != "" {
			log.WithFields(log.Fields{
				"Separator": sep,
				"Raw":       value,
			}).Debug("Multi price")

			splitted := strings.Split(value, sep)
			for _, element := range splitted {
				num := convertToNum(element)
				if num > 0 {
					prices = append(prices, int(num))
				}
			}
		} else {
			log.WithFields(log.Fields{
				"Raw": value,
			}).Debug("One price")
			num := convertToNum(value)
			if num > 0 {
				prices = append(prices, num)
			}
		}
	}

	log.WithFields(log.Fields{
		"Content": prices,
		"Type":    reflect.TypeOf(prices).String(),
	}).Debug("List of prices")

	if prices != nil {
		price = slices.Min(prices)
	} else {
		fmt.Println("Empty")
	}
	return price
}
