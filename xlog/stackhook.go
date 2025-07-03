package xlog

import (
	"fmt"
	"github.com/rs/zerolog"
	"runtime"
	"sync"
)

// getGID возвращает ID текущей горутины, извлекая его из runtime.Stack
func getGID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var gid uint64
	fmt.Sscanf(string(buf[:n]), "goroutine %d ", &gid)
	return gid
}

// GIDStorage — интерфейс для хранения значений по ID горутины
type GIDStorage interface {
	Store(uint64, any)
	Load(uint64) (any, bool)
	Delete(uint64)
}

// SimpleMap — простая потокобезопасная реализация GIDStorage на основе map
type SimpleMap struct {
	mu sync.RWMutex
	m  map[uint64]any
}

// NewSimpleMap создаёт новый SimpleMap
func NewSimpleMap() *SimpleMap {
	return &SimpleMap{m: make(map[uint64]any)}
}

func (s *SimpleMap) Store(gid uint64, value any) {
	s.mu.Lock()
	s.m[gid] = value
	s.mu.Unlock()
}

func (s *SimpleMap) Load(gid uint64) (any, bool) {
	s.mu.RLock()
	v, ok := s.m[gid]
	s.mu.RUnlock()
	return v, ok
}

func (s *SimpleMap) Delete(gid uint64) {
	s.mu.Lock()
	delete(s.m, gid)
	s.mu.Unlock()
}

// Регистрация хранилищ и хуков
var (
	registryMu sync.RWMutex
	storages   = make(map[string]GIDStorage)
)

// RegisterHook создаёт новый StackValueSetterHook с авто-managed хранилищем
// Если rawMap != nil, создаётся SimpleMap, обёрнутый вокруг rawMap.
func RegisterHook(name string, rawMap map[uint64]any) (*StackValueSetterHook, error) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := storages[name]; exists {
		return nil, fmt.Errorf("storage %s already registered", name)
	}
	var storage GIDStorage
	if rawMap != nil {
		storage = &SimpleMap{m: rawMap}
	} else {
		storage = NewSimpleMap()
	}
	storages[name] = storage
	return &StackValueSetterHook{name: name, storage: storage}, nil
}

// RegisterHookWithStorage позволяет использовать своё хранилище, реализующее GIDStorage
func RegisterHookWithStorage(name string, storage GIDStorage) (*StackValueSetterHook, error) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := storages[name]; exists {
		return nil, fmt.Errorf("storage %s already registered", name)
	}
	storages[name] = storage
	return &StackValueSetterHook{name: name, storage: storage}, nil
}

// StackValueSetterHook добавляет в событие значение из хранилища по GID
type StackValueSetterHook struct {
	name    string
	storage GIDStorage
}

// Run реализует zerolog.Hook: берёт текущее значение из storage и вставляет в поле с именем name
func (h *StackValueSetterHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	gid := getGID()
	if v, ok := h.storage.Load(gid); ok {
		e.Interface(h.name, v)
	}
}

// SetValue сохраняет value для текущей горутины в хранилище name
func SetValue(name string, value any) error {
	registryMu.RLock()
	storage, exists := storages[name]
	registryMu.RUnlock()
	if !exists {
		return fmt.Errorf("storage %s not found", name)
	}
	storage.Store(getGID(), value)
	return nil
}

// DeleteValue удаляет значение для текущей горутины из хранилища name
func DeleteValue(name string) error {
	registryMu.RLock()
	storage, exists := storages[name]
	registryMu.RUnlock()
	if !exists {
		return fmt.Errorf("storage %s not found", name)
	}
	storage.Delete(getGID())
	return nil
}
