package screenshot

import "testing"

func TestRedrawSurfaceDefersWhenAllRenderSlotsBusy(t *testing.T) {
	screenBuf, err := CreateShmBuffer(2, 2, 8)
	if err != nil {
		t.Fatalf("CreateShmBuffer() error = %v", err)
	}
	defer screenBuf.Close()

	os := &OutputSurface{
		screenBuf:  screenBuf,
		slotsReady: true,
	}
	for i := range os.slots {
		os.slots[i] = &RenderSlot{busy: true}
	}

	var selector RegionSelector
	selector.redrawSurface(os)

	if !os.needsRedraw {
		t.Fatal("redrawSurface should mark a pending redraw when all render slots are busy")
	}
}

func TestCreateDimmedBufferCopyDimsSourcePixels(t *testing.T) {
	src, err := CreateShmBuffer(1, 1, 4)
	if err != nil {
		t.Fatalf("CreateShmBuffer() error = %v", err)
	}
	defer src.Close()

	source := []byte{100, 150, 200, 255}
	copy(src.Data(), source)

	dimmed, err := createDimmedBufferCopy(src)
	if err != nil {
		t.Fatalf("createDimmedBufferCopy() error = %v", err)
	}
	defer dimmed.Close()

	want := []byte{60, 90, 120, 255}
	for i, got := range dimmed.Data()[:4] {
		if got != want[i] {
			t.Fatalf("dimmed.Data()[%d] = %d, want %d", i, got, want[i])
		}
	}

	for i, got := range src.Data()[:4] {
		if got != source[i] {
			t.Fatalf("src.Data()[%d] = %d, want %d", i, got, source[i])
		}
	}
}
