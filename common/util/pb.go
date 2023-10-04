package util

import (
	"google.golang.org/protobuf/types/known/structpb"
)

func ToStructPb(v any) (result *structpb.Struct, err error) {
	tmp := make(map[string]any)
	if err = ParseAnyToAny(v, &tmp); err != nil {
		return nil, err
	}
	result, err = structpb.NewStruct(tmp)
	return
}
