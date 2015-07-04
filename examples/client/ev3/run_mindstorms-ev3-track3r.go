package main

import (
	. "github.com/zubairhamed/betwixt"
	. "github.com/zubairhamed/betwixt/examples/client/ev3/objects"
)

func main() {
	registry := NewDefaultObjectRegistry()
	c := NewDefaultClient(":0", "192.168.1.212:5683", registry)

	setupEv3Resources(c, registry)

	c.OnStartup(func() {
		c.Register("track3r")
	})

	c.Start()
}

func setupEv3Resources(client LWM2MClient, reg Registry) {
	client.SetEnabler(OMA_OBJECT_LWM2M_SECURITY, NewExampleSecurityObject(reg))
	client.AddObjectInstances(OMA_OBJECT_LWM2M_SECURITY, 0, 1, 2)

	client.SetEnabler(OMA_OBJECT_LWM2M_SERVER, NewExampleServerObject(reg))
	client.AddObjectInstance(OMA_OBJECT_LWM2M_SERVER, 1)

	client.SetEnabler(OMA_OBJECT_LWM2M_DEVICE, NewExampleDeviceObject(reg))
	client.AddObjectInstance(OMA_OBJECT_LWM2M_DEVICE, 0)

	client.EnableObject(OMA_OBJECT_LWM2M_ACCESS_CONTROL, NewExampleAccessControlObject(reg))
	client.AddObjectInstances(OMA_OBJECT_LWM2M_ACCESS_CONTROL, 0, 1, 2)

	client.EnableObject(OMA_OBJECT_LWM2M_CONNECTIVITY_MONITORING, NewExampleConnectivityMonitoringObject(reg))
	client.AddObjectInstance(OMA_OBJECT_LWM2M_CONNECTIVITY_MONITORING, 0)

	client.EnableObject(OMA_OBJECT_LWM2M_FIRMWARE_UPDATE, NewExampleFirmwareUpdateObject(reg))
	client.AddObjectInstance(OMA_OBJECT_LWM2M_FIRMWARE_UPDATE, 0)

	client.EnableObject(OMA_OBJECT_LWM2M_LOCATION, NewExampleLocationObject(reg))
	client.EnableObject(OMA_OBJECT_LWM2M_CONNECTIVITY_STATISTICS, NewExampleConnectivityStatisticsObject(reg))
}

