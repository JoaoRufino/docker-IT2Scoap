package main

//#cgo LDFLAGS:-L./pkg/amd64 -lit2s-asn-static-cam -lit2s-asn-static-denm -lssl -lcrypto
/*#include <stdlib.h>
#include <errno.h>
#include <syslog.h>
#include <CAM.h>
#include <time.h>
#include <INTEGER.h>

#define syslog_emerg(msg, ...) syslog(LOG_EMERG, "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)
#define syslog_err(msg, ...)   syslog(LOG_ERR  , "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)

#ifndef NDEBUG
#define syslog_debug(msg, ...) syslog(LOG_DEBUG, "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)
#else
#define syslog_debug(msg, ...)
#endif

int decodeCam(uint8_t *buffer,int file_size, void *msg){
	asn_dec_rval_t dec;
	asn_enc_rval_t er;
	asn_codec_ctx_t *opt_codec_ctx = 0;

	if(msg == NULL) {
		syslog_emerg("calloc() failed: %m");
	}
	dec =  uper_decode_complete(opt_codec_ctx,&asn_DEF_CAM, &msg, buffer, file_size);
		switch(dec.code) {
		case RC_OK:
		//	xer_fprint(stdout, &asn_DEF_CAM, msg);
			return 1;
		case RC_FAIL:
			syslog_debug("Error decoding: RC_FAIL");
			xer_fprint(stdout, &asn_DEF_CAM, msg);
			return 0;
		case RC_WMORE:
			syslog_debug("ERROR decoding: RC_WMORE");
			xer_fprint(stdout, &asn_DEF_CAM, msg);
			return -1;
	}
}
int createCAM(uint8_t *buffer, CAM_t *cam) {
	asn_enc_rval_t retval_enc;
	memset(buffer,0x00,2360);
	//xer_fprint(stdout, &asn_DEF_CAM, cam);
		retval_enc = uper_encode_to_buffer(&asn_DEF_CAM, cam, buffer,2360);
		if(retval_enc.encoded == -1) {
			return -1;
		} else {
			return 1;
		}
	}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func (msg *frontendMsg) populate(payload []byte, cam C.CAM_t) C.int {
	cam = C.CAM_t{}
	er := C.decodeCam((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), unsafe.Pointer(&cam))
	if er <= 0 {
		fmt.Println(er)
		return er
	}
	//Header
	msg.Header.ProtocolVersion = int(cam.header.protocolVersion)
	msg.Header.MessageID = int(cam.header.messageID)
	msg.Header.StationID = int(cam.header.stationID)
	//
	if msg.Header.MessageID != 2 {
		fmt.Printf("Incorrect message type! Sent %d to /CAM\n", msg.Header.MessageID)
		return -1
	}

	msg.Cam.GenerationDeltaTime = int(cam.cam.generationDeltaTime)

	//basicContainer
	msg.Cam.CamParameters.BasicContainer.StationType = int(cam.cam.camParameters.basicContainer.stationType)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude = int(cam.cam.camParameters.basicContainer.referencePosition.latitude)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude = int(cam.cam.camParameters.basicContainer.referencePosition.longitude)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMajorConfidence = int(cam.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMajorConfidence)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMinorConfidence = int(cam.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMinorConfidence)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMajorOrientation = int(cam.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMajorOrientation)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.Altitude.AltitudeValue = int(cam.cam.camParameters.basicContainer.referencePosition.altitude.altitudeValue)
	msg.Cam.CamParameters.BasicContainer.ReferencePosition.Altitude.AltitudeConfidence = int(cam.cam.camParameters.basicContainer.referencePosition.altitude.altitudeConfidence)

	//highFrequencyContainer
	if cam.cam.camParameters.highFrequencyContainer.present == C.HighFrequencyContainer_PR_basicVehicleContainerHighFrequency {
		//MAGIC FUCKING HAPPENING

		var addr *byte = &cam.cam.camParameters.highFrequencyContainer.choice[0]
		var cast **C.BasicVehicleContainerHighFrequency_t = (**C.BasicVehicleContainerHighFrequency_t)(unsafe.Pointer(&addr))
		var basicContainer *C.BasicVehicleContainerHighFrequency_t = *cast

		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Heading.HeadingValue = int(basicContainer.heading.headingValue)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Heading.HeadingConfidence = int(basicContainer.heading.headingConfidence)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Speed.SpeedValue = int(basicContainer.speed.speedValue)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Speed.SpeedConfidence = int(basicContainer.speed.speedConfidence)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.DriveDirection = int(basicContainer.driveDirection)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VehicleLength.VehicleLengthValue = int(basicContainer.vehicleLength.vehicleLengthValue)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VehicleWidth = int(basicContainer.vehicleWidth)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LongitudinalAcceleration.LongitudinalAccelerationConfidence = int(basicContainer.longitudinalAcceleration.longitudinalAccelerationValue)
		msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LongitudinalAcceleration.LongitudinalAccelerationValue = int(basicContainer.longitudinalAcceleration.longitudinalAccelerationValue)

		if basicContainer.accelerationControl != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.AccelerationControl = int(*basicContainer.accelerationControl.buf)
		}
		if basicContainer.lanePosition != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LanePosition = int(*basicContainer.lanePosition)
		}
		if basicContainer.steeringWheelAngle != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.SteeringWheelAngle.SteeringWheelAngleValue = int(basicContainer.steeringWheelAngle.steeringWheelAngleValue)
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.SteeringWheelAngle.SteeringWheelAngleConfidence = int(basicContainer.steeringWheelAngle.steeringWheelAngleConfidence)
		}
		if basicContainer.lateralAcceleration != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LateralAcceleration.LateralAccelerationValue = int(basicContainer.lateralAcceleration.lateralAccelerationValue)
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LateralAcceleration.LateralAccelerationConfidence = int(basicContainer.lateralAcceleration.lateralAccelerationConfidence)
		}

		if basicContainer.verticalAcceleration != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VerticalAcceleration.VerticalAccelerationValue = int(basicContainer.verticalAcceleration.verticalAccelerationValue)
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VerticalAcceleration.VerticalAccelerationConfidence = int(basicContainer.verticalAcceleration.verticalAccelerationConfidence)
		}

		if basicContainer.performanceClass != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.PerformanceClass = int(*basicContainer.performanceClass)
		}

		if basicContainer.cenDsrcTollingZone != nil {
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.CenDsrcTollingZone.ProtectedZoneLatitude = int(basicContainer.cenDsrcTollingZone.protectedZoneLatitude)
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.CenDsrcTollingZone.ProtectedZoneLongitude = int(basicContainer.cenDsrcTollingZone.protectedZoneLongitude)
			msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.CenDsrcTollingZone.CenDsrcTollingZoneID = int(*basicContainer.cenDsrcTollingZone.cenDsrcTollingZoneID)
		}

	}
	if cam.cam.camParameters.lowFrequencyContainer != nil {
		if cam.cam.camParameters.lowFrequencyContainer.present == C.LowFrequencyContainer_PR_basicVehicleContainerLowFrequency {
			var addr *byte = &cam.cam.camParameters.lowFrequencyContainer.choice[0]
			var cast **C.BasicVehicleContainerLowFrequency_t = (**C.BasicVehicleContainerLowFrequency_t)(unsafe.Pointer(&addr))
			var lowContainer *C.BasicVehicleContainerLowFrequency_t = *cast

			msg.Cam.CamParameters.LowFrequencyContainer.BasicVehicleContainerLowFrequency.VehicleRole = int(lowContainer.vehicleRole)
			msg.Cam.CamParameters.LowFrequencyContainer.BasicVehicleContainerLowFrequency.ExteriorLights = int(*lowContainer.exteriorLights.buf)

			var ptr **C.PathPoint_t = lowContainer.pathHistory.list.array
			for i := 0; i < int(lowContainer.pathHistory.list.count); i++ {
				temp := pathPoint{PathDeltaTime: int(*(*ptr).pathDeltaTime)}
				temp.PathPosition.DeltaLatitude = int((*ptr).pathPosition.deltaLatitude)
				temp.PathPosition.DeltaLongitude = int((*ptr).pathPosition.deltaLongitude)
				temp.PathPosition.DeltaAltitude = int((*ptr).pathPosition.deltaAltitude)
				msg.Cam.CamParameters.LowFrequencyContainer.BasicVehicleContainerLowFrequency.PathHistory.PathPoint = append(msg.Cam.CamParameters.LowFrequencyContainer.BasicVehicleContainerLowFrequency.PathHistory.PathPoint, temp)
				ptr = (**C.PathPoint_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*ptr)))
			}
		}
	}
	if cam.cam.camParameters.specialVehicleContainer != nil {
		var addr *byte = &cam.cam.camParameters.specialVehicleContainer.choice[0]
		if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_emergencyContainer {
			var cast **C.EmergencyContainer_t = (**C.EmergencyContainer_t)(unsafe.Pointer(&addr))
			var emergency *C.EmergencyContainer_t = *cast
			fmt.Println(emergency)
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_dangerousGoodsContainer {
			var cast **C.DangerousGoodsContainer_t = (**C.DangerousGoodsContainer_t)(unsafe.Pointer(&addr))
			var dangerousGoods *C.DangerousGoodsContainer_t = *cast
			fmt.Println(dangerousGoods)
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_safetyCarContainer {
			var cast **C.SafetyCarContainer_t = (**C.SafetyCarContainer_t)(unsafe.Pointer(&addr))
			var safetyCar *C.SafetyCarContainer_t = *cast
			fmt.Println(safetyCar)
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_publicTransportContainer {
			var cast **C.PublicTransportContainer_t = (**C.PublicTransportContainer_t)(unsafe.Pointer(&addr))
			var publicTransport *C.PublicTransportContainer_t = *cast
			fmt.Println(publicTransport)
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_rescueContainer {
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_roadWorksContainerBasic {
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_safetyCarContainer {
		} else if cam.cam.camParameters.specialVehicleContainer.present == C.SpecialVehicleContainer_PR_specialTransportContainer {
		}
	}

	return 1
}

func (msg *frontendMsg) send(file []byte, camTX C.CAM_t) C.int {
	camTX = C.CAM_t{}
	camTX.header.protocolVersion = C.long(msg.Header.ProtocolVersion) //retificar isto
	camTX.header.messageID = C.long(msg.Header.MessageID)
	camTX.header.stationID = C.ulong(msg.Header.StationID)

	camTX.cam.generationDeltaTime = 21676

	camTX.cam.camParameters.basicContainer.stationType = C.long(msg.Cam.CamParameters.BasicContainer.StationType)

	/* read gps modem data */
	//retval = it2s_gps_read(&gps_data);
	camTX.cam.camParameters.basicContainer.referencePosition.altitude.altitudeValue = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.Altitude.AltitudeValue)
	camTX.cam.camParameters.basicContainer.referencePosition.altitude.altitudeConfidence = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.Altitude.AltitudeConfidence)
	camTX.cam.camParameters.basicContainer.referencePosition.latitude = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude)
	camTX.cam.camParameters.basicContainer.referencePosition.longitude = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude)
	camTX.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMinorConfidence = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMinorConfidence)
	camTX.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMajorConfidence = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMajorConfidence)
	camTX.cam.camParameters.basicContainer.referencePosition.positionConfidenceEllipse.semiMajorOrientation = C.long(msg.Cam.CamParameters.BasicContainer.ReferencePosition.PositionConfidenceEllipse.SemiMajorOrientation)

	// Possible highFrequencyContainer modes: HighFrequencyContainer_PR_NOTHING, HighFrequencyContainer_PR_basicVehicleContainerHighFrequency and HighFrequencyContainer_PR_rsuContainerHighFrequency
	//f msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency {
	camTX.cam.camParameters.highFrequencyContainer.present = C.HighFrequencyContainer_PR_basicVehicleContainerHighFrequency
	var addr *byte = &camTX.cam.camParameters.highFrequencyContainer.choice[0]
	var cast **C.BasicVehicleContainerHighFrequency_t = (**C.BasicVehicleContainerHighFrequency_t)(unsafe.Pointer(&addr))
	var basicVehicleContainerHighFrequency *C.BasicVehicleContainerHighFrequency_t = *cast

	basicVehicleContainerHighFrequency.heading.headingValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Heading.HeadingValue))
	basicVehicleContainerHighFrequency.heading.headingConfidence = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Heading.HeadingConfidence))

	basicVehicleContainerHighFrequency.speed.speedValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Speed.SpeedValue)) // cm/s
	basicVehicleContainerHighFrequency.speed.speedConfidence = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Speed.SpeedConfidence))

	basicVehicleContainerHighFrequency.vehicleWidth = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VehicleWidth))

	basicVehicleContainerHighFrequency.vehicleLength.vehicleLengthValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VehicleLength.VehicleLengthValue))
	basicVehicleContainerHighFrequency.vehicleLength.vehicleLengthConfidenceIndication = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.VehicleLength.VehicleLengthConfidenceIndication))

	basicVehicleContainerHighFrequency.longitudinalAcceleration.longitudinalAccelerationValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LongitudinalAcceleration.LongitudinalAccelerationValue))
	basicVehicleContainerHighFrequency.longitudinalAcceleration.longitudinalAccelerationConfidence = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.LongitudinalAcceleration.LongitudinalAccelerationConfidence))

	basicVehicleContainerHighFrequency.curvatureCalculationMode = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.CurvatureCalculationMode))

	basicVehicleContainerHighFrequency.driveDirection = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.DriveDirection))
	basicVehicleContainerHighFrequency.curvature.curvatureValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Curvature.CurvatureValue))
	basicVehicleContainerHighFrequency.curvature.curvatureConfidence = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.Curvature.CurvatureConfidence))

	basicVehicleContainerHighFrequency.yawRate.yawRateValue = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.YawRate.YawRateValue))
	basicVehicleContainerHighFrequency.yawRate.yawRateConfidence = (C.long(msg.Cam.CamParameters.HighFrequencyContainer.BasicVehicleContainerHighFrequency.YawRate.YawRateConfidence))

	res := C.createCAM((*C.uint8_t)(unsafe.Pointer(&file[0])), (*C.CAM_t)(unsafe.Pointer(&camTX)))
	return res
	/*} else if msg.Cam.CamParameters.HighFrequencyContainer.RsuContainerHighFrequency != nil {
		camTX.cam.camParameters.highFrequencyContainer.present = C.HighFrequencyContainer_PR_rsuContainerHighFrequency
	}

	if msg.Cam.CamParameters.LowFrequencyContainer != nil {

	}

	camTX.cam.camParameters.lowFrequencyContainer = C.calloc(1, sizeof(C.LowFrequencyContainer_t))
	camTX.cam.camParameters.lowFrequencyContainer.present = LowFrequencyContainer_PR_basicVehicleContainerLowFrequency
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.vehicleRole = VehicleRole_default
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.exteriorLights.buf = C.calloc(1, 1)
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.exteriorLights.bits_unused = 0
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.exteriorLights.size = 1
	*(camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.exteriorLights.buf) = 0x00
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.pathHistory.list.count = 0
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.pathHistory.list.size = 40
	camTX.cam.camParameters.lowFrequencyContainer.choice.basicVehicleContainerLowFrequency.pathHistory.list.array = C.calloc(1, sizeof(*C.PathPoint_t))
	*/
	return 1
}
