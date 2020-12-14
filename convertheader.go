package convertheader

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)


const Uint64ToHex = "uint64tohex"
const Uint64ToInt64 = "uint64toint64"

type ReplaceValue struct {
	OldValue	string		`json:"oldValue,omitempty" toml:"oldValue,omitempty" yaml:"oldValue,omitempty"`
	NewValue	string		`json:"newValue,omitempty" toml:"newValue,omitempty" yaml:"newValue,omitempty"`
}

type Config struct {
	FromHeader		string	`json:"fromHeader,omitempty" toml:"fromHeader,omitempty" yaml:"fromHeader,omitempty"`
	CreateHeader	string	`json:"createHeader,omitempty" toml:"createHeader,omitempty" yaml:"createHeader,omitempty"`
	ConvertType		string	`json:"convertType,omitempty" toml:"convertType,omitempty" yaml:"convertType,omitempty"`
	ReplaceValues	[]ReplaceValue
	Prefix			string 	`json:"prefix,omitempty" toml:"prefix,omitempty" yaml:"prefix,omitempty"`
	Postfix			string 	`json:"postfix,omitempty" toml:"postfix,omitempty" yaml:"postfix,omitempty"`
}


// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ReplaceValues: []ReplaceValue{},
	}
}

// ConvertHeader holds the necessary components of a Traefik plugin
type ConvertHeader struct {
	next  	http.Handler
	config	Config
	name  	string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.FromHeader == "" || config.CreateHeader == "" {
		return nil, fmt.Errorf("some required fields are empty")
	}
	if config.ConvertType != "" && config.ConvertType != Uint64ToHex && config.ConvertType != Uint64ToInt64 {
		return nil, fmt.Errorf("not allowed value for convertType")
	}

	return &ConvertHeader{
		next:   next,
		config: *config,
		name:   name,
	}, nil
}

// Iterate over every headers to match the ones specified in the config and
// return nothing if regexp failed.
func (u *ConvertHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	header := req.Header.Get(u.config.FromHeader)
	newHeader := req.Header.Get(u.config.CreateHeader)
	if newHeader != "" {
		req.Header.Del(u.config.CreateHeader)
	}
	if u.config.ReplaceValues != nil {
		for _, replaceValue := range u.config.ReplaceValues {
			header = strings.ReplaceAll(header, replaceValue.OldValue, replaceValue.NewValue)
		}
	}
	if u.config.ConvertType != "" && header != ""{
		switch u.config.ConvertType {
		case Uint64ToHex:
			headerUInt64, err := strconv.ParseUint(header, 10, 64)
			if err != nil {
				log.Println(err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			header = strconv.FormatUint(headerUInt64, 16)
		case Uint64ToInt64:
			headerUInt64, err := strconv.ParseUint(header, 10, 64)
			if err != nil {
				log.Println(err.Error())
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			headerInt64 := int64(headerUInt64)
			header = strconv.FormatInt(headerInt64, 10)
		}
	}
	if u.config.Prefix != "" {
		header = u.config.Prefix + header
	}
	if u.config.Postfix != "" {
		header = header + u.config.Postfix
	}
	req.Header.Add(u.config.CreateHeader, header)
	u.next.ServeHTTP(rw, req)
}
