package neo4jx

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type BaseModel struct {
	Driver    neo4j.Driver
	Type      reflect.Type
	dsn       string
	LabelName string
	jsonTags  []string
}

func NewBaseModel(dsn string, data any) (*BaseModel, error) {
	v, e := GetDriverFromPool(dsn)
	if e != nil {
		return nil, e
	}
	label, e := ToLabelName(data)
	if e != nil {
		return nil, e
	}

	t := reflect.TypeOf(data)
	model := &BaseModel{
		Driver:    v,
		LabelName: label,
		dsn:       dsn,
		Type:      t,
	}

	if model.Type.Kind() == reflect.Ptr {
		return nil, errors.New("data must be struct type")
	}

	// check fields
	indexes := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if i == 0 {
			switch field.Type.Kind() {
			case reflect.Uint,
				reflect.Uint64,
				reflect.Uint32,
				reflect.Uint16,
				reflect.String:
			default:
				return nil, errors.New("The first field " + field.Name + "'s type must be one of uint,uint32,uint64,uint16,string")
			}
		}

		//dbTag
		dbTag, ok := field.Tag.Lookup(TagName)
		if !ok {
			return nil, errors.New("field " + field.Name + " has no `" + TagName + "` tag specified")
		}
		if dbTag != strcase.ToLowerCamel(dbTag) {
			return nil, errors.New("Field '" + field.Name + "'s `" + TagName + "` tag is not in low camel case")
		}

		//index
		if index, ok := field.Tag.Lookup("index"); ok {
			indexes[dbTag] = index
		}

		model.jsonTags = append(model.jsonTags, dbTag)
	}

	return model, nil
}

func MustNewBaseModel(dsn string, data any) *BaseModel {
	v, e := NewBaseModel(dsn, data)
	if e != nil {
		log.Fatal(e)
	}
	return v
}

func (b *BaseModel) FindBy(key string, value any) (any, error) {
	query := fmt.Sprintf("MATCH (n:%s{%s:$value}) return n", b.LabelName, key)

	ses := b.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer ses.Close()
	println(query)
	res, e := ses.Run(query, map[string]any{
		"value": value,
	})
	if e != nil {
		log.Println(e)
		return nil, e
	}

	record, e := res.Single()
	if e != nil {
		log.Println(e)
		return nil, e
	}
	v := reflect.New(b.Type)
	e = UnmarshalRecord(v.Interface(), record, "n")
	if e != nil {
		log.Println(e)
		return nil, e
	}

	return v.Interface(), nil
}
