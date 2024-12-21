// Camera helper functions
package glf

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
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
    Position        mgl64.Vec3
    Front           mgl64.Vec3
    Up              mgl64.Vec3
    Right           mgl64.Vec3

    WorldUp         mgl64.Vec3

    Yaw             float64
    Pitch           float64
    MovementAccel   float64
    MouseSens       float64
    Zoom            float64
}

func NewCamera(camPosition, worldUp mgl64.Vec3, camYaw, camPitch, camAccel, mouseSens float64) *Camera {
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

func (camera *Camera) UpdateCamera(direction Direction, deltaT, xoffset, yoffset float64) {
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

func (camera *Camera) GetViewMatrix() mgl64.Mat4 {
    center := camera.Position.Add(camera.Front)
    return mgl64.LookAt(camera.Position.X(), camera.Position.Y(), camera.Position.Z(),  center.X(), center.Y(), center.Z(), camera.Up.X(), camera.Up.Y(), camera.Up.Z())
}

func (camera *Camera) updateVectors() {
    front := mgl64.Vec3{math.Cos(mgl64.DegToRad(camera.Yaw) * math.Cos(mgl64.DegToRad(camera.Pitch))), 
    math.Sin(mgl64.DegToRad(camera.Pitch)), 
    math.Sin(mgl64.DegToRad(camera.Yaw)) * math.Cos(mgl64.DegToRad(camera.Pitch)), 
    }
    camera.Front = front.Normalize()
    camera.Right = camera.Front.Cross(camera.WorldUp).Normalize()
    camera.Up = camera.Right.Cross(camera.Front).Normalize()
}

