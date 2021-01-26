// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"github.com/Safecast/TTDefs"
)

// Filter classes
const filterClassRadiation = "rad"
const filterClassAir = "air"

func filterMatches(data TTDefs.SafecastData, target string, args map[string]string) bool {
	return true
}

type filterEvent struct {
	class    string
	percent  float64
	summary  string
	country  string
	city     string
	distance float64
	lat      float64
	lon      float64
	device   string
}

func fev(sd TTDefs.SafecastData, ipinfo IPInfoData, class string, summary string, unit string, fmin float64, fmax float64, f float64) (e filterEvent) {

	// Basic info
	e.class = class

	// Location-related
	if sd.Loc != nil && sd.Loc.Lat != nil && sd.Loc.Lon != nil {
		if sd.Loc.LocName != nil {
			e.city = *sd.Loc.LocName
		}
		if sd.Loc.LocCountry != nil {
			e.country = *sd.Loc.LocCountry
		}
		e.lat = *sd.Loc.Lat
		e.lon = *sd.Loc.Lon
		e.distance = distance(ipinfo.Latitude, ipinfo.Longitude, *sd.Loc.Lat, *sd.Loc.Lon)
	}

	// Generate a summary
	e.device = sd.DeviceUID
	if sd.DeviceSN == "" {
		e.summary = "(" + sd.DeviceUID + ")"
	} else {
		e.summary = "(" + sd.DeviceSN + ")"
	}
	e.summary += fmt.Sprintf(" %s %.1f%s", summary, f, unit)
	if f > fmax {
		e.percent = 1.0
	} else if f < fmin {
		e.percent = 0.0
	} else {
		e.percent = (f - fmin) / (fmax - fmin)
	}
	return
}

func filterClassify(sd TTDefs.SafecastData, ipinfo IPInfoData) (events []filterEvent) {

	// Classify radiation
	if sd.Lnd != nil {
		if sd.Lnd.U7318 != nil {
			events = append(events, fev(sd, ipinfo, filterClassRadiation, "U7318", "cpm", 20.0, 50.0, *sd.Lnd.U7318))
		}
		if sd.Lnd.C7318 != nil {
			events = append(events, fev(sd, ipinfo, filterClassRadiation, "C7318", "cpm", 20.0, 50.0, *sd.Lnd.C7318))
		}
		if sd.Lnd.EC7128 != nil {
			events = append(events, fev(sd, ipinfo, filterClassRadiation, "EC7128", "cpm", 5.0, 20.0, *sd.Lnd.EC7128))
		}
		if sd.Lnd.U712 != nil {
			events = append(events, fev(sd, ipinfo, filterClassRadiation, "U712", "cpm", 20.0, 50.0, *sd.Lnd.U712))
		}
		if sd.Lnd.W78017 != nil {
			events = append(events, fev(sd, ipinfo, filterClassRadiation, "W78017", "cpm", 10.0, 50.0, *sd.Lnd.W78017))
		}
	}

	// Classify Air PMS
	if sd.Pms != nil {
		if sd.Pms.Pm02_5 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM2.5", "ug", 0.0, 10.0, *sd.Pms.Pm02_5))
		}
		if sd.Pms.Pm10_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM10.0", "ug", 0.0, 10.0, *sd.Pms.Pm10_0))
		}
		if sd.Pms.Pm01_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM1.0", "ug", 0.0, 10.0, *sd.Pms.Pm01_0))
		}
	}
	if sd.Pms2 != nil {
		if sd.Pms2.Pm02_5 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM2.5", "ug", 0.0, 10.0, *sd.Pms2.Pm02_5))
		}
		if sd.Pms2.Pm10_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM10.0", "ug", 0.0, 10.0, *sd.Pms2.Pm10_0))
		}
		if sd.Pms2.Pm01_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "PMS PM1.0", "ug", 0.0, 10.0, *sd.Pms2.Pm01_0))
		}
	}

	// Classify Air PMS
	if sd.Opc != nil {
		if sd.Opc.Pm02_5 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "OPC PM2.5", "ug", 0.0, 10.0, *sd.Opc.Pm02_5))
		}
		if sd.Opc.Pm10_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "OPC PM10.0", "ug", 0.0, 10.0, *sd.Opc.Pm10_0))
		}
		if sd.Opc.Pm01_0 != nil {
			events = append(events, fev(sd, ipinfo, filterClassAir, "OPC PM1.0", "ug", 0.0, 10.0, *sd.Opc.Pm01_0))
		}
	}

	return

}

// Earth parameters for fast computation
const earthRadiusMeters = 6378100
const earthRadiusMetersDoubled = earthRadiusMeters * 2

// haversin(θ) function
// http://en.wikipedia.org/wiki/Haversine_formula
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS
func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	var la1, lo1, la1Cos, la2, lo2 float64

	// Compute info about the base lat/lon
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la1Cos = math.Cos(la1)

	// convert to radians
	// must cast radius as float to multiply later
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	// calculate
	h := hsin(la2-la1) + la1Cos*math.Cos(la2)*hsin(lo2-lo1)

	return earthRadiusMetersDoubled * math.Asin(math.Sqrt(h))

}
