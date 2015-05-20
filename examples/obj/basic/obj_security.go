package basic


import (
    . "github.com/zubairhamed/lwm2m/api"
    "github.com/zubairhamed/lwm2m/core"
)


type Security struct {

}
func (o *Security) OnRead(instanceId int, resourceId int) (ResourceValue) {
    return core.NewEmptyValue()
}

/*
[0]
LWM2M Server URI        0   coap://bootstrap.example.com
Bootstrap Server        1   true
Security Mode           2   0
Public Key or Identity  3   [identity string]
Secret Key              4   [secret key data]
Short Server ID         10  0
Client Hold Off Time    11  3600

[1]
LWM2M Server URI        0   coap://server1.example.com
Bootstrap Server        1   false
Security Mode           2   0
Public Key or Identity  3   [identity string]
Secret Key              4   [secret key data]
Short Server ID         10  101
Client Hold Off Time    11  0

[2]
LWM2M Server URI        0   coap://server2.example.com
Bootstrap Server        1   false
Security Mode           2   0
Public Key or Identity  3   [identity string]
Secret Key              4   [secret key data]
Short Server ID         10  102
Client Hold Off Time    11  0
*/