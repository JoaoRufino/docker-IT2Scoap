/*
 * Generated by asn1c-0.9.28 (http://lionet.info/asn1c)
 * From ASN.1 module "ITS-Container"
 * 	found in "/usr/local/share/asn1c/standard-modules/ts_102_894_2_v1.2.1.asn1"
 * 	`asn1c -gen-PER`
 */

#ifndef	_DangerousGoodsExtended_H_
#define	_DangerousGoodsExtended_H_


#include <asn_application.h>

/* Including external dependencies */
#include "DangerousGoodsBasic.h"
#include <NativeInteger.h>
#include <BOOLEAN.h>
#include <IA5String.h>
#include <UTF8String.h>
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* DangerousGoodsExtended */
typedef struct DangerousGoodsExtended {
	DangerousGoodsBasic_t	 dangerousGoodsType;
	long	 unNumber;
	BOOLEAN_t	 elevatedTemperature;
	BOOLEAN_t	 tunnelsRestricted;
	BOOLEAN_t	 limitedQuantity;
	IA5String_t	*emergencyActionCode	/* OPTIONAL */;
	IA5String_t	*phoneNumber	/* OPTIONAL */;
	UTF8String_t	*companyName	/* OPTIONAL */;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} DangerousGoodsExtended_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_DangerousGoodsExtended;

#ifdef __cplusplus
}
#endif

#endif	/* _DangerousGoodsExtended_H_ */
#include <asn_internal.h>