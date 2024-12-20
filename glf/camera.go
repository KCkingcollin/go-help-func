package glf

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
    Position        mgl32.Vec3
    Front           mgl32.Vec3
    Up              mgl32.Vec3
    Right           mgl32.Vec3

    WorldUp         mgl32.Vec3

    Yaw             float32
    Pitch           float32
    MovementSpeed   float32
    MouseSens       float32
    Zoom            float32
}

func NewCamera(position, worldUp mgl32.Vec3, yaw, pitch float32) *Camera {
    camera := Camera{}
    camera.Position = position
    camera.WorldUp = worldUp
    camera.Yaw = yaw
    camera.Pitch = pitch
    camera.updateVectors()
    return &camera
}

func (camera *Camera) updateVectors() {
    front := mgl32.Vec3{float32(math.Cos(float64(mgl32.DegToRad(camera.Yaw)) * math.Cos(float64(mgl32.DegToRad(camera.Pitch))))), 
    float32(math.Sin(float64(mgl32.DegToRad(camera.Pitch)))), 
    float32(math.Sin(float64(mgl32.DegToRad(camera.Yaw)) * math.Cos(float64(mgl32.DegToRad(camera.Pitch))))), 
    }
    camera.Front = front.Normalize()
    camera.Right = camera.Front.Cross(camera.WorldUp).Normalize()
    camera.Up = camera.Right.Cross(camera.Front).Normalize()
}

func (camera *Camera) GetViewMatrix() mgl32.Mat4 {
    center := camera.Position.Add(camera.Front)
    return mgl32.LookAt(camera.Position.X(), camera.Position.Y(), camera.Position.Z(),  center.X(), center.Y(), center.Z(), camera.Up.X(), camera.Up.Y(), camera.Up.Z())
}
