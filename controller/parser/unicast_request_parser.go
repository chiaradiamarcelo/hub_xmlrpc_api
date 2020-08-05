package parser

import (
	"log"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/xmlrpc"
)

func UnicastRequestParser(request *xmlrpc.ServerRequest, output interface{}) error {
	parsedArgs, ok := output.(*controller.UnicastRequest)
	if !ok {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultInvalidParams
	}

	args := request.Params
	if len(args) < 2 {
		log.Printf("Error ocurred when parsing arguments")
		return controller.FaultWrongArgumentsNumber
	}

	hubSessionKey, ok := args[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing hubSessionKey argument")
		return controller.FaultInvalidParams
	}

	serverID, ok := args[1].(int64)
	if !ok {
		log.Printf("Error ocurred when parsing serverID argument")
		return controller.FaultInvalidParams
	}

	rest := args[2:len(args)]
	serverArgs := make([]interface{}, len(rest))
	for i, arg := range rest {
		serverArgs[i] = (interface{})(arg)
	}

	method, err := removeNamespace(request.MethodName)
	if err != nil {
		return err
	}

	*parsedArgs = controller.UnicastRequest{HubSessionKey: hubSessionKey, Call: method, ServerID: serverID, Args: serverArgs}
	return nil
}
