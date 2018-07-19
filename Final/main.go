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


	if(msg_type==0)
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
}*/
import "C"

type frontendMsg struct {
	Msg_type     int `json:"Msg_type"`
	ID           int `json:"ID"`
	Lat          int `json:"Lat"`
	Lng          int `json:"Lng"`
	CauseCode    int `json:"CauseCode"`
	SubCauseCode int `json:"SubCauseCode"`
	Timestamp    int `json:"Timestamp"`
	Termination  int `json:"Termination"`
}

var (
	receive = make(chan *frontendMsg)
	conn    canopus.Connection
)

func main() {
	var ser = "nats://" + os.Getenv("NATS") + ":4222"
	fmt.Println(ser)
	nc, er := nats.Connect(ser)
	check(er)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	check(err)
	defer c.Close()
	c.Subscribe("msg.frontend", func(p *frontendMsg) {
		fmt.Printf("Received a new message %+v\n", p)
		fmt.Println(p)
		receive <- p
	})
	server := canopus.NewServer()
	var rsu = os.Getenv("RSU")
	fmt.Println(rsu)
	conn, er = canopus.Dial(rsu)
	check(er)
	server.Post("/CAM", func(req canopus.Request) canopus.Response {
		cam := CAM{}
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")

		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()
		er := C.decode((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), (unsafe.Pointer(&cam)), 1)
		checkC((int)(er))
		cam_msg := &frontendMsg{
			Msg_type:     1,
			ID:           int(cam.Header.MessageID),
			Lat:          int(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude),
			Lng:          int(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude),
			CauseCode:    0,
			SubCauseCode: 0,
			Timestamp:    int(cam.Cam.GenerationDeltaTime)}
		fmt.Println(cam_msg)
		c.Publish("msg.frontend", cam_msg)
		c.Flush()

		//		changeVal := fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude) + "|" + fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude)
		//		server.NotifyChange("/cam/pos", changeVal, false)
		return res
	})

	server.Post("/DENM", func(req canopus.Request) canopus.Response {
		denm := C.DENM_t{}
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")

		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()

		er := C.decode((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), (unsafe.Pointer(&denm)), 0)
		checkC((int)(er))
		termination := 0
		if denm.denm.management.termination != nil {
			termination = 1
		}
		denm_msg := &frontendMsg{
			Msg_type:     0,
			ID:           7,
			Lat:          int(denm.denm.management.eventPosition.latitude),
			Lng:          int(denm.denm.management.eventPosition.longitude),
			CauseCode:    int(denm.denm.situation.eventType.causeCode),
			SubCauseCode: int(denm.denm.situation.eventType.subCauseCode),
			Timestamp:    int(*denm.denm.management.detectionTime.buf),
			Termination:  termination}

		fmt.Println(denm_msg)
		c.Publish("msg.frontend", denm_msg)
		c.Flush()

		//		changeVal := fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Latitude) + "|" + fmt.Sprint(cam.Cam.CamParameters.BasicContainer.ReferencePosition.Longitude)
		//		server.NotifyChange("/cam/pos", changeVal, false)
		return res
	})

	server.OnMessage(func(msg canopus.Message, inbound bool) {
		//canopus.PrintMessage(msg)
	})

	server.OnObserve(func(resource string, msg canopus.Message) {
		fmt.Println("Observe Requested for " + resource)
	})
	server.OnBlockMessage(func(msg canopus.Message, inbound bool) {
		canopus.PrintMessage(msg)
	})
	go handleMessages()
	server.ListenAndServe(":5683")
	<-make(chan struct{})
}

func checkC(i int) {
	if i <= 0 {
		panic("ERROR ON DECODING!!!")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handleMessages() {
	for {
		file := make([]byte, 200)
		// Grab the next message from the broadcast channel
		msg := <-receive
		log.Println("message from nats")
		// Send it out to every client that is currently connected
		if msg.Msg_type == 0 {
			er := C.createDenm((*C.uint8_t)(unsafe.Pointer(&file[0])), (*C.int32_t)(unsafe.Pointer(&msg.Lat)), (*C.int32_t)(unsafe.Pointer(&msg.Lng)), (C.int)(msg.Termination), (*C.uint8_t)(unsafe.Pointer(&msg.CauseCode)), (*C.uint8_t)(unsafe.Pointer(&msg.SubCauseCode)), C.int(len(file)))
			checkC((int)(er))
			req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
			req.SetRequestURI("/DENM")
			req.SetPayload(file)
			_, err := conn.Send(req)
			check(err)
		}
	}
}
