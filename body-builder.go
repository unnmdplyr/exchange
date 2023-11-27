package main

type BodyBuilder interface {
	build(b *Body)
}

type ToJsonFormBuilder struct {
}

type ToFormUrlEncodedFormBuilder struct {
}

func (builder *ToJsonFormBuilder) build(b *Body) {

}
