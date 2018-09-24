package main

import (
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/JoaoRufino/canopus"
	nats "github.com/nats-io/go-nats"
)

//#cgo LDFLAGS:-L./pkg/amd64 -lit2s-asn-static-cam -lit2s-asn-static-denm -lssl -lcrypto
/*#include <stdlib.h>
#include <errno.h>
#include <syslog.h>
#include <CAM.h>
#include <DENM.h>
#include <time.h>
#include <INTEGER.h>

#define syslog_emerg(msg, ...) syslog(LOG_EMERG, "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)
#define syslog_err(msg, ...)   syslog(LOG_ERR  , "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)

#ifndef NDEBUG
#define syslog_debug(msg, ...) syslog(LOG_DEBUG, "%s(%s:%d) [" msg "]", __FILE__, __func__, __LINE__, ##__VA_ARGS__)
#else
#define syslog_debug(msg, ...)
#endif

int decode(uint8_t *buffer,int file_size, void *msg, int msg_type){
	asn_dec_rval_t dec;
	asn_enc_rval_t er;
	asn_codec_ctx_t *opt_codec_ctx = 0;

	if(msg == NULL) {
		syslog_emerg("calloc() failed: %m");
	}


	if(msg_type==1)
	{
	dec =  uper_decode_complete(opt_codec_ctx,&asn_DEF_DENM, &msg, buffer, file_size);
		switch(dec.code) {
		case RC_OK:
			xer_fprint(stdout, &asn_DEF_DENM, msg);
			return 1;
		case RC_FAIL:
			syslog_debug("Error decoding: RC_FAIL");
			xer_fprint(stdout, &asn_DEF_DENM, msg);
			return 0;
		case RC_WMORE:
			syslog_debug("ERROR decoding: RC_WMORE");
			xer_fprint(stdout, &asn_DEF_DENM, msg);
			return -1;
	}
	}else{
	dec =  uper_decode_complete(opt_codec_ctx,&asn_DEF_CAM, &msg, buffer, file_size);
		switch(dec.code) {
		case RC_OK:
			xer_fprint(stdout, &asn_DEF_CAM, msg);
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
}

int createDenm(uint8_t *buffer, int32_t *lat, int32_t *lng, int termination ,uint8_t *causeCode, uint8_t *subCauseCode, int size) {
	DENM_t *denm_tx;
	int retval, i;
	struct timespec *systemtime;
	asn_enc_rval_t retval_enc;
	uint16_t seq_number=0;
	uint32_t msg_size=0;
	long timestampits;
	uint8_t *headers, *ptr, *print;
	INTEGER_t *tmp_value; // One INTEGER_t per usage warning!
	INTEGER_t *tmp_value2; // One INTEGER_t per usage warning!

	// Allocate the INTEGER_t
	tmp_value = calloc(1, sizeof(INTEGER_t));
	if(tmp_value == NULL) {
		syslog_err("unable to allocate tmp_value: %m");
		exit(1);
	}

	// Allocate the INTEGER_t
	tmp_value2 = calloc(1, sizeof(INTEGER_t));
	if(tmp_value2 == NULL) {
		syslog_err("unable to allocate tmp_value: %m");
		exit(1);
	}

	// Allocate the struct timespec
	systemtime = (struct timespec *)malloc(sizeof(struct timespec));
	if (systemtime == NULL) {
		syslog_err("malloc() failed: %m");
		exit(1);
	}

	// Allocate the DENM_t denm_tx
	denm_tx = calloc(1, sizeof(DENM_t)); // not malloc!
	if(denm_tx == NULL) {
		syslog_err("calloc() failed");
		exit(1);
	}

		denm_tx->header.protocolVersion = 1;
		denm_tx->header.messageID = 1;
		denm_tx->header.stationID = 101;
		denm_tx->denm.management.actionID.originatingStationID=999;
		denm_tx->denm.management.actionID.sequenceNumber=1;

	retval = clock_gettime(CLOCK_REALTIME,systemtime);
	if (retval != 0){
		perror("clock_gettime() failed");
		exit(1);
	}

	//TODO: Update detectionTime to standard
		timestampits = (long) (systemtime->tv_sec * 1000 + systemtime->tv_nsec/1E6);
		timestampits = timestampits - 1072915200000; // Convert EPOCH to 2004/01/01 00:00:000

		retval = asn_long2INTEGER(tmp_value, timestampits);
		if (retval != 0){
			perror("asn_long2INTEGER() FAILED");
		}
		denm_tx->denm.management.detectionTime = *tmp_value;

		//TODO: Update referenceTime to standard
		retval = asn_long2INTEGER(tmp_value2, timestampits-(60*1000)); // 1 minute before
		if (retval != 0){
			perror("asn_long2INTEGER() FAILED");
		}
		denm_tx->denm.management.referenceTime = *tmp_value2;

		denm_tx->denm.management.eventPosition.latitude = *lat;
		denm_tx->denm.management.eventPosition.longitude = *lng;
		denm_tx->denm.management.eventPosition.altitude.altitudeValue = (int32_t)(32);
		denm_tx->denm.management.eventPosition.altitude.altitudeConfidence = AltitudeConfidence_alt_050_00;
		denm_tx->denm.management.eventPosition.positionConfidenceEllipse.semiMinorConfidence = SemiAxisLength_unavailable;
		denm_tx->denm.management.eventPosition.positionConfidenceEllipse.semiMajorConfidence = SemiAxisLength_unavailable;
		denm_tx->denm.management.eventPosition.positionConfidenceEllipse.semiMajorOrientation = HeadingValue_unavailable;
		denm_tx->denm.management.stationType = StationType_passengerCar;

		denm_tx->denm.management.relevanceDistance = calloc(1,1);
		*(denm_tx->denm.management.relevanceDistance) = RelevanceDistance_lessThan5km;

		denm_tx->denm.management.relevanceTrafficDirection = calloc(1,1);
		*(denm_tx->denm.management.relevanceTrafficDirection) = RelevanceTrafficDirection_allTrafficDirections;

		denm_tx->denm.management.validityDuration = calloc(1,1);
		*(denm_tx->denm.management.validityDuration) = 30; //0.5	 minute

		if(termination){
		denm_tx->denm.management.termination = calloc(1,sizeof(long));
       *(denm_tx->denm.management.termination) = Termination_isCancellation;
       }

		denm_tx->denm.situation = calloc(1,sizeof(SituationContainer_t));
		denm_tx->denm.situation->informationQuality = InformationQuality_lowest;

		// cause and subcause codes!!!!!!
		denm_tx->denm.situation->eventType.causeCode = *causeCode;
		denm_tx->denm.situation->eventType.subCauseCode = *subCauseCode;
		denm_tx->denm.location = calloc(1,sizeof(LocationContainer_t));
		denm_tx->denm.location->traces.list.count = 1;
		denm_tx->denm.location->traces.list.size = 1;
		denm_tx->denm.location->traces.list.array = calloc(1,sizeof(PathHistory_t));
		denm_tx->denm.location->traces.list.array[0] = calloc(1,sizeof(PathHistory_t));
		denm_tx->denm.location->traces.list.array[0]->list.size=1;
		denm_tx->denm.location->traces.list.array[0]->list.count=1;
		denm_tx->denm.location->traces.list.array[0]->list.array = calloc(1,sizeof(PathPoint_t));
		denm_tx->denm.location->traces.list.array[0]->list.array[0] = calloc(1,sizeof(PathPoint_t));
		denm_tx->denm.location->traces.list.array[0]->list.array[0]->pathPosition.deltaLatitude = DeltaLatitude_unavailable;
		denm_tx->denm.location->traces.list.array[0]->list.array[0]->pathPosition.deltaLongitude = DeltaLongitude_unavailable;
		denm_tx->denm.location->traces.list.array[0]->list.array[0]->pathPosition.deltaAltitude = DeltaAltitude_unavailable;

		xer_fprint(stdout, &asn_DEF_DENM, denm_tx);
		memset(buffer,0x00,200);
		retval_enc = uper_encode_to_buffer(&asn_DEF_DENM, denm_tx, buffer,200);
		if(retval_enc.encoded == -1) {
			return -1;
		} else {
			return 1;
		}
	 	int createCAM(uint8_t *buffer){
		CAM_t *cam_tx;
		int retval, i;
		struct timespec *systemtime;
		asn_enc_rval_t retval_enc;
		uint16_t seq_number=0;
		uint32_t msg_size=0;
		long timestampits,generationdeltatime;
		uint8_t *headers, *ptr, *print;
		INTEGER_t *tmp_value; // One INTEGER_t per usage warning!
		INTEGER_t *tmp_value2; // One INTEGER_t per usage warning!

				// Allocate the INTEGER_t
		tmp_value = calloc(1, sizeof(INTEGER_t));
		if(tmp_value == NULL) {
			syslog_err("unable to allocate tmp_value: %m");
			exit(1);
		}

		// Allocate the INTEGER_t
		tmp_value2 = calloc(1, sizeof(INTEGER_t));
		if(tmp_value2 == NULL) {
			syslog_err("unable to allocate tmp_value: %m");
			exit(1);
		}

		// Allocate the struct timespec
		systemtime = (struct timespec *)malloc(sizeof(struct timespec));
		if (systemtime == NULL) {
			syslog_err("malloc() failed: %m");
			exit(1);
		}

		// Allocate the DENM_t denm_tx
		denm_tx = calloc(1, sizeof(DENM_t)); // not malloc!
		if(denm_tx == NULL) {
			syslog_err("calloc() failed");
			exit(1);
		}

		denm_tx->header.protocolVersion = 1;
		denm_tx->header.messageID = 1;
		denm_tx->header.stationID = 101;


		timestampits = (long) (systemtime->tv_sec * 1000 + systemtime->tv_nsec/1E6);
		timestampits = timestampits - 1072915200000; // EPOCH -> 2004/01/01 00:00:000
		generationdeltatime = timestampits % 64536; // generationDeltaTime = TimestampIts mod 65 536
		cam_tx->cam.generationDeltaTime = generationdeltatime;

		cam_tx->cam.camParameters.basicContainer.stationType = StationType_roadSideUnit;

		cam_tx->cam.camParameters.highFrequencyContainer.present = HighFrequencyContainer_PR_basicVehicleContainerHighFrequency;

		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.heading.headingValue = ((uint32_t)(100*10));
			cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.heading.headingConfidence = 126;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.speed.speedValue = ((uint32_t)(1000*100)); // cm/s
			cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.speed.speedConfidence = 126;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.vehicleWidth = 20;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.vehicleLength.vehicleLengthValue = 46;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.vehicleLength.vehicleLengthConfidenceIndication = VehicleLengthConfidenceIndication_noTrailerPresent;
		//cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.accelerationControl = calloc(1,1);
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.longitudinalAcceleration.longitudinalAccelerationValue = LongitudinalAccelerationValue_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.longitudinalAcceleration.longitudinalAccelerationConfidence = AccelerationConfidence_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.curvatureCalculationMode = CurvatureCalculationMode_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.driveDirection = DriveDirection_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.curvature.curvatureValue = CurvatureValue_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.curvature.curvatureConfidence = CurvatureConfidence_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.yawRate.yawRateValue = YawRateValue_unavailable;
		cam_tx->cam.camParameters.highFrequencyContainer.choice.basicVehicleContainerHighFrequency.yawRate.yawRateConfidence = YawRateConfidence_unavailable;

		cam_tx->cam.camParameters.lowFrequencyContainer = calloc(1,sizeof(LowFrequencyContainer_t));
		cam_tx->cam.camParameters.lowFrequencyContainer->present = LowFrequencyContainer_PR_basicVehicleContainerLowFrequency;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.vehicleRole = VehicleRole_default;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.exteriorLights.buf = calloc(1,1);
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.exteriorLights.bits_unused = 0;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.exteriorLights.size = 1;
		*(cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.exteriorLights.buf) = 0x00;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.pathHistory.list.count = 0;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.pathHistory.list.size = 40;
		cam_tx->cam.camParameters.lowFrequencyContainer->choice.basicVehicleContainerLowFrequency.pathHistory.list.array = calloc(1, sizeof(PathPoint_t *));

		return 1;
	 	}
}*/
import "C"

var (
	receive    = make(chan *frontendMsg)
	conn       canopus.Connection
	natsServer = "nats://" + os.Getenv("NATS") + ":4222"
	coapServer = canopus.NewServer()
	rsu        = os.Getenv("RSU")
	port       = os.Getenv("PORT")
	once       = true
	camMsg     = &frontendMsg{}
	camTX      C.CAM_t
	file       = make([]byte, 2360)
)

func main() {

	//Inicializing nats server
	fmt.Println(natsServer)
	nc, er := nats.Connect(natsServer)
	check(er)
	//With json messages
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	check(err)
	defer c.Close()

	// subscribe to frontend messages
	c.Subscribe("msg.frontend", func(p *frontendMsg) {
		fmt.Println("NATS - MESSAGE RECEIVED")
		//fmt.Printf("Received a new message %+v\n", p)
		//fmt.Println(p)
		receive <- p
	})

	// initializng rsu server
	fmt.Println(rsu)
	conn, er = canopus.Dial(rsu)
	check(er)

	coapServer.Post("/CAM", func(req canopus.Request) canopus.Response {
		fmt.Println("COAP - CAM RECEBIDA")
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")
		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()
		er := camMsg.populate(payload, camTX)
		checkC((int)(er), 2)
		//		changeVal := fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude) + "|" + fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude)
		//		server.NotifyChange("/cam/pos", changeVal, false)
		fmt.Println("CAM enviada para nats")
		c.Publish("msg.frontend", camMsg)
		c.Flush()
		return res
	})

	coapServer.Post("/DENM", func(req canopus.Request) canopus.Response {
		denm := C.DENM_t{}
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")

		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()

		er := C.decode((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), (unsafe.Pointer(&denm)), 1)
		checkC((int)(er), 1) /*
			termination := 0
			if denm.denm.management.termination != nil {
				termination = 1*/
		denm_msg := &frontendMsg{}
		/*
			Msg_type:     0,
			ID:           7,
			Lat:          int(denm.denm.management.eventPosition.latitude),
			Lng:          int(denm.denm.management.eventPosition.longitude),
			CauseCode:    int(denm.denm.situation.eventType.causeCode),
			SubCauseCode: int(denm.denm.situation.eventType.subCauseCode),
			Timestamp:    int(*denm.denm.management.detectionTime.buf),
			Termination:  termination}
		*/
		fmt.Println(denm_msg)
		c.Publish("msg.frontend", denm_msg)
		c.Flush()

		//		changeVal := fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude) + "|" + fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude)
		//		server.NotifyChange("/cam/pos", changeVal, false)
		return res
	})

	coapServer.OnMessage(func(msg canopus.Message, inbound bool) {
		//canopus.PrintMessage(msg)
	})

	coapServer.OnObserve(func(resource string, msg canopus.Message) {
		fmt.Println("Observe Requested for " + resource)
	})

	coapServer.OnBlockMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})
	go handleMessages()
	coapServer.ListenAndServe(":" + port)
	<-make(chan struct{})
}

func checkC(i int, j int) {
	if i <= 0 {
		fmt.Println("ERROR ON DECODING!!!")
		switch j {
		case 1:
			fmt.Println("NOT A DENM!!!")
			break
		case 2:
			fmt.Println("NOT A CAM!!!")
			break
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-receive
		log.Println("message from nats")
		// Send it out to every client that is currently connected
		if msg.Header.MessageID == 1 {
			/*er := C.createDenm((*C.uint8_t)(unsafe.Pointer(&file[0])), (*C.int32_t)(unsafe.Pointer(&msg.Lat)), (*C.int32_t)(unsafe.Pointer(&msg.Lng)), (C.int)(msg.Termination), (*C.uint8_t)(unsafe.Pointer(&msg.CauseCode)), (*C.uint8_t)(unsafe.Pointer(&msg.SubCauseCode)), C.int(len(file)))
			checkC((int)(er))
			req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
			req.SetRequestURI("/DENM")
			req.SetPayload(file)
			_, err := conn.Send(req)
			check(err)*/
		} else if msg.Header.MessageID == 2 {
			fmt.Println("PIPE - CAM message Received")
			er := msg.send(file, camTX)
			checkC((int)(er), 2)
			req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
			req.SetRequestURI("/CAM")
			req.SetPayload(file)
			_, err := conn.Send(req)
			check(err)
			fmt.Println("ENVIADA")
		}
	}
}
