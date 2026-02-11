package bootstrapper

import (
	"testing"
)

type MockService struct {
	bootingCalled bool
	bootedCalled  bool
}

func (m *MockService) Booting() {
	m.bootingCalled = true
}

func (m *MockService) Booted() {
	m.bootedCalled = true
}

func TestRegister(t *testing.T) {
	registry = nil
	mock := &MockService{}
	Register(mock, OrderService)

	if len(registry) != 1 {
		t.Errorf("Expected registry length 1, got %d", len(registry))
	}
	if registry[0].order != OrderService {
		t.Errorf("Expected order %d, got %d", OrderService, registry[0].order)
	}
}

func TestSortRegistry(t *testing.T) {
	registry = nil
	Register(&MockService{}, OrderService)
	Register(&MockService{}, OrderConfig)
	Register(&MockService{}, OrderDB)

	sortRegistry()

	if registry[0].order != OrderConfig {
		t.Errorf("Expected first order %d, got %d", OrderConfig, registry[0].order)
	}
	if registry[2].order != OrderService {
		t.Errorf("Expected last order %d, got %d", OrderService, registry[2].order)
	}
}

func TestShareAndGet(t *testing.T) {
	Shared = make(map[string]any)
	testKey := "testKey"
	testValue := "testValue"

	Share(testKey, testValue)
	result := Get(testKey)

	if result != testValue {
		t.Errorf("Expected %v, got %v", testValue, result)
	}
}

func TestGetNonExistent(t *testing.T) {
	Shared = make(map[string]any)
	result := Get("nonexistent")

	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestStartBooting(t *testing.T) {
	registry = nil
	mock1 := &MockService{}
	mock2 := &MockService{}

	Register(mock1, OrderConfig)
	Register(mock2, OrderService)

	StartBooting()

	if !mock1.bootingCalled {
		t.Error("Expected mock1 Booting to be called")
	}
	if !mock2.bootingCalled {
		t.Error("Expected mock2 Booting to be called")
	}
}

func TestStartBooted(t *testing.T) {
	registry = nil
	mock1 := &MockService{}
	mock2 := &MockService{}

	Register(mock1, OrderConfig)
	Register(mock2, OrderService)

	StartBooted()

	if !mock1.bootedCalled {
		t.Error("Expected mock1 Booted to be called")
	}
	if !mock2.bootedCalled {
		t.Error("Expected mock2 Booted to be called")
	}
}

func TestCustomOrder(t *testing.T) {
	registry = nil
	customOrder := 100
	mock := &MockService{}

	Register(mock, customOrder)

	if registry[0].order != customOrder {
		t.Errorf("Expected custom order %d, got %d", customOrder, registry[0].order)
	}
}

func TestConcurrentShareAndGet(t *testing.T) {
	Shared = make(map[string]any)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(n int) {
			key := string(rune('a' + n))
			Share(key, n)
			val := Get(key)
			if val != n {
				t.Errorf("Expected %d, got %v", n, val)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
