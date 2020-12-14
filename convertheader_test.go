package convertheader_test

import (
	"context"
	ch "github.com/valebedeva/convertheader"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConvertHeader_ServeHTTP(t *testing.T) {
	testCases := []struct{
		desc			string
		header			string
		expectedValue 	string
		expectedCode 	int
		config			ch.Config
	} {
		{
			desc:				"Cut and convert to hex",
			header:				"SerialNumber%3D%229876543345678%22",
			expectedValue:    	"8fb8fdb940e",
			expectedCode:		http.StatusOK,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "uint64tohex",
				ReplaceValues: []ch.ReplaceValue{
					{
						OldValue: "SerialNumber%3D",
						NewValue: "",
					},
					{
						OldValue: "%22",
						NewValue: "",
					},
				},
			},
		},
		{
			desc:				"Cut and convert to int64 and add prefix and postfix",
			header:				"SerialNumber%3D%229876543334765445678%22",
			expectedValue:    	"SN=-8570200738944105938;",
			expectedCode:		http.StatusOK,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "uint64toint64",
				Prefix:        "SN=",
				Postfix:       ";",
				ReplaceValues: []ch.ReplaceValue{
					{
						OldValue: "SerialNumber%3D",
						NewValue: "",
					},
					{
						OldValue: "%22",
						NewValue: "",
					},
				},
			},
		},
		{
			desc:				"Parsing value out of range",
			header:				"2345678908765432123456789",
			expectedValue:    	"",
			expectedCode:		http.StatusInternalServerError,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "uint64toint64",
			},
		},
		{
			desc:				"Empty header value",
			header:				"",
			expectedValue:    	"",
			expectedCode:		http.StatusOK,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "uint64toint64",
			},
		},
		{
			desc:				"Can't parse string Parsing",
			header:				"hello",
			expectedValue:    	"",
			expectedCode:		http.StatusInternalServerError,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "uint64toint64",
			},
		},
		{
			desc:				"Cut and convert to hex",
			header:				"SomeValue",
			expectedValue:    	"SomeValue/new",
			expectedCode:		http.StatusOK,
			config: 			ch.Config{
				FromHeader:    "TEST-HEADER",
				CreateHeader:  "EXISTING-HEADER",
				Postfix:       "/new",
			},
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			fat, _ := ch.New(context.Background(), next, &test.config, "TestConvertHeader")
			r := httptest.NewRequest("GET", "http://example.com", nil)
			r.Header.Set(test.config.FromHeader, test.header)
			r.Header.Set("EXISTING-HEADER", "OldValue")
			w := httptest.NewRecorder()
			fat.ServeHTTP(w, r)
			if test.expectedCode != w.Code {
				t.Errorf("Expexted code: %d; Recieved: %d", test.expectedCode, w.Code)
			}
			value := r.Header.Get(test.config.CreateHeader)
			if test.expectedValue != value {
				t.Errorf("Expexted new header value: %s; Recieved value: %s", test.expectedValue, r.Header.Get(test.config.CreateHeader))
			}
		})
	}
}

func TestConvertHeader_New(t *testing.T) {
	testCases := []struct{
		desc			string
		expectedError	error
		config			ch.Config
	} {
		{
			desc: 		        "Ok header and empty convert type",
			expectedError:		nil,
			config:  		    ch.Config{
				FromHeader:    "TEST-HAEDER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "",
				ReplaceValues: nil,
				Prefix:        "",
				Postfix:       "",
			},
		},
		{
			desc: 		        "Ok header and not allowed convert type",
			expectedError:		fmt.Errorf(""),
			config:  		    ch.Config{
				FromHeader:    "TEST-HAEDER",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "hello",
				ReplaceValues: nil,
				Prefix:        "",
				Postfix:       "",
			},
		},
		{
			desc: 		        "Empty from header",
			expectedError:		fmt.Errorf(""),
			config:  		    ch.Config{
				FromHeader:    "",
				CreateHeader:  "NEW-HEADER",
				ConvertType:   "",
				ReplaceValues: nil,
				Prefix:        "",
				Postfix:       "",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
			_, err := ch.New(ctx, next, &test.config, "Test creating middleware")
			if test.expectedError == nil{
				if err != nil {
					t.Errorf("Expexted nil error; Recieved not nil: %s", err)
				}
			}
			if test.expectedError != nil{
				if err == nil {
					t.Errorf("Expexted error; Recieved nil error")
				}
			}
		})
	}
}