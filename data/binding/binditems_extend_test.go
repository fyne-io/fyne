package binding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBindTime(t *testing.T) {
	val := time.Now()
	f := bindTime(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, val.Unix(), v.Unix())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	newTime := val.Add(time.Hour)
	called = false
	err = f.Set(newTime)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, newTime.Unix(), val.Unix())
	assert.True(t, called)

	newTime = newTime.Add(time.Minute)
	called = false
	val = newTime
	_ = f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, newTime.Unix(), v.Unix())
}

func TestNewTime(t *testing.T) {
	f := newTime()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, time.Unix(0, 0), v)

	now := time.Now()
	err = f.Set(now)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, now.Unix(), v.Unix())
}

type timeBinding struct {
	Int
	src *time.Time
}

func bindTime(t *time.Time) *timeBinding {
	return &timeBinding{Int: NewInt(), src: t}
}

func newTime() *timeBinding {
	return &timeBinding{Int: NewInt()}
}

func (t *timeBinding) Get() (time.Time, error) {
	if t.src != nil {
		return *t.src, nil
	}

	i, err := t.Int.Get()
	return time.Unix(int64(i), 0), err
}

func (t *timeBinding) Reload() error {
	if t.src == nil {
		return nil
	}

	return t.Set(*t.src)
}

func (t *timeBinding) Set(time time.Time) error {
	if t.src != nil {
		*t.src = time
	}

	i := time.Unix()
	return t.Int.Set(int(i))
}

func TestBindUserType(t *testing.T) {
	val := user{name: "Unnamed"}
	u := bindUserType(&val)
	v, err := u.GetUser()
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
	v, err = u.GetUser()
	assert.Nil(t, err)
	assert.Equal(t, "User: Direct", v.String())
}

func TestNewUserType(t *testing.T) {
	u := newUserType()
	v, err := u.GetUser()
	assert.Nil(t, err)
	assert.Equal(t, "User: ", v.String())

	err = u.Set(user{name: "Dave"})
	assert.Nil(t, err)
	v, err = u.GetUser()
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

func (t *userType) GetUser() (user, error) {
	val, err := t.Get()
	return val.(user), err
}

func (t *userType) SetUser(u user) error {
	return t.Set(u)
}

type externalUserType struct {
	*userType
	*user
}

func bindUserType(u *user) *externalUserType {
	return &externalUserType{userType: newUserType(), user: u}
}

func (t *externalUserType) GetUser() (user, error) {
	val, err := t.Get()
	return val.(user), err
}

func (t *externalUserType) SetUser(u user) error {
	return t.Set(u)
}

func (t *externalUserType) Get() (interface{}, error) {
	return *t.user, nil
}

func (t *externalUserType) Reload() error {
	return t.Untyped.Set(*t.user)
}

func (t *externalUserType) Set(u interface{}) error {
	*t.user = u.(user)
	return t.userType.Set(u)
}
