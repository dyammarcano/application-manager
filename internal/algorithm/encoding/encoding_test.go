package encoding

import (
	"github.com/dyammarcano/application-manager/internal/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mm = "{\n  \"id\": 1,\n  \"first_name\": \"Trace\",\n  \"last_name\": \"Morena\",\n  \"email\": \"tmorena0@behance.net\",\n  \"gender\": \"Male\",\n  \"ip_address\": \"79.238.201.131\",\n  \"rfid\": \"6528560a9037d5962c0bc86f\",\n  \"index\": 1,\n  \"guid\": \"f98c9b31-42cd-4491-b733-56792707b53c\",\n  \"isActive\": false,\n  \"balance\": \"$2,537.69\",\n  \"picture\": \"https://placehold.it/32x32\",\n  \"age\": 25,\n  \"eyeColor\": \"green\",\n  \"name\": \"Lucia Molina\",\n  \"company\": \"BUZZMAKER\",\n  \"phone\": \"+1 (821) 517-2438\",\n  \"address\": \"770 Blake Avenue, Washington, Delaware, 2378\",\n  \"about\": \"Sit dolor eu magna ea et. Exercitation consequat aute eiusmod adipisicing reprehenderit. Ea ullamco ex minim deserunt voluptate qui ad sint Lorem voluptate exercitation. In non sunt laboris ad aliqua labore ex laboris laboris nostrud aliquip. Proident et incididunt id mollit ut cupidatat enim adipisicing veniam anim ea minim. Ut cupidatat deserunt dolor reprehenderit do pariatur ex do occaecat incididunt commodo.\\r\\n\",\n  \"registered\": \"2016-08-21T06:09:02 +03:00\",\n  \"latitude\": 7.273269,\n  \"longitude\": -1.804389,\n  \"data1\": {\n    \"title\": \"{{faker 'lorem.sentence'}}\",\n    \"content\": \"{{faker 'lorem.sentences'}}\",\n    \"media\": \"{{faker 'image.nature'}}\",\n    \"author\": {\n      \"name\": \"{{faker 'name.firstName'}} {{faker 'name.firstName'}}\",\n      \"avatar\": \"{{faker 'image.avatar'}}\"\n    },\n    \"comments\": {\n      \"id\": \"{{faker 'datatype.uuid'}}\",\n      \"content\": \"{{faker 'lorem.sentence'}}\",\n      \"author\": {\n        \"name\": \"{{faker 'name.firstName'}} {{faker 'name.firstName'}}\",\n        \"avatar\": \"{{faker 'image.avatar'}}\"\n      }\n    }\n  },\n  \"data2\": {\n    \"title\": \"{{faker 'lorem.sentence'}}\",\n    \"content\": \"{{faker 'lorem.sentences'}}\",\n    \"media\": \"{{faker 'image.nature'}}\",\n    \"author\": {\n      \"name\": \"{{faker 'name.firstName'}} {{faker 'name.firstName'}}\",\n      \"avatar\": \"{{faker 'image.avatar'}}\"\n    },\n    \"comments\": {\n      \"id\": \"{{faker 'datatype.uuid'}}\",\n      \"content\": \"{{faker 'lorem.sentence'}}\",\n      \"author\": {\n        \"name\": \"{{faker 'name.firstName'}} {{faker 'name.firstName'}}\",\n        \"avatar\": \"{{faker 'image.avatar'}}\"\n      }\n    }\n  }\n}\n"

func TestEncoding1kChars(t *testing.T) {
	serialized, err := Serialize(mm)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mm)
}

func TestEncoding2kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message2kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message2kChars)
}

func TestEncoding3kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message3kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message3kChars)
}

func TestEncoding4kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message4kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message4kChars)
}

func TestEncoding5kChars(t *testing.T) {
	serialized, err := Serialize(mock.Message5kChars)
	assert.Nil(t, err)

	deserialized, err := Deserialize(serialized)
	assert.Nil(t, err)

	assert.Equal(t, deserialized, mock.Message5kChars)
}
