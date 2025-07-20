package xlog

import (
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/petermattis/goid"
	"github.com/rs/zerolog"

	"github.com/xakepp35/pkg/xerrors"
)

var GetGID func() uint64 = GetGIDSafe

// EnableUnsafe для stack hook включает unsafe метод получения goid
func EnableUnsafe() {
	GetGID = GetGIDUnsafe
}

// GetGIDUnsafe быстро возвращает ID текущей горутины (unsafe, но быстро)
func GetGIDUnsafe() uint64 {
	return uint64(goid.Get())
}

func GetGIDSafe() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := string(buf[:n])

	const prefix = "goroutine "
	if !strings.HasPrefix(stack, prefix) {
		return 0
	}

	stack = stack[len(prefix):]
	end := strings.IndexByte(stack, ' ')
	if end < 0 {
		return 0
	}

	gid, err := strconv.ParseUint(stack[:end], 10, 64)
	if err != nil {
		return 0
	}

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
	defer s.mu.Unlock()
	s.m[gid] = value
}

func (s *SimpleMap) Load(gid uint64) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[gid]
	return v, ok
}

func (s *SimpleMap) Delete(gid uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, gid)
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
		return nil, xerrors.Err(nil).Str("storage", name).Msg("storage already registered")
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
		return nil, xerrors.Err(nil).Str("storage", name).Msg("storage already registered")
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
	gid := GetGID()
	if v, ok := h.storage.Load(gid); ok {
		e.Interface(h.name, v)
	}
}

// SetValue сохраняет value для текущей горутины в хранилище name
func SetValue(name string, value any) error {
	registryMu.RLock()
	defer registryMu.RUnlock()
	storage, exists := storages[name]
	if !exists {
		return xerrors.Err(nil).Str("storage", name).Msg("storage not found")
	}
	storage.Store(GetGID(), value)
	return nil
}

// DeleteValue удаляет значение для текущей горутины из хранилища name
func DeleteValue(name string) error {
	registryMu.RLock()
	defer registryMu.RUnlock()
	storage, exists := storages[name]
	if !exists {
		return xerrors.Err(nil).Str("storage", name).Msg("storage not found")
	}
	storage.Delete(GetGID())
	return nil
}
