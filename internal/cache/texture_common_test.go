package cache

import "testing"

func TestDropTexturesFor(t *testing.T) {
	c := &dummyCanvas{}
	obj := &dummyWidget{}
	SetTexture(obj, TextureType(1), c)
	if _, ok := GetTexture(obj); !ok {
		t.Fatal("texture should be cached")
	}
	DropTexturesFor(c)
	if _, ok := GetTexture(obj); ok {
		t.Fatal("DropTexturesFor must drop object textures without GL free")
	}
}
