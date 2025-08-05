package xgrpc

import (
	"google.golang.org/grpc/metadata"
	"sync"
)

var (
	headerStoragePool = sync.Pool{
		New: func() interface{} {
			return &HeaderStorage{
				headers:  make(metadata.MD, 100),
				trailers: make(metadata.MD, 100),
			}
		},
	}
)

func AcquireHeaderStorage() *HeaderStorage {
	h := headerStoragePool.Get().(*HeaderStorage)
	return h
}

func ReleaseHeaderStorage(h *HeaderStorage) {
	clear(h.headers)
	clear(h.trailers)
	h.method = ""
	headerStoragePool.Put(h)
}

type HeaderStorage struct {
	method   string
	headers  metadata.MD
	trailers metadata.MD
}

func (h *HeaderStorage) Method() string {
	return h.method
}

func (h *HeaderStorage) SetMethod(method string) {
	h.method = method
}

func (h *HeaderStorage) Headers() metadata.MD {
	return h.headers
}

func (h *HeaderStorage) Trailers() metadata.MD {
	return h.trailers
}

func (h *HeaderStorage) SetHeader(md metadata.MD) error {
	if len(md) == 0 {
		return nil
	}

	for k, v := range md {
		h.headers[k] = append(h.headers[k], v...)
	}

	return nil
}

func (h *HeaderStorage) SendHeader(md metadata.MD) error {
	return h.SetHeader(md)
}

func (h *HeaderStorage) SetTrailer(md metadata.MD) error {
	if len(md) == 0 {
		return nil
	}

	for k, v := range md {
		h.trailers[k] = append(h.trailers[k], v...)
	}

	return nil
}
