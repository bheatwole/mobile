// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sensor provides sensor events from various movement sensors.
package sensor // import "golang.org/x/mobile/exp/sensor"

import (
	"errors"
	"time"
)

// Type represents a sensor type.
type Type int

var sensorNames = map[Type]string{
	Accelerometer: "Accelerometer",
	Gyroscope:     "Gyrsocope",
	Magnetometer:  "Magnetometer",
}

// String returns the string representation of the sensor type.
func (t Type) String() string {
	if n, ok := sensorNames[t]; ok {
		return n
	}
	return "Unknown sensor"
}

const (
	Accelerometer = Type(0)
	Gyroscope     = Type(1)
	Magnetometer  = Type(2)
	nTypes        = Type(3)
)

// Event represents a sensor event.
type Event struct {
	// Sensor is the type of the sensor the event is coming from.
	Sensor Type

	// Timestamp is a device specific event time in nanoseconds.
	// Timestamps are not Unix times, they represent a time that is
	// only valid for the device's default sensor.
	Timestamp int64

	// Data is the event data.
	//
	// If the event source is Accelerometer,
	//  - Data[0]: acceleration force in x axis in m/s^2
	//  - Data[1]: acceleration force in y axis in m/s^2
	//  - Data[2]: acceleration force in z axis in m/s^2
	//
	// If the event source is Gyroscope,
	//  - Data[0]: rate of rotation around the x axis in rad/s
	//  - Data[1]: rate of rotation around the y axis in rad/s
	//  - Data[2]: rate of rotation around the z axis in rad/s
	//
	// If the event source is Magnetometer,
	//  - Data[0]: force of gravity along the x axis in m/s^2
	//  - Data[1]: force of gravity along the y axis in m/s^2
	//  - Data[2]: force of gravity along the z axis in m/s^2
	//
	Data []float64
}

// TODO(jbd): Move Sender interface definition to a top-level package.

// Sender sends an event.
type Sender interface {
	Send(event interface{})
}

// m is the underlying platform-specific sensor manager.
var m = newManager()

// Enable enables the specified sensor type with the given delay rate.
// Sensor events will be sent to s, a typical example of Sender
// implementations is app.App.
// Enable is not safe for concurrent use.
func Enable(s Sender, t Type, delay time.Duration) error {
	if t < 0 || int(t) >= len(sensorNames) {
		return errors.New("sensor: unknown sensor type")
	}
	return m.enable(s, t, delay)
}

// Disable disables to feed the manager with the specified sensor.
// Disable is not safe for concurrent use.
func Disable(t Type) error {
	if t < 0 || int(t) >= len(sensorNames) {
		return errors.New("sensor: unknown sensor type")
	}
	return m.disable(t)
}

func newManager() *manager {
	// TODO(jbd): manager type is unnecessary, flatten out the
	// platform specific implementation.
	mgr := new(manager)
	mgr.initialize()
	return mgr
}
