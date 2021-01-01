// Copyright 2021 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/safecast/ttdata"
)

// Filter classes
const filterClassUnknown = ""
const filterClassRadiation = "rad"
const filterClassAir = "air"

func filterMatches(data ttdata.SafecastData, target string, args map[string]string) bool {
	return true
}

func filterClassify(sd ttdata.SafecastData) (class string, percent float64, summary string) {

	// Begin generating summary
	summary = fmt.Sprintf("%s", sd.DeviceUID)
	if sd.DeviceSN != "" {
		summary += " (" + sd.DeviceSN + ")"
	}

	// Classify radiation
	if sd.Lnd != nil {

		var cpm float64
		if sd.Lnd.U7318 != nil {
			cpm = *sd.Lnd.U7318
			summary += fmt.Sprintf(" U7318 CPM: %0.02f", cpm)
		} else if sd.Lnd.C7318 != nil {
			cpm = *sd.Lnd.C7318
			summary += fmt.Sprintf(" C7318 CPM: %0.02f", cpm)
		} else if sd.Lnd.EC7128 != nil {
			cpm = *sd.Lnd.EC7128
			summary += fmt.Sprintf(" EC7128 CPM: %0.02f", cpm)
		} else if sd.Lnd.U712 != nil {
			cpm = *sd.Lnd.U712
			summary += fmt.Sprintf(" EC712 CPM: %0.02f", cpm)
		} else if sd.Lnd.W78017 != nil {
			cpm = *sd.Lnd.W78017
			summary += fmt.Sprintf(" W78017 CPM: %0.02f", cpm)
		}

		if cpm != 0 {
			class = filterClassRadiation
			cpmMax := float64(80.0)
			cpmMin := float64(20.0)
			if cpm > cpmMax {
				percent = 1.0
			} else if cpm < cpmMin {
				percent = 0.0
			} else {
				percent = (cpm - cpmMin) / (cpmMax - cpmMin)
			}
			return
		}

	}

	// Classify Air PMS
	if sd.Pms != nil {
		var pm float64
		if sd.Pms.Pm02_5 != nil {
			pm = *sd.Pms.Pm02_5
			summary += fmt.Sprintf(" PMS PM 2.5: %0.02f", pm)
		}
		if sd.Pms.Pm10_0 != nil {
			pm = *sd.Pms.Pm10_0
			summary += fmt.Sprintf(" PMS PM 10.0: %0.02f", pm)
		}
		if sd.Pms.Pm01_0 != nil {
			pm = *sd.Pms.Pm01_0
			summary += fmt.Sprintf(" PMS PM 1.0: %0.02f", pm)
		}
		if pm != 0 {
			class = filterClassAir
			pmMax := float64(90.0)
			pmMin := float64(0.0)
			if pm > pmMax {
				percent = 1.0
			} else if pm < pmMin {
				percent = 0.0
			} else {
				percent = (pm - pmMin) / (pmMax - pmMin)
			}
			return
		}

	}

	// Classify Air PMS2
	if sd.Pms2 != nil {
		var pm float64
		if sd.Pms2.Pm02_5 != nil {
			pm = *sd.Pms2.Pm02_5
			summary += fmt.Sprintf(" PMS2 PM 2.5: %0.02f", pm)
		}
		if sd.Pms2.Pm10_0 != nil {
			pm = *sd.Pms2.Pm10_0
			summary += fmt.Sprintf(" PMS2 PM 10.0: %0.02f", pm)
		}
		if sd.Pms2.Pm01_0 != nil {
			pm = *sd.Pms2.Pm01_0
			summary += fmt.Sprintf(" PMS2 PM 1.0: %0.02f", pm)
		}
		if pm != 0 {
			class = filterClassAir
			pmMax := float64(90.0)
			pmMin := float64(0.0)
			if pm > pmMax {
				percent = 1.0
			} else if pm < pmMin {
				percent = 0.0
			} else {
				percent = (pm - pmMin) / (pmMax - pmMin)
			}
			return
		}

	}

	// Unknown
	class = filterClassUnknown
	percent = 0.0
	summary = ""
	return

}
