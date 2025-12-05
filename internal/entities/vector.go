package entities

type Vector2 struct {
	X, Y int
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{v.X + other.X, v.Y + other.Y}
}
func (v Vector2) Equal(other Vector2) bool {
	return v.X == other.X && v.Y == other.Y
}
