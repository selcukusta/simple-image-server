package helper

import (
	"regexp"
	"strconv"
	"strings"
)

//IsRouteFit is using to check route is hitting the provided patterns
func IsRouteFit(patterns [4]string, url string) (bool, map[string]string) {
	variables := make(map[string]string)
	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(url)
		if len(matches) <= 0 {
			continue
		}

		for i, name := range regex.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			if split := strings.Split(name, "_"); len(split) == 4 && split[1] == "r" {
				min, err := strconv.Atoi(split[2])
				if err != nil {
					return false, nil
				}

				max, err := strconv.Atoi(split[3])
				if err != nil {
					return false, nil
				}

				if !ValidateRange(matches[i], min, max) {
					return false, nil
				}

				variables[split[0]] = matches[i]
			} else {
				variables[name] = matches[i]
			}
		}
		return true, variables

	}
	return false, nil
}
