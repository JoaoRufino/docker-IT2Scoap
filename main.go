package main

import (
	"fmt"
	"unsafe"

	"github.com/JoaoRufino/canopus"
	nats "github.com/nats-io/go-nats"
)

//#cgo LDFLAGS:-L/home/jplr/gitrepos/Golang-projects/asn/Final/libs -lit2s-asn-cam -lit2s-asn-denm
//#cgo CFLAGS:-I/home/jplr/gitrepos/Golang-projects/asn/Final/libs
/*#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <stdint.h>
#include <errno.h>
#include <syslog.h>
#include <CAM.h>
#include <DENM.h>
#include <INTEGER.h>
#include <asn_application.h>
#include <CamParameters.h>
#include <ItsPduHeader.h>
#include <EventPoint.h>
#include <EventHistory.h>
#include <DangerousGoodsExtended.h>
#include <ItineraryPath.h>
#include <PositionOfPillars.h>
#include <RestrictedTypes.h>
#include <Traces.h>
#include <VehicleIdentification.h>
#include <constr_TYPE.h>
#include <DENM.h>
#include <DecentralizedEnvironmentalNotificationMessage.h>

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
*/
import "C"

type frontendMsg struct {
	Msg_type     int `json:"Msg_type"`
	ID           int `json:"ID"`
	Lat          int `json:"Lat"`
	Lng          int `json:"Lng"`
	CauseCode    int `json:"CauseCode"`
	SubCauseCode int `json:"SubCauseCode"`
	Timestamp    int `json:"Timestamp"`
}

func main() {
	nc, er := nats.Connect("nats://193.136.93.80:4222")
	check(er)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	check(err)
	defer c.Close()
	/*c.Subscribe("cam.frontend", func(p *frontendMsg) {
		fmt.Printf("Received a person on subject %+v\n", p)
		fmt.Println(p)
	})*/
	server := canopus.NewServer()

	server.Post("/cam", func(req canopus.Request) canopus.Response {
		cam := CAM{}
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")

		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()
		er := C.decode((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), (unsafe.Pointer(&cam)), 1)
		checkC((int)(er))
		fmt.Println(cam.Cam.CamParameters.HighFrequencyContainer)
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

	server.Post("/denm", func(req canopus.Request) canopus.Response {
		denm := C.DENM_t{}
		fmt.Println(denm)
		msg := canopus.ContentMessage(req.GetMessage().GetMessageId(), canopus.MessageAcknowledgment)
		msg.SetStringPayload("DONE! ")

		res := canopus.NewResponse(msg, nil)
		payload := req.GetMessage().GetPayload().GetBytes()

		er := C.decode((*C.uint8_t)(unsafe.Pointer(&payload[0])), (C.int)(len(payload)), (unsafe.Pointer(&denm)), 0)
		checkC((int)(er))
		denm_msg := &frontendMsg{
			Msg_type:     0,
			ID:           int(denm.header.messageID),
			Lat:          int(denm.denm.management.eventPosition.latitude),
			Lng:          int(denm.denm.management.eventPosition.longitude),
			CauseCode:    int(denm.denm.situation.eventType.causeCode),
			SubCauseCode: int(denm.denm.situation.eventType.subCauseCode),
			Timestamp:    int(*denm.denm.management.detectionTime.buf)}
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
