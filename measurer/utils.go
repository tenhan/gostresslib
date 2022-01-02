package measurer

import "fmt"

func NanosecondsToReadable(nanoseconds int64) string {
	if nanoseconds < 1000{
		return fmt.Sprintf("%dns",nanoseconds)
	}
	us := nanoseconds /1000
	if us < 1000 {
		return fmt.Sprintf("%dus",us)
	}
	ms := us /1000
	if ms < 1000 {
		return fmt.Sprintf("%dms",ms)
	}
	s := ms/1000
	sMs := ms%1000
	if s < 1000{
		return fmt.Sprintf("%ds%dms",s,sMs)
	}
	minutes := s/60
	minutesS := s%60
	if minutes < 60{
		return fmt.Sprintf("%dm%ds",minutes,minutesS)
	}
	hour := minutes/60
	hourM := minutes%60
	if hour < 24{
		return fmt.Sprintf("%dh%dm",hour,hourM)
	}
	day := hour/24
	dayH := hour%24
	return fmt.Sprintf("%dd%dh",day,dayH)
}
