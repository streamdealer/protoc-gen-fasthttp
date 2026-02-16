package http

import (
	"bytes"

	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ToProto(ctx *fasthttp.RequestCtx, msg proto.Message) error {
	contentType := ctx.Request.Header.ContentType()
	if len(ctx.PostBody()) > 0 {
		switch {
		case bytes.HasPrefix(contentType, []byte("application/json")):
			err := UnmarshalerCtx(ctx).Unmarshal(ctx.PostBody(), msg)
			if err != nil {
				ctx.Error("Invalid JSON body", fasthttp.StatusBadRequest)
				return err
			}
		}
	}

	m := msg.ProtoReflect()
	mDesc := m.Descriptor()

	ctx.QueryArgs().All()(func(k, v []byte) bool {
		if fd := mDesc.Fields().ByName(protoreflect.Name(k)); fd != nil {
			setProtoField(m, fd, string(v))
		}

		return true
	})

	ctx.VisitUserValues(func(k []byte, v any) {
		if fd := mDesc.Fields().ByName(protoreflect.Name(k)); fd != nil {
			setProtoField(m, fd, v.(string))
		}
	})

	return nil
}

func setProtoField(msg protoreflect.Message, fd protoreflect.FieldDescriptor, val string) {
	switch fd.Kind() {
	case protoreflect.BytesKind:
		msg.Set(fd, protoreflect.ValueOfBytes([]byte(val)))
	case protoreflect.BoolKind:
		msg.Set(fd, protoreflect.ValueOfBool(cast.ToBool(val)))
	case protoreflect.EnumKind:
		msg.Set(fd, protoreflect.ValueOfEnum(protoreflect.EnumNumber(cast.ToInt32(val))))
	case protoreflect.Int32Kind:
		msg.Set(fd, protoreflect.ValueOfInt32(cast.ToInt32(val)))
	case protoreflect.Sint32Kind:
		msg.Set(fd, protoreflect.ValueOfUint32(cast.ToUint32(val)))
	case protoreflect.Uint32Kind:
		msg.Set(fd, protoreflect.ValueOfUint32(cast.ToUint32(val)))
	case protoreflect.Int64Kind:
		msg.Set(fd, protoreflect.ValueOfInt64(cast.ToInt64(val)))
	case protoreflect.Sint64Kind:
		msg.Set(fd, protoreflect.ValueOfUint64(cast.ToUint64(val)))
	case protoreflect.Uint64Kind:
		msg.Set(fd, protoreflect.ValueOfUint64(cast.ToUint64(val)))
	case protoreflect.Sfixed32Kind:
		msg.Set(fd, protoreflect.ValueOfFloat32(cast.ToFloat32(val)))
	case protoreflect.Fixed32Kind:
		msg.Set(fd, protoreflect.ValueOfFloat32(cast.ToFloat32(val)))
	case protoreflect.FloatKind:
		msg.Set(fd, protoreflect.ValueOfFloat32(cast.ToFloat32(val)))
	case protoreflect.Sfixed64Kind:
		msg.Set(fd, protoreflect.ValueOfFloat64(cast.ToFloat64(val)))
	case protoreflect.Fixed64Kind:
		msg.Set(fd, protoreflect.ValueOfFloat64(cast.ToFloat64(val)))
	case protoreflect.DoubleKind:
		msg.Set(fd, protoreflect.ValueOfFloat64(cast.ToFloat64(val)))
	case protoreflect.StringKind:
		msg.Set(fd, protoreflect.ValueOfString(val))
	}
}
