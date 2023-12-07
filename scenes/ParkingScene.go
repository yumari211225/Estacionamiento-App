package scenes

import (
	"image/color"
	"math/rand"
	"parking/models"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

var (
	spots = []*models.ParkingSpot{
		models.NewParkingSpot(100, 220, 130, 250, 1, 1),
		models.NewParkingSpot(170, 220, 200, 250, 1, 2),
		models.NewParkingSpot(240, 220, 270, 250, 1, 3),
		models.NewParkingSpot(310, 220, 340, 250, 1, 4),
		models.NewParkingSpot(380, 220, 410, 250, 1, 5),
		models.NewParkingSpot(450, 220, 480, 250, 1, 6),
		models.NewParkingSpot(520, 220, 550, 250, 1, 7),
	}
	parking    = models.NewParking(spots)
	doorMutex  sync.Mutex
	carManager = models.NewCarManager()
)

type ParkingScene struct {
}

type Row struct {
}

func NewParkingScene() *ParkingScene {
	return &ParkingScene{}
}

func (ps *ParkingScene) Start() {
	isFirstTime := true

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			_ = ctx.Window.SetBorderless(false)
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for i := 0; i < 100; i++ {
					go carCycle(ctx)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})
}

func setUpScene(ctx *scene.Context) {

	enclosedArea := floatgeom.NewRect2(50, 10, 600, 330)
	entities.New(ctx, entities.WithRect(enclosedArea), entities.WithColor(color.RGBA{100, 100, 100, 255}), entities.WithDrawLayers([]int{0}))

	var sprite *render.Sprite
	area := floatgeom.NewRect2(170, 0, 280, 400)
	sprite, _ = render.LoadSprite("assets/h2.jpg")
	entities.New(ctx, entities.WithRect(area), entities.WithRenderable(sprite), entities.WithDrawLayers([]int{1, 2}))

	parkinglot := floatgeom.NewRect2(202, 202, 202, 255)
	entities.New(ctx, entities.WithRect(parkinglot), entities.WithColor(color.RGBA{34, 139, 34, 255}), entities.WithDrawLayers([]int{0}))
	//145, 10, 425, 400

	green := floatgeom.NewRect2(100, 800, 400, 1000)
	entities.New(ctx, entities.WithRect(green), entities.WithColor(color.RGBA{27, 162, 62, 255}), entities.WithDrawLayers([]int{0}))

	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(172, 170, 455, 175)), entities.WithColor(color.RGBA{50, 50, 50, 255}), entities.WithDrawLayers([]int{0}))
	//entities.New(ctx, entities.WithRect(floatgeom.NewRect2(172, 168, 455, 605)), entities.WithColor(color.RGBA{50, 50, 50, 255}), entities.WithDrawLayers([]int{0}))
	//entities.New(ctx, entities.WithRect(floatgeom.NewRect2(172, 70, 145, 400)), entities.WithColor(color.RGBA{50, 50, 50, 255}), entities.WithDrawLayers([]int{0}))
	//entities.New(ctx, entities.WithRect(floatgeom.NewRect2(455, 10, 430, 400)), entities.WithColor(color.RGBA{50, 50, 50, 255}), entities.WithDrawLayers([]int{0}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(spot.GetArea().Min.X(), spot.GetArea().Min.Y(), spot.GetArea().Min.X()+5, spot.GetArea().Max.Y())), entities.WithColor(color.RGBA{202, 202, 202, 255}), entities.WithDrawLayers([]int{1}))
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(spot.GetArea().Max.X(), spot.GetArea().Min.Y(), spot.GetArea().Max.X()-5, spot.GetArea().Max.Y())), entities.WithColor(color.RGBA{202, 202, 202, 255}), entities.WithDrawLayers([]int{1}))
	}
}

func carCycle(ctx *scene.Context) {
	car := models.NewCar(ctx)

	carManager.AddCar(car)

	car.Enqueue(carManager)

	spotAvailable := parking.GetParkingSpotAvailable()

	doorMutex.Lock()

	car.JoinDoor(carManager)

	doorMutex.Unlock()

	car.Park(spotAvailable, carManager)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(40000, 50000)))

	car.LeaveSpot(carManager)

	parking.ReleaseParkingSpot(spotAvailable)

	car.Leave(spotAvailable, carManager)

	doorMutex.Lock()

	car.ExitDoor(carManager)

	doorMutex.Unlock()

	car.GoAway(carManager)

	car.Remove()

	carManager.RemoveCar(car)
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
