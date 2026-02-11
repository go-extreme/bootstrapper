package bootstrapper

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

type registeredItem struct {
	instance any
	order    int
}

const (
	OrderConfig = iota + 1 // 1
	OrderDB                // 2
	OrderCache             // 3
	OrderQueue             // 4
	WorkerPool             // 5
	OrderRepository        // 6
	OrderService           // 7
	OrderController        // 8
)

var registry []registeredItem

var Shared = make(map[string]any)

var mutex sync.RWMutex

func Register(instance any, order int) {
	registry = append(registry, registeredItem{
		instance: instance,
		order:    order,
	})
}

// RegisterWithCustomOrder allows registering with a custom order value
func RegisterWithCustomOrder(instance any, customOrder int) {
	Register(instance, customOrder)
}

// Call this when you want the registry sorted by order
func sortRegistry() {
	sort.Slice(registry, func(i, j int) bool {
		return registry[i].order < registry[j].order
	})

}

func Share(key string, value any) {
	mutex.Lock()         // write lock for writing
	defer mutex.Unlock() // unlock after function returns

	Shared[key] = value
}

func Get(key string) any {
	mutex.RLock()         // read lock for reading
	defer mutex.RUnlock() // unlock after function returns

	return Shared[key]
}

func StartBooting() {
	sortRegistry()

	for _, obj := range registry {
		// Ensure obj is a pointer to struct
		val := reflect.ValueOf(obj.instance)
		method := val.MethodByName("Booting")

		// Check if the method exists and is callable
		if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 0 {
			// Call the method
			results := method.Call(nil)

			// Expecting one return value of type error
			if len(results) == 1 && !results[0].IsNil() {
				err := results[0].Interface().(error)
				fmt.Printf("contracts.IBooting %s occurred error: %s\n", val.Type().String(), err.Error())
			}
		}
	}
}

func StartBooted() {

	for _, obj := range registry {
		// Ensure obj is a pointer to struct
		val := reflect.ValueOf(obj.instance)
		method := val.MethodByName("Booted")

		// Check if the method exists and is callable
		if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 0 {
			// Call the method
			results := method.Call(nil)

			// Expecting one return value of type error
			if len(results) == 1 && !results[0].IsNil() {
				err := results[0].Interface().(error)
				fmt.Printf("contracts.Booted %s occurred error: %s\n", val.Type().String(), err.Error())
			}
		}
	}
}
