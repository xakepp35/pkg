package xgrpc

import (
	"google.golang.org/grpc/metadata"
	"reflect"
	"sync"
	"testing"
)

func TestMethod_SetMethod(t *testing.T) {
	h := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h)

	if got := h.Method(); got != "" {
		t.Errorf("initial Method = %q, want empty", got)
	}

	h.SetMethod("/svc/meth")
	if want := "/svc/meth"; h.Method() != want {
		t.Errorf("Method() = %q, want %q", h.Method(), want)
	}
}

func TestSetHeader_MergesAndAppends(t *testing.T) {
	h := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h)

	in1 := metadata.MD{"a": {"1", "2"}, "b": {"x"}}
	if err := h.SetHeader(in1); err != nil {
		t.Fatal(err)
	}

	// Повторный вызов с новыми значениями
	in2 := metadata.MD{"a": {"3"}, "c": {"z"}}
	h.SetHeader(in2)

	want := metadata.MD{"a": {"1", "2", "3"}, "b": {"x"}, "c": {"z"}}
	if !reflect.DeepEqual(h.Headers(), want) {
		t.Errorf("Headers() = %+v; want %+v", h.Headers(), want)
	}

	// SetHeader(nil) не должен привести к панике и не изменить map
	h.SetHeader(nil)
	if !reflect.DeepEqual(h.Headers(), want) {
		t.Errorf("after nil SetHeader, headers changed: %+v", h.Headers())
	}
}

func TestSendHeader_EquivalentToSetHeader(t *testing.T) {
	h := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h)

	in := metadata.MD{"k": {"v"}}
	if err := h.SendHeader(in); err != nil {
		t.Fatal(err)
	}
	if got := h.Headers()["k"]; !reflect.DeepEqual(got, []string{"v"}) {
		t.Errorf("SendHeader did not set header, got %v", got)
	}
}

func TestSetTrailer_MergesAndAppends(t *testing.T) {
	h := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h)

	in1 := metadata.MD{"s": {"alpha"}}
	h.SetTrailer(in1)

	in2 := metadata.MD{"s": {"beta"}, "t": {"gamma"}}
	h.SetTrailer(in2)

	want := metadata.MD{"s": {"alpha", "beta"}, "t": {"gamma"}}
	if !reflect.DeepEqual(h.Trailers(), want) {
		t.Errorf("Trailers() = %+v; want %+v", h.Trailers(), want)
	}

	h.SetTrailer(metadata.MD{})
	if !reflect.DeepEqual(h.Trailers(), want) {
		t.Errorf("after empty SetTrailer, trailers changed: %+v", h.Trailers())
	}
}

func TestPoolReuse_ClearsBetweenUses(t *testing.T) {
	// Acquire, mutate, release
	h1 := AcquireHeaderStorage()
	h1.SetMethod("X")
	h1.SetHeader(metadata.MD{"k": {"v1"}})
	h1.SetTrailer(metadata.MD{"t": {"u1"}})
	ReleaseHeaderStorage(h1)

	// Acquire again — может вернуться тот же объект
	h2 := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h2)

	if h2.Method() != "" {
		t.Errorf("after Release, Method = %q; want empty", h2.Method())
	}
	if len(h2.Headers()) != 0 {
		t.Errorf("after Release, Headers len = %d; want 0", len(h2.Headers()))
	}
	if len(h2.Trailers()) != 0 {
		t.Errorf("after Release, Trailers len = %d; want 0", len(h2.Trailers()))
	}
}

func TestConcurrentAcquireRelease(t *testing.T) {
	const goroutines = 20
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			h := AcquireHeaderStorage()
			h.SetMethod("M")
			h.SetHeader(metadata.MD{"k": {string(byte('A' + i))}})
			h.SetTrailer(metadata.MD{"t": {string(byte('a' + i))}})
			ReleaseHeaderStorage(h)
		}(i)
	}
	wg.Wait()
	// После массового параллельного использования пул остаётся работоспособным
	h := AcquireHeaderStorage()
	defer ReleaseHeaderStorage(h)
	if h.Method() != "" || len(h.Headers()) != 0 || len(h.Trailers()) != 0 {
		t.Errorf("pool state not reset properly: Method=%q len(headers)=%d len(trailers)=%d",
			h.Method(), len(h.Headers()), len(h.Trailers()))
	}
}
