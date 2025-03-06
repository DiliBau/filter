package filter_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/DiliBau/filter"
)

type FilterSuite struct {
	suite.Suite
	filter *filter.Filter
}

type Model1 struct {
	SomeProp string
}

type Model2 struct {
	Id string
}

type Model3 struct {
	Ones []Model1
	Twos []Model2
}

func (f *FilterSuite) SetupTest() {
	f.filter = filter.NewFilter()

	// remove TfmJourneyKey
	f.filter.Register(reflect.TypeOf(&Model1{}), func(value reflect.Value) error {
		value.Elem().FieldByName("SomeProp").SetString("new-prop")

		return nil
	})

	// add id- prefix to route
	f.filter.Register(reflect.TypeOf(&Model2{}), func(value reflect.Value) error {
		id := value.Elem().FieldByName("Id")
		id.SetString("id-" + id.String())

		return nil
	})
}

func (f *FilterSuite) TestApply() {
	res := Model3{
		Ones: []Model1{
			{SomeProp: "some-prop"},
			{SomeProp: "other-prop"},
		},
		Twos: []Model2{
			{Id: "some-id"},
			{Id: "other-id"},
		},
	}

	err := f.filter.Apply(reflect.ValueOf(res.Ones[0]))
	f.Error(err)

	err = f.filter.Apply(reflect.ValueOf(&res))
	f.NoError(err)

	f.Equal("new-prop", res.Ones[0].SomeProp)
	f.Equal("id-some-id", res.Twos[0].Id)
	f.Equal("new-prop", res.Ones[1].SomeProp)
	f.Equal("id-other-id", res.Twos[1].Id)

	res2 := []interface{}{
		&Model1{SomeProp: "some-prop"},
		&Model1{SomeProp: "other-prop"},
		&Model2{Id: "some-id"},
		&Model2{Id: "other-id"},
	}

	// derefencing should trigger error
	err = f.filter.Apply(reflect.ValueOf(*res2[0].(*Model1)))
	f.Error(err)

	err = f.filter.Apply(reflect.ValueOf(&res2))
	f.NoError(err)

	f.Equal("new-prop", res2[0].(*Model1).SomeProp)
	f.Equal("id-some-id", res2[2].(*Model2).Id)
	f.Equal("new-prop", res2[1].(*Model1).SomeProp)
	f.Equal("id-other-id", res2[3].(*Model2).Id)
}

func TestFilterSuite(t *testing.T) {
	suite.Run(t, new(FilterSuite))
}
