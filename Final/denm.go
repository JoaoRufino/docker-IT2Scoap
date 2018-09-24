package main

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

int decodeDenm(uint8_t *buffer,int file_size, void *msg){
	asn_dec_rval_t dec;
	asn_enc_rval_t er;
	asn_codec_ctx_t *opt_codec_ctx = 0;

	if(msg == NULL) {
		syslog_emerg("calloc() failed: %m");
	}
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
}

int createDenm2(uint8_t *buffer, int32_t *lat, int32_t *lng, int termination ,uint8_t *causeCode, uint8_t *subCauseCode, int size) {
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
		memset(buffer,0x00,2360);
		retval_enc = uper_encode_to_buffer(&asn_DEF_DENM, denm_tx, buffer,2360);
		if(retval_enc.encoded == -1) {
			return -1;
		} else {
			return 1;
		}
}*/
import "C"
