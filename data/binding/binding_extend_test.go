package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustBinding(t *testing.T) {
	type someStruct struct {
		Exposed   string
		unexposed bool
	}
	structVal := someStruct{
		Exposed:   "I'm an exposed string",
		unexposed: true,
	}
	sbind := BindStruct(&structVal)

	assert.NotPanics(t, func() {
		boundExposed := sbind.MustGetItem("Exposed").(String)
		assert.Equal(t, "I'm an exposed string", boundExposed.MustGet())
	}, "Exposed fields are readable by reflection, and thus bindable")
	assert.NotPanics(t, func() {
		assert.Equal(t, "I'm an exposed string", sbind.MustGetValue("Exposed").(string))
	}, "MustGetValue gets the struct field value directly")
	assert.Panics(t, func() {
		_ = sbind.MustGetItem("unexposed")
	}, "Unexposed fields should not be readable by reflection")
}

func TestBindUserType(t *testing.T) {
	val := user{name: "Unnamed"}
	u := bindUserType(&val)
	v, err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Unnamed", v.String())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	u.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = u.Set(user{name: "Replace"})
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, "User: Replace", val.String())
	assert.True(t, called)

	called = false
	val = user{name: "Direct"}
	_ = u.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Direct", v.String())

	called = false
	val.name = "FieldSet"
	_ = u.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: FieldSet", v.String())
}

func TestNewUserType(t *testing.T) {
	u := newUserType()
	v, err := u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: ", v.String())

	err = u.Set(user{name: "Dave"})
	assert.Nil(t, err)
	v, err = u.Get()
	assert.Nil(t, err)
	assert.Equal(t, "User: Dave", v.String())
}

type user struct {
	name string
}

func (u *user) String() string {
	return "User: " + u.name
}

type userType struct {
	Untyped
}

func newUserType() *userType {
	ret := &userType{Untyped: NewUntyped()}
	_ = ret.Set(user{})
	return ret
}

func (t *userType) Get() (user, error) {
	val, err := t.Untyped.Get()
	return val.(user), err
}

func (t *userType) Set(u user) error {
	return t.Untyped.Set(u)
}

type externalUserType struct {
	ExternalUntyped
}

func bindUserType(u *user) *externalUserType {
	return &externalUserType{ExternalUntyped: BindUntyped(u)}
}

func (t *externalUserType) Get() (user, error) {
	val, err := t.ExternalUntyped.Get()
	return val.(user), err
}

func (t *externalUserType) Set(u user) error {
	return t.ExternalUntyped.Set(u)
}
