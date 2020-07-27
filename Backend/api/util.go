package api

import (
	"math"
	"math/rand"
)

// cribbed from the javascript on https://www.random.org/geographic-coordinates
// definitely puts the flag on a random place on earth, and not inside the coordinates
// not the important part of the demo, but fix later
func getRandomLocation(lat float64, lng float64, radius uint8) (float64, float64) {
	/*
	  floats = xmlhttp.responseText.split("\n");
	  x = floats[0] * 2 * Math.PI - Math.PI;
	  y = floats[1] * 2 - 1;
	  lng = rad2deg(x).toFixed(5);
	  latrad = Math.PI/2 - Math.acos(y);
	  lat = rad2deg(latrad).toFixed(5);
	  distortion = Math.pow(sec(latrad), 2).toFixed(2);
	*/

	float0 := rand.Float64()
	float1 := rand.Float64()

	x := float0*2*math.Pi - math.Pi
	y := float1*2 - 1
	flagLng := rad2deg(x)
	latrad := math.Pi/2 - math.Acos(y)
	flagLat := rad2deg(latrad)
	//distortion := math.Pow(seconds(latrad), 2)

	return flagLat, flagLng
}

func rad2deg(arg float64) float64 {
	return 360 * arg / (2 * math.Pi)
}

func seconds(arg float64) float64 {
	return 1 / math.Cos(arg)
}

// taken from https://www.nhc.noaa.gov/gccalc2.js
func computeDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	/*
	lat1=(Math.PI/180)*signlat1*checkField(document.InputFormCD.lat1)
	lat2=(Math.PI/180)*signlat2*checkField(document.InputFormCD.lat2)
	lon1=(Math.PI/180)*signlon1*checkField(document.InputFormCD.lon1)
	lon2=(Math.PI/180)*signlon2*checkField(document.InputFormCD.lon2)

	dc=1 // get distance conversion factor, 1 == nautical mile

	ellipse=new ellipsoid("Sphere", 180*60/Math.PI,"Infinite") //get ellipse

	// spherical code
	cd=crsdist(lat1,lon1,lat2,lon2) // compute crs and distance
	d=cd.d*(180/Math.PI)*60*dc  // go to physical units

	return Math.round(d)

	*/

	lat1 = degrees(lat1)
	lng1 = degrees(lng1)
	lat2 = degrees(lat2)
	lng2 = degrees(lng2)

	cd := crsdist(lat1, lng1, lat2, lng2)
	d := cd * (180 / math.Pi) * 60
	return math.Round(d / feetPerNauticalMile)
}

const feetPerNauticalMile = 6076.118

func degrees(arg float64) float64 {
	return math.Pi / 180 * arg
}

func crsdist(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	/*
		// radian args
		// compute course and distance (spherical)
		d=acos(sin(lat1)*sin(lat2)+cos(lat1)*cos(lat2)*cos(lon1-lon2))

	*/

	return math.Acos(
		math.Sin(lat1)*math.Sin(lat2) +
			math.Cos(lat1)*math.Cos(lat2)*math.Cos(lng1-lng2))
}
