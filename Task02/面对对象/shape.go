package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func RunShapeDemo() {
	r := Rectangle{Width: 3, Height: 4}
	c := Circle{Radius: 5}
	fmt.Println("矩形面积:", r.Area())
	fmt.Println("矩形周长:", r.Perimeter())
	fmt.Println("圆形面积:", c.Area())
	fmt.Println("圆形周长:", c.Perimeter())
}
