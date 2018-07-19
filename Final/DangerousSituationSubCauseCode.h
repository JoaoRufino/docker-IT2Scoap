/*
 * Generated by asn1c-0.9.28 (http://lionet.info/asn1c)
 * From ASN.1 module "ITS-Container"
 * 	found in "/usr/local/share/asn1c/standard-modules/ts_102_894_2_v1.2.1.asn1"
 * 	`asn1c -gen-PER`
 */

#ifndef	_DangerousSituationSubCauseCode_H_
#define	_DangerousSituationSubCauseCode_H_


#include <asn_application.h>

/* Including external dependencies */
#include <NativeInteger.h>

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum DangerousSituationSubCauseCode {
	DangerousSituationSubCauseCode_unavailable	= 0,
	DangerousSituationSubCauseCode_emergencyElectronicBrakeEngaged	= 1,
	DangerousSituationSubCauseCode_preCrashSystemEngaged	= 2,
	DangerousSituationSubCauseCode_espEngaged	= 3,
	DangerousSituationSubCauseCode_absEngaged	= 4,
	DangerousSituationSubCauseCode_aebEngaged	= 5,
	DangerousSituationSubCauseCode_brakeWarningEngaged	= 6,
	DangerousSituationSubCauseCode_collisionRiskWarningEngaged	= 7
} e_DangerousSituationSubCauseCode;

/* DangerousSituationSubCauseCode */
typedef long	 DangerousSituationSubCauseCode_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_DangerousSituationSubCauseCode;
asn_struct_free_f DangerousSituationSubCauseCode_free;
asn_struct_print_f DangerousSituationSubCauseCode_print;
asn_constr_check_f DangerousSituationSubCauseCode_constraint;
ber_type_decoder_f DangerousSituationSubCauseCode_decode_ber;
der_type_encoder_f DangerousSituationSubCauseCode_encode_der;
xer_type_decoder_f DangerousSituationSubCauseCode_decode_xer;
xer_type_encoder_f DangerousSituationSubCauseCode_encode_xer;
per_type_decoder_f DangerousSituationSubCauseCode_decode_uper;
per_type_encoder_f DangerousSituationSubCauseCode_encode_uper;

#ifdef __cplusplus
}
#endif

#endif	/* _DangerousSituationSubCauseCode_H_ */
#include <asn_internal.h>