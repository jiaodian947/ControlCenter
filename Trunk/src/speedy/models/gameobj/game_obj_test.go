package gameobj

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/binary"
	"io"
	"speedy/protocol"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNewGameObjectFromBinary(t *testing.T) {
	db, err := sql.Open("mysql", "sa:abc@tcp(192.168.1.180:3306)/nx_merge_test?charset=utf8")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r, err := db.Query("select `lb_save_data` from player_binary limit 0, 1")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if !r.Next() {
		t.Fatal("not found ")
	}
	var data []byte
	r.Scan(&data)
	type args struct {
		data     []byte
		compress bool
	}
	tests := []struct {
		name string
		args args
		want *GameObject
	}{
		{
			name: "load binary",
			args: args{data: data},
			want: NewGameObject(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewGameObjectFromBinary(tt.args.data)
			data := make([]byte, 0, MAX_DATA_LEN)
			ar := protocol.NewStoreArchiver(data)
			err := obj.Store(ar)
			if err != nil {
				t.Fatal(err)
			}
			b := bytes.NewBuffer(tt.args.data[8:])
			r, _ := zlib.NewReader(b)
			var out bytes.Buffer
			io.Copy(&out, r)

			if obj.Compress {
				//原始数据测试
				if out.Len() != ar.Len() {
					t.Fatal("origin len not match", out.Len(), ar.Len())
				}

				d1 := out.Bytes()
				d2 := ar.Data()
				for i := 0; i < out.Len(); i++ {
					if d1[i] != d2[i] {
						t.Fatal("data not match")
					}
				}

				//数据压缩测试
				b := bytes.NewBuffer(nil)
				binary.Write(b, binary.LittleEndian, uint32(COMPRESSED_DATA_VERSION))
				binary.Write(b, binary.LittleEndian, uint32(ar.Len()))
				w := zlib.NewWriter(b)
				w.Write(ar.Data())
				w.Close()

				tmpobj := NewGameObjectFromBinary(b.Bytes())
				data1 := make([]byte, 0, MAX_DATA_LEN)
				ar1 := protocol.NewStoreArchiver(data1)
				tmpobj.Store(ar1)

				if ar1.Len() != ar.Len() {
					t.Fatal("len not match", ar1.Len(), ar.Len())
				}

				d1 = ar1.Data()
				d2 = ar.Data()

				for i := 0; i < ar1.Len(); i++ {
					if d1[i] != d2[i] {
						t.Fatal("data not match")
					}
				}

			} else {
				if ar.Len() != len(tt.args.data) {
					t.Fatal("len not match", ar.Len(), len(tt.args.data))
				}

				d := ar.Data()
				for i := 0; i < ar.Len(); i++ {
					if d[i] != tt.args.data[i] {
						t.Fatal("data not match")
					}
				}

			}

		})
	}
}
