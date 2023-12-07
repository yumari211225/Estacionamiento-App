package models

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
	entranceSpotX = 140.00
	speed         = 10
)

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	var sprite *render.Sprite

	area := floatgeom.NewRect2(0, 50, 90, 60)

	carModelNum := int(getRandomNumber(1, 5))

	if carModelNum == 1 {
		sprite, _ = render.LoadSprite("assets/car2.png")
	} else if carModelNum == 2 {
		sprite, _ = render.LoadSprite("assets/car3.png")
	} else if carModelNum == 3 {
		sprite, _ = render.LoadSprite("assets/car5.png")
	} else if carModelNum == 4 {
		sprite, _ = render.LoadSprite("assets/car1.png")
	} else if carModelNum == 5 {
		sprite, _ = render.LoadSprite("assets/car4.png")
	}

	swtch := render.NewSwitch("left", map[string]render.Modifiable{
		"up":    sprite,
		"down":  sprite.Copy().Modify(mod.FlipY),
		"left":  sprite.Copy().Modify(mod.Rotate(90)),
		"right": sprite.Copy().Modify(mod.Rotate(-90)),
	})

	entity := entities.New(ctx, entities.WithRect(area), entities.WithRenderable(swtch), entities.WithDrawLayers([]int{1, 2}))

	return &Car{
		area:   area,
		entity: entity,
	}
}

func (c *Car) Enqueue(manager *CarManager) {
	for c.X() < entranceSpotX {
		if !c.isCollision("right", manager.GetCars()) {
			c.entity.Renderable.(*render.Switch).Set("right")
			c.ShiftX(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) JoinDoor(manager *CarManager) {
	for c.X() < 20 {
		if !c.isCollision("right", manager.GetCars()) {
			c.entity.Renderable.(*render.Switch).Set("right")
			c.ShiftX(1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ExitDoor(manager *CarManager) {
	for c.X() > 300 {
		if !c.isCollision("left", manager.GetCars()) {
			c.entity.Renderable.(*render.Switch).Set("left")
			c.ShiftX(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) Park(spot *ParkingSpot, manager *CarManager) {
	for index := 0; index < len(*spot.GetDirectionsForParking()); index++ {
		directions := *spot.GetDirectionsForParking()
		if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.isCollision("right", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("right")
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.isCollision("down", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("down")
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}
}

func (c *Car) Leave(spot *ParkingSpot, manager *CarManager) {
	for index := 0; index < len(*spot.GetDirectionsForLeaving()); index++ {
		directions := *spot.GetDirectionsForLeaving()
		if directions[index].Direction == "left" {
			for c.X() > directions[index].Point {
				if !c.isCollision("left", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("left")
					c.ShiftX(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "right" {
			for c.X() < directions[index].Point {
				if !c.isCollision("right", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("right")
					c.ShiftX(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "up" {
			for c.Y() > directions[index].Point {
				if !c.isCollision("up", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("up")
					c.ShiftY(-1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		} else if directions[index].Direction == "down" {
			for c.Y() < directions[index].Point {
				if !c.isCollision("down", manager.GetCars()) {
					c.entity.Renderable.(*render.Switch).Set("down")
					c.ShiftY(1)
					time.Sleep(speed * time.Millisecond)
				}
			}
		}
	}
}

func (c *Car) LeaveSpot(manager *CarManager) {
	spotY := c.Y()
	for c.Y() < spotY+30 {
		if !c.isCollision("down", manager.GetCars()) {
			c.entity.Renderable.(*render.Switch).Set("down")
			c.ShiftY(1)
			time.Sleep(speed * time.Millisecond)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (c *Car) GoAway(manager *CarManager) {
	for c.X() > -20 {
		if !c.isCollision("left", manager.GetCars()) {
			c.entity.Renderable.(*render.Switch).Set("left")
			c.ShiftX(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}

func (c *Car) ShiftY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) ShiftX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) isCollision(direction string, cars []*Car) bool {
	minDistance := 35.0
	for _, car := range cars {
		if direction == "left" {
			if c.X() > car.X() && c.X()-car.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "right" {
			if c.X() < car.X() && car.X()-c.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if direction == "up" {
			if c.Y() > car.Y() && c.Y()-car.Y() < minDistance && c.X() == car.X() {
				return true
			}
		} else if direction == "down" {
			if c.Y() < car.Y() && car.Y()-c.Y() < minDistance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
