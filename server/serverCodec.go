package server

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gorilla/rpc"
	"github.com/kolo/xmlrpc"
)

type Codec struct {
	methods                  map[string]string
	defaultMethodByNamespace map[string]string
	defaultMethod            string
	parsers                  map[string]Parser
}

func NewCodec() *Codec {
	return &Codec{
		methods:                  make(map[string]string),
		defaultMethodByNamespace: make(map[string]string),
		defaultMethod:            "",
		parsers:                  make(map[string]Parser),
	}
}

func (c *Codec) RegisterMethod(method string, parser Parser) {
	c.methods[method] = method
	c.parsers[c.resolveMethod(method)] = parser
}

func (c *Codec) RegisterDefaultMethod(method string, parser Parser) {
	c.defaultMethod = method
	c.parsers[method] = parser
}

func (c *Codec) RegisterDefaultMethodForNamespace(namespace, method string, parser Parser) {
	c.defaultMethodByNamespace[namespace] = method
	c.parsers[c.resolveMethod(method)] = parser
}

func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	rawxml, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &CodecRequest{err: err}
	}
	defer r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(rawxml))

	var request ServerRequest
	if err := xml.Unmarshal(rawxml, &request); err != nil {
		return &CodecRequest{err: err}
	}
	request.rawxml = rawxml
	request.Method = c.resolveMethod(request.Method)

	parser := c.parsers[request.Method]

	return &CodecRequest{request: &request, parser: parser}
}

func (c *Codec) resolveMethod(requestMethod string) string {
	namespace, methodStr := getNamespaceAndMethod(requestMethod)
	if _, ok := c.methods[requestMethod]; ok {
		return toLowerCase(namespace, methodStr)
	} else if method, ok := c.defaultMethodByNamespace[namespace]; ok {
		return method
	} else if c.defaultMethod != "" {
		return c.defaultMethod
	}
	return requestMethod
}

func getNamespaceAndMethod(requestMethod string) (string, string) {
	//TODO:
	if len(requestMethod) > 1 {
		parts := strings.Split(requestMethod, ".")
		slice := parts[1:len(parts)]
		return parts[0], strings.Join(slice, ".")
	}
	return "", ""
}

func toLowerCase(namespace, method string) string {
	//TODO:
	if namespace != "" && method != "" {
		r, n := utf8.DecodeRuneInString(method)
		if unicode.IsLower(r) {
			return namespace + "." + string(unicode.ToUpper(r)) + method[n:]
		}
	}
	return namespace + "." + method
}

type ServerRequest struct {
	Name   xml.Name `xml:"methodCall"`
	Method string   `xml:"methodName"`
	rawxml []byte
}

type CodecRequest struct {
	request *ServerRequest
	err     error
	parser  Parser
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	val := reflect.ValueOf(args)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer value passed")
	}

	var argsList []interface{}
	argsList, c.err = xmlrpc.UnmarshalToList(c.request.rawxml)

	err := c.parser.Parse(argsList, args)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, response interface{}, methodErr error) error {
	var xmlstr string
	if c.err != nil {
		//TODO:
		/*	var fault Fault
			switch c.err.(type) {
			case Fault:
				fault = c.err.(Fault)
			default:
				fault = FaultApplicationError
				fault.String += fmt.Sprintf(": %v", c.err)
			}
			xmlstr = fault2XML(fault)*/
		return c.err
	} else if methodErr != nil {

		var fault Fault
		switch methodErr.(type) {
		case Fault:
			fault = methodErr.(Fault)
		default:
			fault = FaultApplicationError
			fault.String += fmt.Sprintf(": %v", methodErr)
		}
		xmlstr = fault2XML(fault)
	} else {
		xmlstr, _ = encodeResponseToXML(response)
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write([]byte(xmlstr))
	return nil
}
