package swagger

import "github.com/bytedance/sonic"

var (
	jsonEnc = sonic.Config{
		EscapeHTML: false,
	}.Froze()
)

func JsonMarshal(v interface{}) ([]byte, error) {
	return jsonEnc.Marshal(v)
}

func JsonUnmarshal(data []byte, v interface{}) error {
	return jsonEnc.Unmarshal(data, v)
}
