/*******************************************************************************
 * Copyright (c) 2014 Institute for Pervasive Computing, ETH Zurich and others.
 * 
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 * 
 * The Eclipse Public License is available at
 *    http://www.eclipse.org/legal/epl-v10.html
 * and the Eclipse Distribution License is available at
 *    http://www.eclipse.org/org/documents/edl-v10.html.
 * 
 * Contributors:
 *    Matthias Kovatsch - creator and main architect
 *    Martin Lanter - architect and initial implementation
 ******************************************************************************/
package org.eclipse.californium.examples;

import java.io.*;
import java.util.*;
import org.eclipse.californium.core.CoapClient;
import org.eclipse.californium.core.CoapHandler;
import org.eclipse.californium.core.CoapObserveRelation;
import org.eclipse.californium.core.CoapResponse;
import org.eclipse.californium.core.coap.MediaTypeRegistry;

public class CoAPClientExample {
	private static final String RX_FIFO = "/dev/shm/it2s-DENM-mknod-Rx.code";
	private static final String DENM_MESSAGE = "/dev/shm/it2s-DENM_Rx.code";

	public static void main(String[] args) {
		try {
			CoapClient client = new CoapClient("coap://127.0.0.1:5683/DENM");
			// CoapClient client = new CoapClient("coap://127.0.0.1:5683/CAM");
			while (1 == 1) {
				try {
					RandomAccessFile pipe = new RandomAccessFile(RX_FIFO, "r");
					// Read data from named fifo
					String res = pipe.readLine();
					if (res.length() > 1) {
						System.out.println("Data Available");
						// Read binary frame from shared memory
						FileInputStream fi1 = new FileInputStream(DENM_MESSAGE);
						byte[] payload_denm = new byte[fi1.available()];
						fi1.read(payload_denm, 0, payload_denm.length);
						client.post(payload_denm, MediaTypeRegistry.APPLICATION_OCTET_STREAM);
						fi1.close();
					}
					// Close named fifo
					pipe.close();

				} catch (IOException ex) {
					System.err.println("IO Exception at buffered read!!");
					System.exit(-1);
				}

			}

		} catch (Exception e) {
			e.printStackTrace();
		}
	}

}
