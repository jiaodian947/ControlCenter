package game

import (
	"fmt"
	"math"
	"robots/utils"

	quicklz "github.com/dgryski/go-quicklz"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadPosInfo(ar *utils.LoadArchive) *PosInfo {
	p := &PosInfo{}
	var x, y, z int16
	CheckErr(ar.Read(&x))
	CheckErr(ar.Read(&y))
	CheckErr(ar.Read(&z))
	p.Vec.X = float32(x) / 100
	p.Vec.Y = float32(y) / 100
	p.Vec.Z = float32(z) / 100
	CheckErr(ar.Read(&p.Orient))
	return p
}

func LoadDestInfo(ar *utils.LoadArchive) *DestInfo {
	p := &DestInfo{}
	var x, y, z, movespeed int16

	CheckErr(ar.Read(&x))
	CheckErr(ar.Read(&y))
	CheckErr(ar.Read(&z))
	CheckErr(ar.Read(&p.Orient))
	CheckErr(ar.Read(&movespeed))
	//CheckErr(ar.Read(&p.RotateSpeed))
	//CheckErr(ar.Read(&p.JumpSpeed))
	CheckErr(ar.Read(&p.Mode))
	p.Vec.X = float32(x) / 100
	p.Vec.Y = float32(y) / 100
	p.Vec.Z = float32(z) / 100
	p.MoveSpeed = float32(movespeed) / 100
	return p
}

func OutputObject(obj *GameObject) {
	fmt.Print("output begin\nobject:", obj.ObjId, "config:", obj.ConfigId, "\n")
	for _, v := range obj.Attr.Attr {
		fmt.Println(v.Key, ":", v.Val)
	}
	fmt.Println("end")
}

func Decompress(ar *utils.LoadArchive) *utils.LoadArchive {
	data, err := quicklz.Decompress(ar.Source()[ar.Position():])
	if err != nil {
		panic(err)
	}
	ar = utils.NewLoadArchiver(data)
	return ar
}

type Vector3D struct {
	X float32
	Y float32
	Z float32
}

// 乘法
func (v Vector3D) Mul(val float32) Vector3D {
	v.X *= val
	v.Y *= val
	v.Z *= val
	return v
}

// 除法
func (v Vector3D) Div(val float32) Vector3D {
	v.X /= val
	v.Y /= val
	v.Z /= val
	return v
}

// 相加
func (v Vector3D) Add(other Vector3D) Vector3D {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// 相减
func (v Vector3D) Sub(other Vector3D) Vector3D {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// 计算两个向量的点乘积
func (lhs Vector3D) Dot(rhs Vector3D) float32 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z
}

// 计算两个向量的叉乘
func (lhs Vector3D) Cross(rhs Vector3D) Vector3D {
	x := lhs.Y*rhs.Z - rhs.Y*lhs.Z
	y := lhs.Z*rhs.X - rhs.Z*lhs.X
	z := lhs.X*rhs.Y - rhs.X*lhs.Y
	return Vector3D{x, y, z}
}

// 求夹角
func (from Vector3D) Angle(to Vector3D) float32 {
	dot := from.Dot(to)
	return float32(180 * math.Acos(float64(dot/(from.Magnitude()*to.Magnitude()))) / math.Pi)
}

func (v Vector3D) Normalize() Vector3D {
	magSq := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	if magSq > 0 {
		oneOverMag := 1.0 / float32(math.Sqrt(float64(magSq)))
		v.X *= oneOverMag
		v.Y *= oneOverMag
		v.Z *= oneOverMag
	}
	return v
}

//  两点距离
func (v Vector3D) Distance(other Vector3D) float32 {
	return v.Sub(other).Magnitude()
}

// 求模
func (v Vector3D) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func Lerp(from Vector3D, to Vector3D, t float32) Vector3D {
	if t <= 0 {
		return from
	} else if t >= 1 {
		return to
	}
	return to.Mul(t).Add(from.Mul(1 - t))
}
