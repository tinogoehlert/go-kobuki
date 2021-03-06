package kobuki

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/tinogoehlert/gobuki"
	"github.com/tinogoehlert/gobuki/commands"
	"github.com/tinogoehlert/gobuki/sensors"
	"gobot.io/x/gobot"
)

const (
	// FeedbackEvent event
	FeedbackEvent = "Feedback"
	// GyroEvent event
	GyroEvent = "Gyro"
	// CliffEvent event
	CliffEvent = "Cliff"
	// WheelsDropEvent event
	WheelsDropEvent = "WheelsDrop"
	// WheelsEncoderEvent event
	WheelsEncoderEvent = "WheelsEncoder"
	// WheelsPWMEvent event
	WheelsPWMEvent = "wheelsPWM"
	// BumperEvent event
	BumperEvent = "Bumper"
	// ButtonsEvent event
	ButtonsEvent = "Button"
	// ChargeStateEvent event
	ChargeStateEvent = "ChargeState"
	// DockingIREvent event
	DockingIREvent = "DockingIR"
	// CliffADCEvent event
	CliffADCEvent = "CliffADC"
	// InertialEvent event
	InertialEvent = "Inertial"
	// BatteryVoltageEvent event
	BatteryVoltageEvent = "BatteryVoltage"
)

const (
	// SoundOn sound
	SoundOn = commands.ON
	// SoundOff sound
	SoundOff = commands.OFF
	// SoundRecharge sound
	SoundRecharge = commands.RECHARGE
	// SoundCleaningStart sound
	SoundCleaningStart = commands.CLEANINGSTART
	// SoundCleaningEnd sound
	SoundCleaningEnd = commands.CLEANINGEND
)

// Feedback holds feedback data
type Feedback struct {
	Name string
	Data interface{}
}

// Driver is the interface that describes a driver in gobot
type Driver struct {
	name    string
	adaptor *Adaptor
	started func()
	gobot.Eventer
	mover   *gobuki.Mover
	encoder gob.Encoder
}

// NewDriver creates a new Kobuki Bot driver
func NewDriver(a *Adaptor) *Driver {

	d := &Driver{
		name:    "Kobuki",
		adaptor: a,
		Eventer: gobot.NewEventer(),
		started: func() {},
	}

	d.AddEvent(FeedbackEvent)
	return d
}

// Name returns the label for the Driver
func (d *Driver) Name() string {
	return d.name
}

// SetName sets the label for the Driver
func (d *Driver) SetName(s string) {
	d.name = s
}

// Start initiates the Driver
func (d *Driver) Start() error {
	d.mover = gobuki.NewMover(d.adaptor.bot)
	d.adaptor.bot.AddModule(d.mover)
	d.adaptor.bot.Start()
	d.started()
	d.adaptor.bot.OnAll(func(name string, data interface{}) {
		d.Publish(FeedbackEvent, Feedback{
			Name: name,
			Data: data,
		})
	})

	go func() {
		for {
			fmt.Println(d.adaptor.bot.LogChannel())
		}
	}()

	return nil
}

// Halt terminates the Driver
func (d *Driver) Halt() error {
	fmt.Println("halt called")
	return nil
}

// SetCliffADCTolerance set tolerance for Cliff ADC
func (d *Driver) SetCliffADCTolerance(t int) error {
	if d.adaptor.bot == nil {
		return errors.New("driver not started")
	}
	d.adaptor.bot.SetCliffADCTolerance(t)
	return nil
}

// SetGyroTolerance set tolerance for gyroscope
func (d *Driver) SetGyroTolerance(t float64) error {
	if d.adaptor.bot == nil {
		return errors.New("driver not started")
	}
	d.adaptor.bot.SetGyroTolerance(t)
	return nil
}

// SetCurrentWheelsTolerance set tolerance for wheels current
func (d *Driver) SetCurrentWheelsTolerance(t int) error {
	if d.adaptor.bot == nil {
		return errors.New("driver not started")
	}
	d.adaptor.bot.SetCurrentWheelsTolerance(t)
	return nil
}

// Move moves the robot
func (d *Driver) Move(vx, wz float64) {
	d.mover.Move(vx, wz)
}

// MoveRaw moves the robot also
func (d *Driver) MoveRaw(speed, radius int16) {
	d.mover.MoveRaw(speed, radius)
}

// PlaySoundSequence plays a sound sequence on the robot
func (d *Driver) PlaySoundSequence(sequence commands.SoundSequence) {
	d.adaptor.bot.Send(commands.SoundSequenceCmd(sequence))
}

// Connection returns the Connection associated with the Driver
func (d *Driver) Connection() gobot.Connection {
	return d.adaptor
}

// OnStart triggered when driver was started, useful for setup
func (d *Driver) OnStart(f func()) {
	d.started = f
}

// OnGyro new GyroData available
func (d *Driver) OnGyro(f func(*sensors.GyroData)) {
	d.adaptor.bot.On(GyroEvent, func(data interface{}) {
		f(data.(*sensors.GyroData))
	})
}

// OnCliff Cliff sensor changed
func (d *Driver) OnCliff(f func(*sensors.Cliff)) {
	d.adaptor.bot.On(CliffEvent, func(data interface{}) {
		f(data.(*sensors.Cliff))
	})
}

// OnWheelEncoder WheelEncoder data changed
func (d *Driver) OnWheelEncoder(f func(*sensors.WheelsEncoder)) {
	d.adaptor.bot.On(WheelsEncoderEvent, func(data interface{}) {
		f(data.(*sensors.WheelsEncoder))
	})
}

// OnWheelDrop WheelDrop data changed
func (d *Driver) OnWheelDrop(f func(*sensors.WheelsDrop)) {
	d.adaptor.bot.On(WheelsDropEvent, func(data interface{}) {
		f(data.(*sensors.WheelsDrop))
	})
}

// OnWheelPWM WheelPWM data changed
func (d *Driver) OnWheelPWM(f func(*sensors.WheelsPWM)) {
	d.adaptor.bot.On(WheelsPWMEvent, func(data interface{}) {
		f(data.(*sensors.WheelsPWM))
	})
}

// OnInertial Inertial data changed
func (d *Driver) OnInertial(f func(*sensors.Inertial)) {
	d.adaptor.bot.On(InertialEvent, func(data interface{}) {
		f(data.(*sensors.Inertial))
	})
}

// OnBumper Bumper state changed
func (d *Driver) OnBumper(f func(*sensors.Bumper)) {
	d.adaptor.bot.On(BumperEvent, func(data interface{}) {
		f(data.(*sensors.Bumper))
	})
}

// OnButtons Buttons state changed
func (d *Driver) OnButtons(f func(*sensors.Buttons)) {
	d.adaptor.bot.On(ButtonsEvent, func(data interface{}) {
		f(data.(*sensors.Buttons))
	})
}

// OnChargeState Charge state changed
func (d *Driver) OnChargeState(f func(*sensors.ChargeState)) {
	d.adaptor.bot.On(ChargeStateEvent, func(data interface{}) {
		f(data.(*sensors.ChargeState))
	})
}

// OnBatteryVoltage Charge state changed
func (d *Driver) OnBatteryVoltage(f func(*uint8)) {
	d.adaptor.bot.On(ChargeStateEvent, func(data interface{}) {
		f(data.(*uint8))
	})
}

// OnWheelsCurrent Wheels current state changed
func (d *Driver) OnWheelsCurrent(f func(*sensors.CurrentWheels)) {
	d.adaptor.bot.On(ChargeStateEvent, func(data interface{}) {
		f(data.(*sensors.CurrentWheels))
	})
}

// OnDockingIR DockingIR data changed
func (d *Driver) OnDockingIR(f func(*sensors.DockingIR)) {
	d.adaptor.bot.On(DockingIREvent, func(data interface{}) {
		f(data.(*sensors.DockingIR))
	})
}
