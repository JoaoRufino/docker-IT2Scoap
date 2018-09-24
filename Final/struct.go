package main

type frontendMsg struct {
	Header struct {
		ProtocolVersion int `json:"Version"`
		MessageID       int `json:"MessageID"`
		StationID       int `json:"StationID"`
	} `json:"header"`

	Cam struct {
		GenerationDeltaTime int `json:"generationdeltatime"`
		CamParameters       struct {
			BasicContainer struct {
				StationType       int `json:"stationtype"`
				ReferencePosition struct {
					Latitude                  int `json:"latitude"`
					Longitude                 int `json:"longitude"`
					PositionConfidenceEllipse struct {
						SemiMajorConfidence  int `json:"semimajorconfidence"`
						SemiMinorConfidence  int `json:"semiminorconfidence"`
						SemiMajorOrientation int `json:"SemiMajorOrientation"`
					} `json:"PositionConfidenceEllipse"`
					Altitude struct {
						AltitudeValue      int `json:"altitudevalue"`
						AltitudeConfidence int `json:"altitudeconfidence"`
					} `json:"altitude"`
				} `json:"referenceposition"`
			} `json:"basiccontainer"`
			HighFrequencyContainer struct {
				BasicVehicleContainerHighFrequency struct {
					Heading struct {
						HeadingValue      int `json:"HeadingValue"`
						HeadingConfidence int `json:HeadingConfidence"`
					} `json:"Heading"`
					Speed struct {
						SpeedValue      int `json:"speedValue"`
						SpeedConfidence int `json:"speedconfidence"`
					} `json:"speed"`
					DriveDirection int `json:"drivedirection"`
					VehicleLength  struct {
						VehicleLengthValue                int `json:"vehiclelengthvalue"`
						VehicleLengthConfidenceIndication int `json:"vehiclelengthconfidenceIndication"`
					} `json:"vehiclelength"`
					VehicleWidth             int `json:"vehiclewidth"`
					LongitudinalAcceleration struct {
						LongitudinalAccelerationValue      int `json:"longitudinalaccelerationvalue"`
						LongitudinalAccelerationConfidence int `json:"longitudinalaccelerationconfidence"`
					} `json:"longitudinalacceleration"`
					Curvature struct {
						CurvatureValue      int `json:"curvaturevalue"`
						CurvatureConfidence int `json:"curvatureconfidence"`
					} `json:"curvature"`
					CurvatureCalculationMode int `json:"curvaturecalculationmode"`
					YawRate                  struct {
						YawRateValue      int `son:"yawratevalue"`
						YawRateConfidence int `json:"yawrateconfidence"`
					} `json:"yawrate"`
					AccelerationControl int `json:"accelerationcontrol,omitempty"`
					//[]bool
					LanePosition       int `json:"laneposition,omitempty"`
					SteeringWheelAngle struct {
						SteeringWheelAngleValue      int `json:"steeringwheelanglevalue"`
						SteeringWheelAngleConfidence int ` json:"steeringwheelangleconfidence"`
					} ` json:"steeringwheelangle,omitempty"`
					LateralAcceleration struct {
						LateralAccelerationValue      int `json:"lateralaccelerationvalue,"`
						LateralAccelerationConfidence int `json:"lateralaccelerationconfidence"`
					} `json:"lateralacceleration,omitempty"`
					VerticalAcceleration struct {
						VerticalAccelerationValue      int `json:"verticalaccelerationvalue,"`
						VerticalAccelerationConfidence int `json:"verticalaccelerationconfidence"`
					} `json:"lateralacceleration,omitempty"`
					PerformanceClass   int `json:"performanceclass"`
					CenDsrcTollingZone struct {
						ProtectedZoneLatitude  int `json:"protectedzonelatitude"`
						ProtectedZoneLongitude int `json:"protectedzonelongitude"`
						CenDsrcTollingZoneID   int `json:"cendsrctollingzoneid,omitempty"`
					}
				} `json:"basicvehiclecontainerhighfrequency,omitempty"`
				RsuContainerHighFrequency struct {
					ProtectedCommunicationZonesRSU struct {
						ProtectedCommunicationZone []struct {
							ProtectedZoneType      int `json:"ProtectedZoneType"`
							ExpiryTime             int `json:TimeStampits,omitempty"`
							ProtectedZoneLatitude  int `json:Latitude"`
							ProtectedZoneLongitude int `json:Longitude"`
							ProtectedZoneRadius    int `json:protectedZoneRadius,omitempty"`
							ProtectedZoneID        int `json:protectedZoneID,omitempty"`
						} `json:"ProtectedCommunicationZone"`
					} `json:"ProtectedCommunicationZonesRSU,omitempty"`
				} `json:"rsucontainerhighfrequency,omitempty"`
			} `json:"highfrequencycontainer"`

			LowFrequencyContainer struct {
				BasicVehicleContainerLowFrequency struct {
					VehicleRole    int `json:"vehiclerole"`
					ExteriorLights int `json:"exteriorlights,omitempty"`
					/*	LowBeamHeadlightsOn    bool `json:"lowBeamHeadlightsOn"`    // (0),
						HighBeamHeadlightsOn   bool `json:"highBeamHeadlightsOn"`   //(1),
						LeftTurnSignalOn       bool `json:"LeftTurnSignalOn"`       //(2),
						RightTurnSignalOn      bool `json:"RightTurnSignalOn"`      //(3),
						DaytimeRunningLightsOn bool `json:"DaytimeRunningLightsOn"` //(4),
						ReverseLightOn         bool `json:"ReverseLightOn"`         //(5),
						FogLightOn             bool `json:"FooLightOn"`             //(6),
						ParkingLightsOn        bool `json:"ParkingLightsOn"`        //(7)
					} */
					PathHistory struct {
						PathPoint []pathPoint `json:"pathpoint"`
					} `json:"pathhistory,omitempty"`
				} `json:"basicvehiclecontainerlowfrequency,omitempty"`
			} `json:"lowfrequencycontainer,omitempty"`
			SpecialVehicleContainer struct {
				PublicTransportContainer struct {
					EmbarkationStatus bool `json:"embarkationstatus"`
					PtActivation      struct {
						ptActivationType int    `json:"ptActivationType"`
						ptActivationData string `json:"ptactivationdata"`
					} `json:"PtActivation,omitempty"`
				} `json:"PublicTransportContainer,omitempty"`
				SpecialTransportContainer struct {
					SpecialTransportType int `json:"specialtransporttype"`
					//heavyLoad(0),
					//excessWidth(1),
					//excessLength(2),
					//excessHeight(3)
					LightBarSirenInUse int `json:"lightbarsirenuse,omitempty"`
				} `json:"specialTransportContainer,omitempty"`

				DangerousGoodsContainer struct {
					DangerousGoodsBasic int `"dangerousgoodbasic"`
				} `json:"dangerousGoodsContainer,omitempty"`

				RoadWorksContainerBasic struct {
					RoadworksSubCauseCode int `"roadworkssubcausecode,omitempty"`
					LightBarSirenInUse    int `json:"lightbarsirenuse,omitempty"`
					ClosedLanes           struct {
						HardShoulderStatus int `json:"hardshoulderstatus,omitempty"`
						DrivingLaneStatus  int `json:"drivinglanestatus"`
						//outermostLaneClosed(1),
						//secondLaneFromOutsideClosed(2)
					} `json:"closedLanes"`
				} `json:"roadworkscontainerbasic,omitempty"`

				RescueContainer struct {
					LightBarSirenInUse int `json:"lightbarsirenuse,omitempty"`
				} `json:"rescuecontainer,omitempty"`

				EmergencyContainer struct {
					LightBarSirenInUse int `json:"lightbarsirenuse"`
					IncidentIndication int `json:"incidentIndication,omitempty"`
					EmergencyPriority  int `json:"emergencyPriority,omitempty"`
					//	requestForRightOfWay(0),
					//requestForFreeCrossingAtATrafficLight(1)
				} `json:"emergencycontainer,omitempty"`
				SafetyCarContainer struct {
					LightBarSirenInUse int `json:lightbarsirenuse"`
					IncidentIndication int `json:"incidentIndication",omitempty"`
					TrafficRule        int `json:"TrafficRule,omitempty"`
					SpeedLimit         int `json:"SpeedLimit"`
				} `json:"safetyCarContainer,omitempty"`
			} `json:"SpecialVehicleContainer,omitempty"`
		} `json:"camparameters"`
	} `json:"cam"`
}

type pathPoint struct {
	PathPosition struct {
		DeltaLatitude  int `json:"deltalatitude,omitempty"`
		DeltaLongitude int `json:"deltalongitude"`
		DeltaAltitude  int `json:"deltaaltitude"`
	} `json:"pathposition"`
	PathDeltaTime int `json:pathDeltaTime,omitempty"`
}
