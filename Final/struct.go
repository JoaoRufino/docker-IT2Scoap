// Created by cgo -godefs - DO NOT EDIT
// cgo --godefs generateStruct.go

package main

type CauseCode struct {
	CauseCode    int64
	SubCauseCode int64
	X_asn_ctx    asn_struct_ctx
}
type EventHistory struct {
	List      [23]EventPoint
	X_asn_ctx asn_struct_ctx
}
type Heading struct {
	HeadingValue      int64
	HeadingConfidence int64
	X_asn_ctx         asn_struct_ctx
}
type ItsPduHeader struct {
	ProtocolVersion int64
	MessageID       int64
	StationID       uint64
	X_asn_ctx       asn_struct_ctx
}
type ReferencePosition struct {
	Latitude                  int64
	Longitude                 int64
	PositionConfidenceEllipse PosConfidenceEllipse
	Altitude                  Altitude
	X_asn_ctx                 asn_struct_ctx
}
type Speed struct {
	SpeedValue      int64
	SpeedConfidence int64
	X_asn_ctx       asn_struct_ctx
}
type Traces struct {
	List      [7]PathHistory
	X_asn_ctx asn_struct_ctx
}

type Altitude struct {
	AltitudeValue      int64
	AltitudeConfidence int64
	X_asn_ctx          asn_struct_ctx
}
type asn_struct_ctx struct {
	Phase   int16
	Step    int16
	Context int32
	Ptr     *byte
	Left    int64
}
type BasicVehicleContainerHighFrequency struct {
	Heading                  Heading
	Speed                    Speed
	DriveDirection           int64
	VehicleLength            VehicleLength
	VehicleWidth             int64
	LongitudinalAcceleration LongitudinalAcceleration
	Curvature                Curvature
	CurvatureCalculationMode int64
	YawRate                  YawRate
	AccelerationControl      *BIT_STRING
	LanePosition             *int64
	SteeringWheelAngle       *SteeringWheelAngle
	LateralAcceleration      *LateralAcceleration
	VerticalAcceleration     *VerticalAcceleration
	PerformanceClass         *int64
	CenDsrcTollingZone       *CenDsrcTollingZone
	X_asn_ctx                asn_struct_ctx
}
type BasicVehicleContainerLowFrequency struct {
	VehicleRole    int64
	ExteriorLights byte
	PathHistory    PathHistory
	X_asn_ctx      asn_struct_ctx
}
type BasicContainer struct {
	StationType       int64
	ReferencePosition ReferencePosition
	X_asn_ctx         asn_struct_ctx
}
type BIT_STRING struct {
	Buf       *uint8
	Size      int32
	Unused    int32
	X_asn_ctx asn_struct_ctx
}
type CoopAwareness struct {
	GenerationDeltaTime int64
	CamParameters       CamParameters
	X_asn_ctx           asn_struct_ctx
}
type CamParameters struct {
	BasicContainer          BasicContainer
	HighFrequencyContainer  HighFrequencyContainer
	LowFrequencyContainer   *LowFrequencyContainer
	SpecialVehicleContainer *SpecialVehicleContainer
	X_asn_ctx               asn_struct_ctx
}
type CenDsrcTollingZone struct {
	ProtectedZoneLatitude  int64
	ProtectedZoneLongitude int64
	CenDsrcTollingZoneID   *int64
	X_asn_ctx              asn_struct_ctx
}
type ClosedLanes struct {
	HardShoulderStatus *int64
	DrivingLaneStatus  BIT_STRING
	X_asn_ctx          asn_struct_ctx
}
type Curvature struct {
	CurvatureValue      int64
	CurvatureConfidence int64
	X_asn_ctx           asn_struct_ctx
}
type DangerousGoodsContainer struct {
	DangerousGoodsBasic int64
	X_asn_ctx           asn_struct_ctx
}
type DangerousGoodsExtended struct {
	DangerousGoodsType  int64
	UnNumber            int64
	ElevatedTemperature int32
	TunnelsRestricted   int32
	LimitedQuantity     int32
	Pad_cgo_0           [4]byte
	EmergencyActionCode *OCTET_STRING
	PhoneNumber         *OCTET_STRING
	CompanyName         *OCTET_STRING
	X_asn_ctx           asn_struct_ctx
}
type DeltaReferencePosition struct {
	DeltaLatitude  int64
	DeltaLongitude int64
	DeltaAltitude  int64
	X_asn_ctx      asn_struct_ctx
}
type EmergencyContainer struct {
	LightBarSirenInUse BIT_STRING
	IncidentIndication *CauseCode
	EmergencyPriority  *BIT_STRING
	X_asn_ctx          asn_struct_ctx
}
type EventPoint struct {
	EventPosition      DeltaReferencePosition
	EventDeltaTime     *int64
	InformationQuality int64
	X_asn_ctx          asn_struct_ctx
}
type HighFrequencyContainer struct {
	Present                            uint32
	BasicVehicleContainerHighFrequency BasicVehicleContainerHighFrequency
	RSUContainerHighFrequency          RSUContainerHighFrequency
	X_asn_ctx                          asn_struct_ctx
}
type ItineraryPath struct {
	List      [40]ReferencePosition
	X_asn_ctx asn_struct_ctx
}
type LateralAcceleration struct {
	LateralAccelerationValue      int64
	LateralAccelerationConfidence int64
	X_asn_ctx                     asn_struct_ctx
}
type LongitudinalAcceleration struct {
	LongitudinalAccelerationValue      int64
	LongitudinalAccelerationConfidence int64
	X_asn_ctx                          asn_struct_ctx
}
type LowFrequencyContainer struct {
	Present                           uint32
	BasicVehicleContainerLowFrequency BasicVehicleContainerLowFrequency
	X_asn_ctx                         asn_struct_ctx
}
type OCTET_STRING struct {
	Buf       *uint8
	Size      int32
	Pad_cgo_0 [4]byte
	X_asn_ctx asn_struct_ctx
}
type PathHistory struct {
	List      [40]PathPoint
	X_asn_ctx asn_struct_ctx
}
type PathPoint struct {
	PathPosition  DeltaReferencePosition
	PathDeltaTime *int64
	X_asn_ctx     asn_struct_ctx
}
type PosConfidenceEllipse struct {
	SemiMajorConfidence  int64
	SemiMinorConfidence  int64
	SemiMajorOrientation int64
	X_asn_ctx            asn_struct_ctx
}
type PositionOfPillars_t struct {
	List      [3]int64
	X_asn_ctx asn_struct_ctx
}
type ProtectedCommunicationZone struct {
	ProtectedZoneType      int64
	ExpiryTime             *TimestampIts
	ProtectedZoneLatitude  int64
	ProtectedZoneLongitude int64
	ProtectedZoneRadius    *int64
	ProtectedZoneID        *int64
	X_asn_ctx              asn_struct_ctx
}
type ProtectedCommunicationZonesRSU struct {
	List      [16]ProtectedCommunicationZone
	X_asn_ctx asn_struct_ctx
}
type PtActivation struct {
	PtActivationType int64
	PtActivationData OCTET_STRING
	X_asn_ctx        asn_struct_ctx
}
type PublicTransportContainer struct {
	EmbarkationStatus int32
	Pad_cgo_0         [4]byte
	PtActivation      *PtActivation
	X_asn_ctx         asn_struct_ctx
}
type RescueContainer struct {
	LightBarSirenInUse BIT_STRING
	X_asn_ctx          asn_struct_ctx
}
type RestrictedTypes struct {
	List      [3]int64
	X_asn_ctx asn_struct_ctx
}
type RoadWorksContainerBasic struct {
	RoadworksSubCauseCode *int64
	LightBarSirenInUse    BIT_STRING
	ClosedLanes           *ClosedLanes
	X_asn_ctx             asn_struct_ctx
}
type RSUContainerHighFrequency struct {
	ProtectedCommunicationZonesRSU *ProtectedCommunicationZonesRSU
	X_asn_ctx                      asn_struct_ctx
}
type SafetyCarContainer struct {
	LightBarSirenInUse BIT_STRING
	IncidentIndication *CauseCode
	TrafficRule        *int64
	SpeedLimit         *int64
	X_asn_ctx          asn_struct_ctx
}
type SpecialTransportContainer struct {
	SpecialTransportType BIT_STRING
	LightBarSirenInUse   BIT_STRING
	X_asn_ctx            asn_struct_ctx
}
type SpecialVehicleContainer struct {
	Present                   uint32
	PublicTransportContainer  PublicTransportContainer
	SpecialTransportContainer SpecialTransportContainer
	DangerousGoodsContainer   DangerousGoodsContainer
	RoadWorksContainerBasic   RoadWorksContainerBasic
	RescueContainer           RescueContainer
	EmergencyContainer        EmergencyContainer
	SafetyCarContainer        SafetyCarContainer
}
type SteeringWheelAngle struct {
	SteeringWheelAngleValue      int64
	SteeringWheelAngleConfidence int64
	X_asn_ctx                    asn_struct_ctx
}
type VehicleIdentification struct {
	WMInumber *OCTET_STRING
	VDS       *OCTET_STRING
	X_asn_ctx asn_struct_ctx
}
type VehicleLength struct {
	VehicleLengthValue                int64
	VehicleLengthConfidenceIndication int64
	X_asn_ctx                         asn_struct_ctx
}
type VerticalAcceleration struct {
	VerticalAccelerationValue      int64
	VerticalAccelerationConfidence int64
	X_asn_ctx                      asn_struct_ctx
}
type YawRate struct {
	YawRateValue      int64
	YawRateConfidence int64
	X_asn_ctx         asn_struct_ctx
}

type CAM struct {
	Header    ItsPduHeader
	Cam       CoopAwareness
	X_asn_ctx asn_struct_ctx
}

type ReferenceDenms struct {
	List      [8]ActionID
	X_asn_ctx asn_struct_ctx
}
type ActionID struct {
	OriginatingStationID uint64
	SequenceNumber       int64
	X_asn_ctx            asn_struct_ctx
}
type StationaryVehicleContainer struct {
	StationarySince        *int64
	StationaryCause        *CauseCode
	CarryingDangerousGoods *DangerousGoodsExtended
	NumberOfOccupants      *int64
	VehicleIdentification  *VehicleIdentification
	EnergyStorageType      *BIT_STRING
	X_asn_ctx              asn_struct_ctx
}
type RoadWorksContainerExtended struct {
	LightBarSirenInUse      *BIT_STRING
	ClosedLanes             *ClosedLanes
	Restriction             *RestrictedTypes
	SpeedLimit              *int64
	IncidentIndication      *CauseCode
	RecommendedPath         *ItineraryPath
	StartingPointSpeedLimit *DeltaReferencePosition
	TrafficFlowRule         *int64
	ReferenceDenms          *ReferenceDenms
	X_asn_ctx               asn_struct_ctx
}
type ImpactReductionContainer struct {
	HeightLonCarrLeft         int64
	HeightLonCarrRight        int64
	PosLonCarrLeft            int64
	PosLonCarrRight           int64
	PositionOfPillars         PositionOfPillars_t
	PosCentMass               int64
	WheelBaseVehicle          int64
	TurningRadius             int64
	PosFrontAx                int64
	PositionOfOccupants       BIT_STRING
	VehicleMass               int64
	RequestResponseIndication int64
	X_asn_ctx                 asn_struct_ctx
}
type AlacarteContainer struct {
	LanePosition        *int64
	ImpactReduction     *ImpactReductionContainer
	ExternalTemperature *int64
	RoadWorks           *RoadWorksContainerExtended
	PositioningSolution *int64
	StationaryVehicle   *StationaryVehicleContainer
	X_asn_ctx           asn_struct_ctx
}
type LocationContainer struct {
	EventSpeed           *Speed
	EventPositionHeading *Heading
	Traces               Traces
	RoadType             *int64
	X_asn_ctx            asn_struct_ctx
}
type SituationContainer struct {
	InformationQuality int64
	EventType          CauseCode
	LinkedCause        *CauseCode
	EventHistory       *EventHistory
	X_asn_ctx          asn_struct_ctx
}
type ManagementContainer struct {
	ActionID                  ActionID
	DetectionTime             TimestampIts
	ReferenceTime             TimestampIts
	Termination               int64
	EventPosition             ReferencePosition
	RelevanceDistance         *int64
	RelevanceTrafficDirection *int64
	ValidityDuration          *int64
	TransmissionInterval      *int64
	StationType               *int64
	Tmp                       [2]int64
	X_asn_ctx                 asn_struct_ctx
}

type DecentralizedEnvironmentalNotificationMessage struct {
	Management ManagementContainer
	Situation  SituationContainer
	Location   *LocationContainer
	Alacarte   *AlacarteContainer
	X_asn_ctx  asn_struct_ctx
}

type DENM struct {
	Header    ItsPduHeader
	Denm      DecentralizedEnvironmentalNotificationMessage
	X_asn_ctx asn_struct_ctx
}

type TimestampIts struct {
	Buf  *uint64
	Size int
}
