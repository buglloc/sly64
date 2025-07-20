package upstream

import (
	"encoding"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
)

var _ yaml.InterfaceUnmarshaler = (*Kind)(nil)
var _ yaml.InterfaceMarshaler = (*Kind)(nil)
var _ json.Unmarshaler = (*Kind)(nil)
var _ json.Marshaler = (*Kind)(nil)
var _ encoding.TextUnmarshaler = (*Kind)(nil)

type Kind uint8

const (
	KindNone Kind = iota
	KindPlain
	KindDoT
)

func (k Kind) String() string {
	switch k {
	case KindNone:
		return "none"
	case KindPlain:
		return "plain"
	case KindDoT:
		return "dot"
	default:
		return fmt.Sprintf("upstream_%d", k)
	}
}

func (k *Kind) fromString(v string) error {
	switch v {
	case "none", "":
		*k = KindNone
	case "plain":
		*k = KindPlain
	case "dot":
		*k = KindDoT
	default:
		return fmt.Errorf("unknown upstream: %s", v)
	}

	return nil
}

func (k Kind) MarshalYAML() (interface{}, error) {
	return k.String(), nil
}

func (k *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	return k.fromString(s)
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *Kind) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}

	return k.fromString(s)
}

func (k *Kind) UnmarshalText(v []byte) error {
	return k.fromString(string(v))
}
