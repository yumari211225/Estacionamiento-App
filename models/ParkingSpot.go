package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type ParkingSpot struct {
	area                 *floatgeom.Rect2
	directionsForParking *[]ParkingSpotDirection
	directionsForLeaving *[]ParkingSpotDirection
	number               int
	isAvailable          bool
}

func NewParkingSpot(x, y, x2, y2 float64, row, number int) *ParkingSpot {
	directionsForParking := getDirectionForParking(x, y, row)
	directionsForLeaving := getDirectionsForLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &ParkingSpot{
		area:                 &area,
		directionsForParking: directionsForParking,
		directionsForLeaving: directionsForLeaving,
		number:               number,
		isAvailable:          true,
	}
}

func getDirectionForParking(x, y float64, row int) *[]ParkingSpotDirection {
	var directions []ParkingSpotDirection

	if row == 1 {
		directions = append(directions, *newParkingSpotDirection("down", 30))
	} else if row == 2 {
		directions = append(directions, *newParkingSpotDirection("down", 140))
	} else if row == 3 {
		directions = append(directions, *newParkingSpotDirection("down", 230))
	} else if row == 4 {
		directions = append(directions, *newParkingSpotDirection("down", 320))
	}

	directions = append(directions, *newParkingSpotDirection("right", x+10))
	directions = append(directions, *newParkingSpotDirection("down", y+2))

	return &directions
}

func getDirectionsForLeaving() *[]ParkingSpotDirection {
	var directions []ParkingSpotDirection

	directions = append(directions, *newParkingSpotDirection("right", 400))
	directions = append(directions, *newParkingSpotDirection("up", 20))
	directions = append(directions, *newParkingSpotDirection("left", 155))

	return &directions
}

func (p *ParkingSpot) GetArea() *floatgeom.Rect2 {
	return p.area
}

func (p *ParkingSpot) GetNumber() int {
	return p.number
}

func (p *ParkingSpot) GetDirectionsForParking() *[]ParkingSpotDirection {
	return p.directionsForParking
}

func (p *ParkingSpot) GetDirectionsForLeaving() *[]ParkingSpotDirection {
	return p.directionsForLeaving
}

func (p *ParkingSpot) GetIsAvailable() bool {

	return p.isAvailable
}

func (p *ParkingSpot) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
