package bootstrapper_test

import (
	"fmt"

	"github.com/ahmad/bootstrapper"
)

type Config struct {
	AppName string
}

func (c *Config) Booting() {
	fmt.Println("Loading configuration...")
	c.AppName = "MyApp"
	bootstrapper.Share("config", c)
}

func (c *Config) Booted() {
	fmt.Println("Configuration loaded")
}

type Database struct {
	Connected bool
}

func (d *Database) Booting() {
	fmt.Println("Connecting to database...")
	d.Connected = true
	bootstrapper.Share("db", d)
}

func (d *Database) Booted() {
	config := bootstrapper.Get("config").(*Config)
	fmt.Printf("Database ready for %s\n", config.AppName)
}

func Example() {
	// Register components with predefined orders
	bootstrapper.Register(&Config{}, bootstrapper.OrderConfig)
	bootstrapper.Register(&Database{}, bootstrapper.OrderDB)

	// Start initialization
	bootstrapper.StartBooting()
	bootstrapper.StartBooted()

	// Access shared state
	db := bootstrapper.Get("db").(*Database)
	fmt.Printf("Database connected: %v\n", db.Connected)

	// Output:
	// Loading configuration...
	// Connecting to database...
	// Configuration loaded
	// Database ready for MyApp
	// Database connected: true
}

type CustomService struct{}

func (c *CustomService) Booting() {
	fmt.Println("Custom service initializing")
}

func ExampleRegisterWithCustomOrder() {

	// Register with custom order
	customOrder := 50
	bootstrapper.RegisterWithCustomOrder(&CustomService{}, customOrder)
}
