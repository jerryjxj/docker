package loghub

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Test_initAndValidate(t *testing.T) {

	var config loghubConfig
	fmt.Print(config)
	fmt.Print(config.Drivers)
}

func Test_extractLabels(t *testing.T) {
	config := map[string]string{
		"labels/node":      "192.168.1.22",
		"labels/namespace": "$(io.kubernetes.pod.namespace)",
		"labels/pod":       "$(io.kubernetes.pod.name)",
	}
	labels := map[string]string{
		"io.kubernetes.pod.namespace": "cy",
		"io.kubernetes.pod.name":      "bss-234766563-df311",
	}
	re := regexp.MustCompile(`\$\((.*?)\)`)
	for k, v := range config {
		if strings.Index(k, "labels/") == 0 {
			if tokens := strings.SplitN(k, "/", 2); len(tokens) == 2 {
				if l := re.FindStringSubmatch(v); len(l) > 0 {
					fmt.Printf("%s(%s) -> %s\n", tokens[1], l, labels[l[1]])
				} else {
					fmt.Printf("%s -> %s\n", tokens[1], v)
				}
			}
		}
	}
}

func Test_validateLabels(t *testing.T) {
	key := "labels/landscape"
	if strings.Index(key, "labels/") != 0 {
		fmt.Printf("unknown log opt '%s' for redis log driver", key)
	} else {
		fmt.Println("ok!")
	}
}
