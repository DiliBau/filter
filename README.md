# filter


## Example
```go
import "github.com/DiliBau/filter"

filter := filter.NewFilter()


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

filter.Register(reflect.TypeOf(&Model1{}), func(value reflect.Value) error {
    // add conditions or custom code here
    value.Elem().FieldByName("SomeProp").SetString("new-prop")

    return nil
})

filter.Register(reflect.TypeOf(&Model2{}), func(value reflect.Value) error {
    // add conditions or custom code here
    id := value.Elem().FieldByName("Id")
    id.SetString("id-" + id.String())

    return nil
})


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

err := filter.Apply(reflect.ValueOf(&res))
// all instances of &Model1 now have SomeProp = new-prod
// all instances of &Model2 now have Id = "id-" + old value of Id
```

## Notes

- changes will be applied until the first error occur, they will not be reverted on error
