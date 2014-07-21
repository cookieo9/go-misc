package avg

// The Averager interface describe a type which maintains
// a running average.
type Averager interface {
	Update(float64)   // Update average w/ new value
	Average() float64 // Get current average
}

// A MovingAverage computes an average of the last
// Size samples.
type MovingAverage struct {
	Size    int
	samples []float64
	sum     float64
}

// Update adds the given sample to the average
// dropping the oldest one.
func (ma *MovingAverage) Update(value float64) {
	ma.samples = append(ma.samples, value)
	ma.sum += value
	for len(ma.samples) > ma.Size {
		ma.sum -= ma.samples[0]
		ma.samples = ma.samples[1:]
	}
}

// Average Computes the current average.
func (ma *MovingAverage) Average() float64 {
	return ma.sum / float64(len(ma.samples))
}

// An AlphaAverage computes a running average
// by using an alpha value to weight the existing
// average with a new sample.
type AlphaAverage struct {
	Alpha, average float64
}

// Update updates the internal average of this
// AlphaAverage using the weighted average
// formula:
//
//	avg = sample * alpha + (1-alpha) * avg
func (aa *AlphaAverage) Update(value float64) {
	aa.average *= 1 - aa.Alpha
	aa.average += aa.Alpha * value
}

// Average returns the current value of the
// running average.
func (aa *AlphaAverage) Average() float64 {
	return aa.average
}
