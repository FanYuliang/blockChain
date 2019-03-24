package config

type Configuration struct {
	ServiceIP		 	string
	ServicePort      	int
	Debug 			 	bool
	DetectionTimeout 	int64
	SuspiciousTimeout 	int64
	FailureTimeout		int64
	LeaveTimeout		int64
	PingNum				int
	TransacCap			int
	PingPeriod 			float64
}