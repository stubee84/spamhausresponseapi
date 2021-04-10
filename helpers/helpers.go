package helpers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var IpReg *regexp.Regexp = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

const Zen string = ".zen.spamhaus.org"

var creds map[string]string = map[string]string{
	"username": os.Getenv("SpamhausUsername"),
	"password": os.Getenv("SpamhausPassword"),
}

func ReverseUsingSeperator(str, sep string) (string, error) {
	splitStr := strings.Split(str, sep)
	result := ""

	for i := len(splitStr) - 1; i >= 0; i-- {
		if ip, err := strconv.Atoi(splitStr[i]); err == nil {
			if ip > 255 {
				return "", fmt.Errorf("error: octect value, %d, is greater than 255", ip)
			} else if err != nil {
				return "", fmt.Errorf("error: could not convert string to int. %s", splitStr[i])
			}
		}
		result = result + sep + splitStr[i]
	}
	return strings.Trim(result, sep), nil
}

func BasicAuth() func(http.Handler) http.Handler {
	return func(handle http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()

			if !ok {
				log.Printf("could not parse authorization\n")
				http.Error(w, "Internal error", http.StatusBadRequest)
				return
			}

			if user != creds["username"] || pass != creds["password"] {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			handle.ServeHTTP(w, r)
		})
	}
}
