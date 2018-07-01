package betwixt

import (
	"errors"
	"log"

	. "github.com/thingspin/canopus"
)

// NewLwm2mClient instantiates a new instance of LWM2M Client
func NewLwm2mClient(name, local, remote string, registry Registry) LWM2MClient {
	conn, _ := Dial(remote)
	//coapServer := NewServer(name, local, remote)
	coapServer := NewServer()

	// Create Mandatory
	c := &DefaultLWM2MClient{
		conn:           conn,
		coapServer:     coapServer,
		enabledObjects: make(map[LWM2MObjectType]Object),
		registry:       registry,
	}

	mandatory := registry.GetMandatory()
	for _, o := range mandatory {
		c.EnableObject(o.GetType(), NewNullEnabler())
	}

	return c
}

// DefaultLWM2MClient is client default.
type DefaultLWM2MClient struct {
	conn           Connection
	coapServer     CoapServer
	registry       Registry
	enabledObjects map[LWM2MObjectType]Object
	path           string

	// Events
	evtOnStartup FnOnStartup
	evtOnRead    FnOnRead
	evtOnWrite   FnOnWrite
	evtOnExecute FnOnExecute
	evtOnError   FnOnError
}

// Register this client to a LWM2M Server instance
// name must be unique and be less than 10 characers
func (c *DefaultLWM2MClient) Register(name string) (string, error) {
	if len(name) > 10 {
		return "", errors.New("Client name can not exceed 10 characters")
	}

	//req := NewRequest(MessageConfirmable, Post, GenerateMessageID())
	req := NewRequest(MessageConfirmable, Post)
	req.SetStringPayload(BuildModelResourceStringPayload(c.enabledObjects))
	req.SetRequestURI("/rd")
	req.SetURIQuery("ep", name)

	resp, err := c.conn.Send(req)

	path := ""
	if err != nil {
		return "", err
	} else {
		path = resp.GetMessage().GetLocationPath()
	}
	c.path = path

	return path, nil
}

// SetEnabler Sets/Defines an Enabler for a given LWM2M Object Type
func (c *DefaultLWM2MClient) SetEnabler(t LWM2MObjectType, e ObjectEnabler) {
	_, ok := c.enabledObjects[t]
	if ok {
		c.enabledObjects[t].SetEnabler(e)
	}
}

// GetEnabledObjects Returns a list of LWM2M Enabled Objects
func (c *DefaultLWM2MClient) GetEnabledObjects() map[LWM2MObjectType]Object {
	return c.enabledObjects
}

// GetRegistry Returns the registry used for looking up LWM2M object type definitions
func (c *DefaultLWM2MClient) GetRegistry() Registry {
	return c.registry
}

// Deregister this client from a LWM2M server which was previously registered
func (c *DefaultLWM2MClient) Deregister() {
	req := NewRequest(MessageConfirmable, Delete /*, GenerateMessageID()*/)

	req.SetRequestURI(c.path)
	_, err := c.conn.Send(req)

	if err != nil {
		log.Println(err)
	}
}

func (c *DefaultLWM2MClient) Update() {

}

func (c *DefaultLWM2MClient) AddResource() {

}

func (c *DefaultLWM2MClient) AddObject() {

}

func (c *DefaultLWM2MClient) UseRegistry(reg Registry) {
	c.registry = reg
}

// Registes an object enabler for a given LWM2M object type
func (c *DefaultLWM2MClient) EnableObject(t LWM2MObjectType, e ObjectEnabler) error {
	_, ok := c.enabledObjects[t]
	if !ok {
		if c.registry == nil {
			return errors.New("No registry found/set")
		}
		c.enabledObjects[t] = NewObject(t, e, c.registry)

		return nil
	} else {
		return errors.New("Object already enabled")
	}
}

// Adds a new object instance for a previously enabled LWM2M object type
func (c *DefaultLWM2MClient) AddObjectInstance(t LWM2MObjectType, instance int) error {
	o := c.enabledObjects[t]
	if o != nil {
		o.AddInstance(instance)

		return nil
	}
	return errors.New("Attempting to add a nil instance")
}

// Adds a list of object instance for a previously enabled LWM2M object type
func (c *DefaultLWM2MClient) AddObjectInstances(t LWM2MObjectType, instances ...int) {
	for _, o := range instances {
		c.AddObjectInstance(t, o)
	}
}

func (c *DefaultLWM2MClient) GetObject(n LWM2MObjectType) Object {
	return c.enabledObjects[n]
}

func (c *DefaultLWM2MClient) validate() {

}

// Start up the LWM2M client, listens to incoming requests and fires the OnStart event
func (c *DefaultLWM2MClient) Start() {
	c.validate()

	s := c.coapServer
	s.OnStart(func(server CoapServer) {
		if c.evtOnStartup != nil {
			c.evtOnStartup()
		}
	})

	s.OnObserve(func(resource string, msg Message) {
		log.Println("Observe Requested")
	})
	/*
		s.Get("/:obj/:inst/:rsrc", c.handleReadRequest)
		s.Get("/:obj/:inst", c.handleReadRequest)
		s.Get("/:obj", c.handleReadRequest)

		s.Put("/:obj/:inst/:rsrc", c.handleWriteRequest)
		s.Put("/:obj/:inst", c.handleWriteRequest)

		s.Delete("/:obj/:inst", c.handleDeleteRequest)

		s.Post("/:obj/:inst/:rsrc", c.handleExecuteRequest)
		s.Post("/:obj/:inst", c.handleCreateRequest)

		c.coapServer.Start()
	*/
}

// Handles LWM2M Create Requests (not to be mistaken for/not the same as  CoAP PUT)
func (c *DefaultLWM2MClient) handleCreateRequest(req CoapRequest) Response {
	log.Println("Create Request")
	attrResource := req.GetAttribute("rsrc")
	objectId := req.GetAttributeAsInt("obj")
	//sooskim instanceId := req.GetAttributeAsInt("inst")

	//sooskim var resourceId = -1

	if attrResource != "" {
		//sooskim resourceId = req.GetAttributeAsInt("rsrc")
	}

	t := LWM2MObjectType(objectId)
	obj := c.GetObject(t)
	enabler := obj.GetEnabler()

	msg := NewMessageOfType(MessageAcknowledgment, req.GetMessage().GetMessageId(), NewEmptyPayload())
	//sooskim msg.Token = req.GetMessage().Token
	//sooskim msg.Payload = NewEmptyPayload()

	if enabler != nil {
		//sooskim lwReq := Default(req, OPERATIONTYPE_CREATE)
		//sooskim response := enabler.OnCreate(instanceId, resourceId, lwReq)
		//sooskim msg.Code = response.GetResponseCode()
	} else {
		//sooskim msg.Code = CoapCodeMethodNotAllowed
	}
	return NewResponseWithMessage(msg)
}

// Handles LWM2M Read Requests (not to be mistaken for/not the same as  CoAP GET)
func (c *DefaultLWM2MClient) handleReadRequest(req CoapRequest) Response {
	attrResource := req.GetAttribute("rsrc")
	objectId := req.GetAttributeAsInt("obj")
	instanceId := req.GetAttributeAsInt("inst")

	var resourceId = -1

	if attrResource != "" {
		resourceId = req.GetAttributeAsInt("rsrc")
	}

	t := LWM2MObjectType(objectId)
	obj := c.GetObject(t)
	enabler := obj.GetEnabler()

	msg := NewMessageOfType(MessageAcknowledgment, req.GetMessage().GetMessageId(), NewEmptyPayload())
	//sooskim msg.Token = req.GetMessage().Token

	if enabler != nil {
		model := obj.GetDefinition()
		resource := model.GetResource(LWM2MResourceType(resourceId))

		if resource == nil {
			// TODO: Return TLV of Object Instance
			//sooskim msg.Code = CoapCodeNotFound
		} else {
			if !IsReadableResource(resource) {
				//sooskim msg.Code = CoapCodeMethodNotAllowed
			} else {
				lwReq := Default(req, OPERATIONTYPE_READ)
				response := enabler.OnRead(instanceId, resourceId, lwReq)

				val := response.GetResponseValue()
				//sooskim msg.Code = response.GetResponseCode()

				msg.AddOption(OptionContentFormat, MediaTypeFromValue(val))
				//sooskim b := EncodeValue(resource.GetId(), resource.MultipleValuesAllowed(), val)
				//sooskim msg.Payload = NewBytesPayload(b)
			}
		}
	} else {
		//sooskim msg.Code = CoapCodeMethodNotAllowed
	}
	return NewResponseWithMessage(msg)
}

// Handles LWM2M Delete Requests (not to be mistaken for/not the same as  CoAP DELETE)
func (c *DefaultLWM2MClient) handleDeleteRequest(req CoapRequest) Response {
	log.Println("Delete Request")
	objectId := req.GetAttributeAsInt("obj")
	instanceId := req.GetAttributeAsInt("inst")

	t := LWM2MObjectType(objectId)
	enabler := c.GetObject(t).GetEnabler()

	msg := NewMessageOfType(MessageAcknowledgment, req.GetMessage().GetMessageId(), NewEmptyPayload())
	//sooskim msg.Token = req.GetMessage().Token
	//sooskim msg.Payload = NewEmptyPayload()

	if enabler != nil {
		lwReq := Default(req, OPERATIONTYPE_DELETE)

		response := enabler.OnDelete(instanceId, lwReq)
		//sooskim msg.Code = response.GetResponseCode()
		msg.AddOption(OptionAccept, response.GetResponseCode()) //sooskim
	} else {
		//sooskim msg.Code = CoapCodeMethodNotAllowed
	}
	return NewResponseWithMessage(msg)
}

func (c *DefaultLWM2MClient) handleDiscoverRequest() {
	log.Println("Discovery Request")
}

func (c *DefaultLWM2MClient) handleObserveRequest() {
	log.Println("Observe Request")
}

// Handles LWM2M Write Requests (not to be mistaken for/not the same as  CoAP POST)
func (c *DefaultLWM2MClient) handleWriteRequest(req CoapRequest) Response {
	log.Println("Write Request")
	attrResource := req.GetAttribute("rsrc")
	objectId := req.GetAttributeAsInt("obj")
	instanceId := req.GetAttributeAsInt("inst")

	var resourceId = -1

	if attrResource != "" {
		resourceId = req.GetAttributeAsInt("rsrc")
	}

	t := LWM2MObjectType(objectId)
	obj := c.GetObject(t)
	enabler := obj.GetEnabler()

	msg := NewMessageOfType(MessageAcknowledgment, req.GetMessage().GetMessageId(), NewEmptyPayload())
	//sooskim msg.Token = req.GetMessage().Token
	//sooskim msg.Payload = NewEmptyPayload()

	if enabler != nil {
		model := obj.GetDefinition()
		resource := model.GetResource(LWM2MResourceType(resourceId))
		if resource == nil {
			// TODO Write to Object Instance
			//sooskim msg.Code = CoapCodeNotFound
		} else {
			if !IsWritableResource(resource) {
				//sooskim msg.Code = CoapCodeMethodNotAllowed
			} else {
				lwReq := Default(req, OPERATIONTYPE_WRITE)
				response := enabler.OnWrite(instanceId, resourceId, lwReq)
				//sooskim msg.Code = response.GetResponseCode()
				msg.AddOption(OptionAccept, response.GetResponseCode()) //sooskim
			}
		}
	} else {
		//sooskim msg.Code = CoapCodeNotFound
	}
	return NewResponseWithMessage(msg)
}

// Handles LWM2M Execute Requests
func (c *DefaultLWM2MClient) handleExecuteRequest(req CoapRequest) Response {
	log.Println("Execute Request")
	attrResource := req.GetAttribute("rsrc")
	objectId := req.GetAttributeAsInt("obj")
	instanceId := req.GetAttributeAsInt("inst")

	var resourceId = -1

	if attrResource != "" {
		resourceId = req.GetAttributeAsInt("rsrc")
	}

	t := LWM2MObjectType(objectId)
	obj := c.GetObject(t)
	enabler := obj.GetEnabler()

	msg := NewMessageOfType(MessageAcknowledgment, req.GetMessage().GetMessageId(), NewEmptyPayload())
	//sooskim msg.Token = req.GetMessage().Token
	//sooskim msg.Payload = NewEmptyPayload()

	if enabler != nil {
		model := obj.GetDefinition()
		resource := model.GetResource(LWM2MResourceType(resourceId))
		if resource == nil {
			//sooskim msg.Code = CoapCodeNotFound
		}

		if !IsExecutableResource(resource) {
			//sooskim msg.Code = CoapCodeMethodNotAllowed
		} else {
			lwReq := Default(req, OPERATIONTYPE_EXECUTE)
			response := enabler.OnExecute(instanceId, resourceId, lwReq)
			//sooskim msg.Code = response.GetResponseCode()
			msg.AddOption(OptionAccept, response.GetResponseCode()) //sooskim
		}
	} else {
		//sooskim msg.Code = CoapCodeNotFound
	}
	return NewResponseWithMessage(msg)
}

// Events
func (c *DefaultLWM2MClient) OnStartup(fn FnOnStartup) {
	c.evtOnStartup = fn
}

func (c *DefaultLWM2MClient) OnRead(fn FnOnRead) {
	c.evtOnRead = fn
}

func (c *DefaultLWM2MClient) OnWrite(fn FnOnWrite) {
	c.evtOnWrite = fn
}

func (c *DefaultLWM2MClient) OnExecute(fn FnOnExecute) {
	c.evtOnExecute = fn
}

func (c *DefaultLWM2MClient) OnError(fn FnOnError) {
	c.evtOnError = fn
}

func (c *DefaultLWM2MClient) OnObserve(fn FnOnError) {

}
