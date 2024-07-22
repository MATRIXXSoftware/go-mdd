package cmdc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

var codec = NewCodec()

func sampleData1() string {
	return "<1,18,0,-6,5222,2>[1,20,<1,2,0,452,5222,2>[100],4]"
}

func TestDecodeSampleData1(t *testing.T) {
	data := sampleData1()
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	t.Logf("decoded: %s", decoded.Dump())

	// Assert decoded
	assert.Equal(t, 1, len(decoded.Containers))
	assert.Equal(t, 1, decoded.Containers[0].Header.Version)
	assert.Equal(t, 18, decoded.Containers[0].Header.TotalField)
	assert.Equal(t, 0, decoded.Containers[0].Header.Depth)
	assert.Equal(t, -6, decoded.Containers[0].Header.Key)
	assert.Equal(t, 5222, decoded.Containers[0].Header.SchemaVersion)
	assert.Equal(t, 2, decoded.Containers[0].Header.ExtVersion)

	assert.Equal(t, 4, len(decoded.Containers[0].Fields))
	assert.Equal(t, "1", decoded.Containers[0].GetField(0).String())
	assert.Equal(t, "20", decoded.Containers[0].GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[100]", decoded.Containers[0].GetField(2).String())
	assert.Equal(t, "4", decoded.Containers[0].GetField(3).String())
}

func sampleData2() string {
	return "<1,8,0,-6,5222,2>[,,2,(5:AMF-1),(4:eMBB),(11:SouthWestUK),1]<1,1,0,-5,5222,2>[1000001]<1,7,0,263,5222,2>[2,{<1,5,1,330,5222,2>[4,200,1,17485760.0,17485824.0]},(6:555555),0,0.0,64.0,200]<1,11,0,626,5222,2>[{1,3,1},,{<1,17,1,624,5222,2>[17485824.0,200,1,(21:Data: Asset + Overage),0:1:5:277,(7:2000000),3,0,,,,,,,,0]},,{(13:HXS0:1:52:409)},{<1,8,1,1000,5222,2>[4,(18:Triple Play Bundle),1,0,,,,0]},{<1,5,1,1277,5222,2>[4,(17:999 - 1200TB Plan),1,<1,6,1,-11,5222,2>[800000.0,1200.0,300000.0,5000000000.0,100000.0,5000000000.0]<1,0,0,1257,5222,2>[]]},,{<1,3,1,1360,5222,2>[,(5:Usage),(5:Usage)]},{<1,2,1,627,5222,2>[(13:HXS0:1:52:408),1]<1,29,0,208,5222,2>[0:1:5:279,(7:1000001),0:1:5:283,,4,0:1:5:278,0:1:5:277,(7:2000000),,,,{<1,14,1,209,5222,2>[,1000,2,1,4,2021-09-07T08:00:25.000000Z,2021-10-07T08:00:25.000000Z,1,0.0,,,-1258291032.242187,,0]},{<1,12,1,567,5222,2>[17485824.0,200,0,0,1,0.0,,1,,1,0]},2021-09-09T16:37:19.000000Z,,(13:HXS0:1:52:409),,,,,1,,,0:1:5:281,,1,2]}]<1,29,0,208,5222,2>[0:1:5:279,(7:1000001),0:1:5:283,,0,0:1:5:280,0:1:5:279,(7:1000001),,,,,,2021-09-09T16:37:19.000000Z,0,(13:HXS0:1:52:408),,,,,1,,,0:1:5:281,,1,1,(26:00000000000000594134:00000)]"
}

func TestDecodeSampleData2(t *testing.T) {
	data := sampleData2()
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	t.Logf("decoded: %s", decoded.Dump())

	// Assert decoded
	assert.Equal(t, 5, len(decoded.Containers))
	assert.Equal(t, -6, decoded.Containers[0].Header.Key)
	assert.Equal(t, -5, decoded.Containers[1].Header.Key)
	assert.Equal(t, 263, decoded.Containers[2].Header.Key)
	assert.Equal(t, 626, decoded.Containers[3].Header.Key)
	assert.Equal(t, 208, decoded.Containers[4].Header.Key)

	assert.Equal(t, 7, len(decoded.Containers[0].Fields))
	assert.Equal(t, 1, len(decoded.Containers[1].Fields))
	assert.Equal(t, 7, len(decoded.Containers[2].Fields))
	assert.Equal(t, 10, len(decoded.Containers[3].Fields))
	assert.Equal(t, 28, len(decoded.Containers[4].Fields))

	assert.Equal(t, "", decoded.Containers[0].GetField(0).String())
	assert.Equal(t, "", decoded.Containers[0].GetField(1).String())
	assert.Equal(t, "2", decoded.Containers[0].GetField(2).String())
	assert.Equal(t, "(5:AMF-1)", decoded.Containers[0].GetField(3).String())
	assert.Equal(t, "(4:eMBB)", decoded.Containers[0].GetField(4).String())
	assert.Equal(t, "(11:SouthWestUK)", decoded.Containers[0].GetField(5).String())
	assert.Equal(t, "1", decoded.Containers[0].GetField(6).String())

	assert.Equal(t, "1000001", decoded.Containers[1].GetField(0).String())

	assert.Equal(t, "2", decoded.Containers[2].GetField(0).String())
	assert.Equal(t, "{<1,5,1,330,5222,2>[4,200,1,17485760.0,17485824.0]}", decoded.Containers[2].GetField(1).String())
	assert.Equal(t, "(6:555555)", decoded.Containers[2].GetField(2).String())
	assert.Equal(t, "0", decoded.Containers[2].GetField(3).String())
	assert.Equal(t, "0.0", decoded.Containers[2].GetField(4).String())
	assert.Equal(t, "64.0", decoded.Containers[2].GetField(5).String())
	assert.Equal(t, "200", decoded.Containers[2].GetField(6).String())

	assert.Equal(t, "{1,3,1}", decoded.Containers[3].GetField(0).String())
	assert.Equal(t, "", decoded.Containers[3].GetField(1).String())
	assert.Equal(t, "{<1,17,1,624,5222,2>[17485824.0,200,1,(21:Data: Asset + Overage),0:1:5:277,(7:2000000),3,0,,,,,,,,0]}", decoded.Containers[3].GetField(2).String())
	assert.Equal(t, "", decoded.Containers[3].GetField(3).String())
	assert.Equal(t, "{(13:HXS0:1:52:409)}", decoded.Containers[3].GetField(4).String())
	assert.Equal(t, "{<1,8,1,1000,5222,2>[4,(18:Triple Play Bundle),1,0,,,,0]}", decoded.Containers[3].GetField(5).String())
	assert.Equal(t, "{<1,5,1,1277,5222,2>[4,(17:999 - 1200TB Plan),1,<1,6,1,-11,5222,2>[800000.0,1200.0,300000.0,5000000000.0,100000.0,5000000000.0]<1,0,0,1257,5222,2>[]]}", decoded.Containers[3].GetField(6).String())
	assert.Equal(t, "", decoded.Containers[3].GetField(7).String())
	assert.Equal(t, "{<1,3,1,1360,5222,2>[,(5:Usage),(5:Usage)]}", decoded.Containers[3].GetField(8).String())
	assert.Equal(t, "{<1,2,1,627,5222,2>[(13:HXS0:1:52:408),1]<1,29,0,208,5222,2>[0:1:5:279,(7:1000001),0:1:5:283,,4,0:1:5:278,0:1:5:277,(7:2000000),,,,{<1,14,1,209,5222,2>[,1000,2,1,4,2021-09-07T08:00:25.000000Z,2021-10-07T08:00:25.000000Z,1,0.0,,,-1258291032.242187,,0]},{<1,12,1,567,5222,2>[17485824.0,200,0,0,1,0.0,,1,,1,0]},2021-09-09T16:37:19.000000Z,,(13:HXS0:1:52:409),,,,,1,,,0:1:5:281,,1,2]}", decoded.Containers[3].GetField(9).String())

	assert.Equal(t, "0:1:5:279", decoded.Containers[4].GetField(0).String())
	assert.Equal(t, "(7:1000001)", decoded.Containers[4].GetField(1).String())
	assert.Equal(t, "0:1:5:283", decoded.Containers[4].GetField(2).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(3).String())
	assert.Equal(t, "0", decoded.Containers[4].GetField(4).String())
	assert.Equal(t, "0:1:5:280", decoded.Containers[4].GetField(5).String())
	assert.Equal(t, "0:1:5:279", decoded.Containers[4].GetField(6).String())
	assert.Equal(t, "(7:1000001)", decoded.Containers[4].GetField(7).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(8).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(9).String())
	assert.Equal(t, "2021-09-09T16:37:19.000000Z", decoded.Containers[4].GetField(13).String())
	assert.Equal(t, "0", decoded.Containers[4].GetField(14).String())
	assert.Equal(t, "(13:HXS0:1:52:408)", decoded.Containers[4].GetField(15).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(16).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(17).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(18).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(19).String())
	assert.Equal(t, "1", decoded.Containers[4].GetField(20).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(21).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(22).String())
	assert.Equal(t, "0:1:5:281", decoded.Containers[4].GetField(23).String())
	assert.Equal(t, "", decoded.Containers[4].GetField(24).String())
	assert.Equal(t, "1", decoded.Containers[4].GetField(25).String())
	assert.Equal(t, "1", decoded.Containers[4].GetField(26).String())
	assert.Equal(t, "(26:00000000000000594134:00000)", decoded.Containers[4].GetField(27).String())
}

func sampleData3() string {
	return "<1,8,0,-6,5222,2>[,,,(5:AMF-1),(4:eMBB),(11:SouthWestUK),1]<1,1,0,-5,5222,2>[1000001]<1,7,0,263,5222,2>[2,{<1,5,1,330,5222,2>[200,4,1],<1,5,1,330,5222,2>[202,6,2]},(6:555555)]<1,11,0,626,5222,2>[{1,3,1},,{<1,17,1,624,5222,2>[17485824.0,200,1,(21:Data: Asset + Overage),0:1:5:277,(7:2000000),3,0,,,,,,,,0]},0:1:5:144,{(13:HXS0:1:52:409)}]<1,29,0,208,5222,2>[]"
}

func TestDecodeSampleData3(t *testing.T) {
	data := sampleData3()
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	t.Logf("decoded: %s", decoded.Dump())

	// Assert decoded
	assert.Equal(t, 5, len(decoded.Containers))
	assert.Equal(t, -6, decoded.Containers[0].Header.Key)
	assert.Equal(t, -5, decoded.Containers[1].Header.Key)
	assert.Equal(t, 263, decoded.Containers[2].Header.Key)
	assert.Equal(t, 626, decoded.Containers[3].Header.Key)
	assert.Equal(t, 208, decoded.Containers[4].Header.Key)

	assert.Equal(t, "1000001", decoded.Containers[1].GetField(0).String())

	assert.Equal(t, "2", decoded.Containers[2].GetField(0).String())
	assert.Equal(t, "{<1,5,1,330,5222,2>[200,4,1],<1,5,1,330,5222,2>[202,6,2]}", decoded.Containers[2].GetField(1).String())
	assert.Equal(t, "(6:555555)", decoded.Containers[2].GetField(2).String())

	assert.Equal(t, "{1,3,1}", decoded.Containers[3].GetField(0).String())
	assert.Equal(t, "", decoded.Containers[3].GetField(1).String())
	assert.Equal(t, "{<1,17,1,624,5222,2>[17485824.0,200,1,(21:Data: Asset + Overage),0:1:5:277,(7:2000000),3,0,,,,,,,,0]}", decoded.Containers[3].GetField(2).String())
	assert.Equal(t, "0:1:5:144", decoded.Containers[3].GetField(3).String())
	assert.Equal(t, "{(13:HXS0:1:52:409)}", decoded.Containers[3].GetField(4).String())
}

func TestDecodeExample1(t *testing.T) {
	data := "<1,18,0,-6,5222,2>[1,,-20,(5:value),{10,20},1-5-1-5,*1234567890#]"
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	decoded.Containers[0].LoadDefinition(&dictionary.ContainerDefinition{
		Fields: []dictionary.FieldDefinition{
			{Type: field.UInt8},
			{Type: field.UInt32},
			{Type: field.Int32},
			{Type: field.String},
			{Type: field.Int32, IsMulti: true},
			{Type: field.ObjectID},
			{Type: field.PhoneNo},
		},
	})

	v, err := decoded.Containers[0].GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint8(1), v)

	v, err = decoded.Containers[0].GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, nil, v)

	v, err = decoded.Containers[0].GetField(2).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-20), v)

	v, err = decoded.Containers[0].GetField(3).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, "value", v)

	v, err = decoded.Containers[0].GetField(4).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, []int32{10, 20}, v)

	v, err = decoded.Containers[0].GetField(5).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, field.MtxObjectID("1-5-1-5"), v)

	v, err = decoded.Containers[0].GetField(6).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, field.MtxPhoneNo("*1234567890#"), v)
}

func TestEncodeExample(t *testing.T) {
	objectIDValue, _ := field.NewObjectID([]byte("1-5-1-5"))
	phoneNoValue, _ := field.NewPhoneNo([]byte("*1234567890#"))
	containers := mdd.Containers{
		Containers: []mdd.Container{
			{
				Header: mdd.Header{Version: 1, TotalField: 18, Depth: 0, Key: -6, SchemaVersion: 5222, ExtVersion: 2},
				Fields: []mdd.Field{
					*mdd.NewBasicField(uint8(1)),
					*mdd.NewNullField(field.UInt32),
					*mdd.NewBasicField(int32(-20)),
					*mdd.NewBasicField("value"),
					*mdd.NewBasicListField([]int32{10, 20}),
					{Type: field.ObjectID, Value: objectIDValue},
					{Type: field.PhoneNo, Value: phoneNoValue},
				},
			},
		},
	}

	expected := "<1,18,0,-6,5222,2>[1,,-20,(5:value),{10,20},1-5-1-5,*1234567890#]"
	encoded, err := codec.Encode(&containers)

	assert.Nil(t, err)
	assert.Equal(t, []byte(expected), encoded)
}

func TestDecodeEncode(t *testing.T) {

	data := "<1,18,0,-6,5222,2>[1,-20,<1,2,0,452,5222,2>[(14:abcdefghijklmn),,100,5.8888,<1,3,0,400,5222,2>[]],,400000000]"
	decoded, err := codec.Decode([]byte(data))
	assert.Nil(t, err)

	// Print decoded decoded
	t.Logf("decoded: %s", decoded.Dump())

	// Assert decoded
	assert.Equal(t, "1", decoded.Containers[0].GetField(0).String())
	assert.Equal(t, "-20", decoded.Containers[0].GetField(1).String())
	assert.Equal(t, "<1,2,0,452,5222,2>[(14:abcdefghijklmn),,100,5.8888,<1,3,0,400,5222,2>[]]", decoded.Containers[0].GetField(2).String())
	assert.Equal(t, "", decoded.Containers[0].GetField(3).String())
	assert.Equal(t, "400000000", decoded.Containers[0].GetField(4).String())

	// Mark field types
	decoded.Containers[0].LoadDefinition(&dictionary.ContainerDefinition{
		Fields: []dictionary.FieldDefinition{
			{Type: field.UInt8},
			{Type: field.Int32},
			{Type: field.Struct},
			{Type: field.UInt32},
			{Type: field.UInt64},
		},
	})

	// Retrieve field 0 as uint8
	field0, err := decoded.Containers[0].GetField(0).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint8(1), field0)

	// Retrieve field 1 as int32
	field1, err := decoded.Containers[0].GetField(1).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, int32(-20), field1)

	// Retrieve field 2 as struct
	field2, err := decoded.Containers[0].GetField(2).GetValue()
	assert.Nil(t, err)

	{
		nested := field2.(*mdd.Containers)
		t.Logf("nested: %s", nested.Dump())

		// Assert Nested
		assert.Equal(t, "(14:abcdefghijklmn)", nested.Containers[0].GetField(0).String())
		assert.Equal(t, "", nested.Containers[0].GetField(1).String())
		assert.Equal(t, "100", nested.Containers[0].GetField(2).String())
		assert.Equal(t, "5.8888", nested.Containers[0].GetField(3).String())
		assert.Equal(t, "<1,3,0,400,5222,2>[]", nested.Containers[0].GetField(4).String())

		// Mark field types
		nested.Containers[0].LoadDefinition(&dictionary.ContainerDefinition{
			Fields: []dictionary.FieldDefinition{
				{Type: field.String},
				{Type: field.String},
				{Type: field.UInt32},
				{Type: field.Decimal},
				{Type: field.Struct},
			},
		})

		// Retrieve nested field 0 as string
		field0, err := nested.Containers[0].GetField(0).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, "abcdefghijklmn", field0)

		// Retrieve nested field 1 as null
		field1, err := nested.Containers[0].GetField(1).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, nil, field1)

		// Retrieve nested field 2 as UInt32
		field2, err := nested.Containers[0].GetField(2).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, uint32(100), field2)

		// Retrieve nested field 3 as Decimal
		field3, err := nested.Containers[0].GetField(3).GetValue()
		assert.Nil(t, err)
		assert.Equal(t, "5.8888", field3.(*big.Float).Text('f', -1))

		// Retrieve nested field 4 as Struct
		field4, err := nested.Containers[0].GetField(4).GetValue()
		assert.Nil(t, err)

		{
			nested2 := field4.(*mdd.Containers)
			t.Logf("nested2: %s", nested2.Dump())

			// Assert Nested2
			assert.Equal(t, "", nested2.Containers[0].GetField(0).String())
			// Field is null
			assert.Equal(t, true, nested2.Containers[0].GetField(0).IsNull)

			field0, err := nested2.Containers[0].GetField(0).GetValue()
			assert.Nil(t, err)
			assert.Equal(t, nil, field0)
		}
	}

	// Retrieve field 3 as uint32
	field3, err := decoded.Containers[0].GetField(3).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, nil, field3)

	// Retrieve field 4 as uint64
	field4, err := decoded.Containers[0].GetField(4).GetValue()
	assert.Nil(t, err)
	assert.Equal(t, uint64(400000000), field4)

	// Re - encode
	encoded, err := codec.Encode(decoded)
	assert.Nil(t, err)
	assert.Equal(t, data, string(encoded))
}
