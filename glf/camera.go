package glf

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Direction int 
const (
    Forward Direction = iota
    Backward
    Left
    Right
    NoWhere
)

type Camera struct {
    Position        mgl32.Vec3
    Front           mgl32.Vec3
    Up              mgl32.Vec3
    Right           mgl32.Vec3

    WorldUp         mgl32.Vec3

    Yaw             float32
    Pitch           float32
    MovementAccel   float32
    MouseSens       float32
    Zoom            float32
}

func NewCamera(camPosition, worldUp mgl32.Vec3, camYaw, camPitch, camAccel, mouseSens float32) *Camera {
    camera := Camera{}
    camera.Position = camPosition
    camera.WorldUp = worldUp
    camera.Yaw = camYaw
    camera.MovementAccel = camAccel
    camera.Pitch = camPitch
    camera.MouseSens = mouseSens
    camera.updateVectors()
    return &camera
}

func (camera *Camera) UpdateCamera(direction Direction, deltaT, xoffset, yoffset float32) {
    magnitude := camera.MovementAccel * deltaT
    switch direction {
    case Forward:
        camera.Position = camera.Position.Add(camera.Front.Mul(magnitude))
    case Backward:
        camera.Position = camera.Position.Sub(camera.Front.Mul(magnitude))
    case Left:
        camera.Position = camera.Position.Sub(camera.Right.Mul(magnitude))
    case Right:
        camera.Position = camera.Position.Add(camera.Right.Mul(magnitude))
    case NoWhere:
    }

    xoffset *= camera.MouseSens
    yoffset *= camera.MouseSens

    camera.Yaw += xoffset
    camera.Pitch += yoffset

    camera.updateVectors()
}

func (camera *Camera) GetViewMatrix() mgl32.Mat4 {
    center := camera.Position.Add(camera.Front)
    return mgl32.LookAt(camera.Position.X(), camera.Position.Y(), camera.Position.Z(),  center.X(), center.Y(), center.Z(), camera.Up.X(), camera.Up.Y(), camera.Up.Z())
}

func (camera *Camera) updateVectors() {
    front := mgl32.Vec3{float32(math.Cos(float64(mgl32.DegToRad(camera.Yaw)) * math.Cos(float64(mgl32.DegToRad(camera.Pitch))))), 
    float32(math.Sin(float64(mgl32.DegToRad(float32(camera.Pitch))))), 
    float32(math.Sin(float64(mgl32.DegToRad(camera.Yaw)) * math.Cos(float64(mgl32.DegToRad(camera.Pitch))))), 
    }
    camera.Front = front.Normalize()
    camera.Right = camera.Front.Cross(camera.WorldUp).Normalize()
    camera.Up = camera.Right.Cross(camera.Front).Normalize()
}

