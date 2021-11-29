package main

import (
	`bytes`
	`encoding/gob`
	`encoding/json`
	`encoding/xml`
	`fmt`
	`os`
	`strings`
)

/*
	在golang中运行时创建的对象和通用数据结构(XML, JSON等)之间转换
	编码：将对象转换为字节数组，从而可以写入流；
	解码：将字节数组还原为对象
*/

// TODO: OpenFile返回的是一个IO流（更具体一点，是一个文件流），将encoder和decoder解读为套接在流上的"奶嘴"，方便数据的存取

type Address struct {
	Type    string
	City    string
	Country string
}

type VCard struct {
	Firstname string
	Lastname  string
	Addresses []*Address
	Remark    string
}

// Encoding2Json 将对象以json格式编码为字节数组/以json格式写入文本文件
func Encoding2Json() {
	pa := &Address{"home", "Hangzhou", "China"}
	wa := &Address{"work", "Boom", "Belgium"}
	vc := VCard{"Narcissus", "Zhao", []*Address{pa, wa}, "A genius"}

	// TODO：以json格式将内存中的变量vc编码为字节数组
	js, _ := json.Marshal(vc)
	fmt.Printf("JSON format:\n%s\n\n", js)

	// TODO：将字节数组js以json的格式解码到内存对象v中
	v := VCard{}
	_ = json.Unmarshal(js, &v)
	fmt.Println(v)

	// 把内存中的对象vc以json格式写入文本文件
	file, _ := os.OpenFile("./coding/vcard.json", os.O_WRONLY|os.O_CREATE, 0644)
	enc := json.NewEncoder(file)
	// 将内存中的变量vc的JSON编码写入自身（即file，是一个流；而enc是一个套在流出入口上的一个"奶嘴"）
	_ = enc.Encode(vc)
	// 流的关闭可以放到defer中执行
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("error %v happend when closing file\n", err)
			os.Exit(1)
		}
	}()
}

// DecodingFromJson 将json字符串解码为golang中的对象
func DecodingFromJson() {
	b := []byte(`{"Name": "Wednesday", "Age": 6, "Parents": ["Gomez", "Lily"]}`)
	var f interface{} // f被声明为一个接口，因此可以被赋予任何值

	// 将json格式的字节数组解码到对象f中
	_ = json.Unmarshal(b, &f)
	// TODO：使用类型判别来访问f（将f在内存中以类型map[string]interface{}的一个实例来解读）
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Printf("%s: %s\n", k, vv)
		case float64:
			fmt.Printf("%s: %f\n", k, vv)
		case []interface{}:
			fmt.Printf("%s: ", k)
			for i, u := range vv {
				fmt.Printf("(%d, %v) ", i, u)
			}
			fmt.Println()
		default:
			fmt.Println(k, "is of a type which I don' know")
		}
	}
}

type VCard2 struct {
	Name      string
	Age       int
	Addresses []*Address
	Remark    string
}

func DecodingFromJson2() {
	byteArr := []byte(`{"Name":"Julia Zhao","Age":24,"Addresses":[{"Type":"work","City":"hangzhou","Country":"China"},{"Type":"home","City":"nanjing","Country":"China"}],"Remark":"this is me"}`)
	vc := VCard2{}
	// TODO：直接解码到一个VCard2实例中
	_ = json.Unmarshal(byteArr, &vc)
	fmt.Printf("Name: %v\n", vc.Name)
	fmt.Printf("Age: %v\n", vc.Age)
	fmt.Printf("Address: ")
	for _, addr := range vc.Addresses {
		fmt.Printf("%v ", *addr)
	}
	fmt.Println()
	fmt.Printf("Remark: %v\n", vc.Remark)
}

// DecodingFromXML 将XML字符串解码为golang中的对象
func DecodingFromXML() {
	input := "<Person><FirstName>Laura</FirstName><LastName>Lynn</LastName></Person>"
	inputReader := strings.NewReader(input)
	p := xml.NewDecoder(inputReader)

	for t, err := p.Token(); err == nil; t, err = p.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			fmt.Printf("Token name: %s\n", name)
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				fmt.Printf("An attribute is: %s %s\n", attrName, attrValue)
			}
		case xml.EndElement:
			fmt.Println("End of token")
		case xml.CharData:
			content := string(token)
			fmt.Printf("This is the content: %s\n", content)
		default:
			fmt.Println("Cannot understand")
		}
	}
}

/* Gob: like pickle (for python) and Serialization (for java),
used for data (param & result) transmission, typically for RPCs. */

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func UseGob() {
	// "自产自销"
	var net bytes.Buffer
	enc := gob.NewEncoder(&net)
	_ = enc.Encode(P{3, 4, 5, "Pythagoras"}) // write to net

	dec := gob.NewDecoder(&net)
	q := Q{}
	_ = dec.Decode(&q) // decoded from net
	fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)
}

// Send the client send request data
func Send() bytes.Buffer {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)

	p := P{
		X:    3,
		Y:    4,
		Z:    5,
		Name: "data transmitted in network",
	}
	_ = enc.Encode(&p)
	return data
}

// Receive the server receives the request and processes
func Receive(buffer bytes.Buffer) Q {
	dec := gob.NewDecoder(&buffer)
	q := Q{}
	_ = dec.Decode(&q)
	fmt.Printf("%d, %d, %s", *q.X, *q.Y, q.Name)
	return q
}

func UseGob2() {
	buffer := Send()
	_ = Receive(buffer)
}

func main() {
	// Encoding2Json()
	// DecodingFromJson()
	// DecodingFromJson2()
	// DecodingFromXML()
	// UseGob()
	UseGob2()
}
