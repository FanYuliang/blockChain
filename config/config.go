package config

type Configuration struct {
	ServiceIP		 	string
	ServicePort      	int
	DetectionTimeout 	int64
	SuspiciousTimeout 	int64
	FailureTimeout		int64
	PingNum				int
	TransacCap			int
	PingPeriod 			float64
}