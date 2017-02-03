/*******************************************************************************
* * Copyright (c) 2017 Instituto de Telecomunicações
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
 *    Bruno Silva <brunofernandes@av.it.pt> - creator and main architect
 ******************************************************************************/
package org.eclipse.californium.examples;

import static org.eclipse.californium.core.coap.CoAP.ResponseCode.CONTENT;
import static org.eclipse.californium.core.coap.CoAP.ResponseCode.CREATED;
import static org.eclipse.californium.core.coap.CoAP.ResponseCode.PRECONDITION_FAILED;
import static org.eclipse.californium.core.coap.MediaTypeRegistry.APPLICATION_JSON;
import static org.eclipse.californium.core.coap.MediaTypeRegistry.APPLICATION_OCTET_STREAM;
import static org.eclipse.californium.core.coap.MediaTypeRegistry.TEXT_PLAIN;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.RandomAccessFile;
import java.net.SocketException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.Timer;
import java.util.TimerTask;

import org.eclipse.californium.core.CoapResource;
import org.eclipse.californium.core.CoapServer;
import org.eclipse.californium.core.coap.CoAP.Type;
import org.eclipse.californium.core.network.Endpoint;
import org.eclipse.californium.core.network.interceptors.MessageTracer;
import org.eclipse.californium.core.server.resources.CoapExchange;

public class IT2SCoapServer extends CoapServer {

	private static final String TX_FIFO = "/dev/shm/it2s-DENM-mknod-Tx.code";
	private static final String DENM_MESSAGE = "/dev/shm/it2s-DENM_Tx.code";
	private static final String CAM_MESSAGE = "/dev/shm/it2s-CAM_Tx.code";
	private static final String tx_string = "IT2S\n";

	public static void main(String[] args) {
		try {
			// create server
			CoapServer server = new IT2SCoapServer();
			server.start();

			// add special intercepter for message traces
			for (Endpoint ep : server.getEndpoints()) {
				ep.addInterceptor(new MessageTracer());
			}

		} catch (SocketException e) {
			System.err.println("Failed to initialize server: " + e.getMessage());
		}
	}

	/*
	 * Constructor for a new IT2S CoAP server. Here, the resources of the server
	 * are initialised.
	 */

	public IT2SCoapServer() throws SocketException {
		// provide an instance of a IT2S-CoAP Server resources
		add(new DENM());
		add(new CAM());
		add(new NumberOfConnectedHosts());
		add(new ConnectedHosts());
		add(new RealTimeSchedullingParameter());
		add(new JSON1609());
		add(new Time());
	}

	/*
	 * Definition of the Time Resource
	 */
	class Time extends CoapResource {

		private String time;

		public Time() {
			// set resource identifier
			super("Time");
			// set display name
			getAttributes().setTitle("Time");
			setObservable(true);
			getAttributes().addResourceType("observe");
			getAttributes().setObservable();
			setObserveType(Type.CON);

			// Set timer task scheduling
			Timer timer = new Timer();
			timer.schedule(new TimeTask(), 0, 1000);

		}

		/*
		 * Defines a new timer task to return the current time
		 */
		private class TimeTask extends TimerTask {

			@Override
			public void run() {
				time = getTime();
				// Call changed to notify subscribers
				changed();
			}
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			if (exchange.getRequestOptions().hasIfNoneMatch()) {
				exchange.respond(PRECONDITION_FAILED);
			} else {
				// respond to the request
				time = getTime();
				exchange.respond(CONTENT, time, TEXT_PLAIN);
			}
		}

		/*
		 * Returns the current date and time
		 */
		private String getTime() {
			DateFormat dateFormat = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss.SSS");
			Date time = new Date();
			return dateFormat.format(time);
		}

	}

	/*
	 * Definition of the 1609 Resource
	 */
	class JSON1609 extends CoapResource {

		public JSON1609() {
			// set resource identifier
			super("JSON1609");
			// set display name
			getAttributes().setTitle("JSON1609");
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			try {
				FileInputStream fi_json1609 = new FileInputStream("/dev/shm/it2s-1609-json.txt");
				byte[] payload_json1609 = new byte[fi_json1609.available()];
				fi_json1609.read(payload_json1609, 0, payload_json1609.length);
				// respond to the request
				exchange.respond(CONTENT, payload_json1609, APPLICATION_JSON);
				// exchange.respond("1609 contents!");
				fi_json1609.close();
			} catch (Exception e) {
				e.printStackTrace();
			}
		}
	}

	/*
	 * Definition of the CAM Resource
	 */
	class CAM extends CoapResource {

		public CAM() {

			// set resource identifier
			super("CAM");

			// set display name
			getAttributes().setTitle("CAM");
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			try {
				FileInputStream fi_cam = new FileInputStream(CAM_MESSAGE);
				byte[] payload_cam = new byte[fi_cam.available()];
				fi_cam.read(payload_cam, 0, payload_cam.length);
				// respond to the request
				exchange.respond(CONTENT, payload_cam, APPLICATION_OCTET_STREAM);
				// exchange.respond("CAM contents!");
				fi_cam.close();
			} catch (Exception e) {
				e.printStackTrace();
			}
		}

	}

	/*
	 * Definition of the DENM Resource
	 */
	class DENM extends CoapResource {
		public DENM() {
			// set resource identifier
			super("DENM");
			// set display name
			getAttributes().setTitle("DENM");
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			if (exchange.getRequestOptions().hasIfNoneMatch()) {
				exchange.respond(PRECONDITION_FAILED);
			} else {
				// respond to the request
				try {
					FileInputStream fi_denm = new FileInputStream(DENM_MESSAGE);
					byte[] payload_denm = new byte[fi_denm.available()];
					fi_denm.read(payload_denm, 0, payload_denm.length);
					// respond to the request
					exchange.respond(CONTENT, payload_denm, APPLICATION_OCTET_STREAM);
					// exchange.respond("DENM cont");
					fi_denm.close();
				} catch (Exception e) {
					e.printStackTrace();
				}
			}

		}

		@Override
		public void handlePOST(CoapExchange exchange) {
			// TODO: Save POST payload to file
			try {
				// System.out.println(Arrays.toString());
				File f = new File(DENM_MESSAGE);
				if (!f.exists())
					f.createNewFile();
				FileOutputStream fi_denm = new FileOutputStream(DENM_MESSAGE);
				fi_denm.write(exchange.getRequestPayload());
				fi_denm.close();

				RandomAccessFile pipe = new RandomAccessFile(TX_FIFO, "rw");
				pipe.writeChars(tx_string);
				pipe.close();

			} catch (IOException e) {
				System.err.println("IO Exception at buffered read!!");
				System.exit(-1);
			}

			// Check: Type, Code, has Content-Type
			exchange.respond(CREATED);
		}
	}

	/*
	 * Definition of the NumberOfConnectedHosts Resource
	 */
	class NumberOfConnectedHosts extends CoapResource {

		public NumberOfConnectedHosts() {

			// set resource identifier
			super("NumberOfConnectedHosts");

			// set display name
			getAttributes().setTitle("NumberOfConnectedHosts");
		}

		@Override
		public void handleGET(CoapExchange exchange) {

			// respond to the request
			exchange.respond("1");
		}
	}

	/*
	 * Definition of the getConnectedHosts Resource
	 */
	class ConnectedHosts extends CoapResource {

		public ConnectedHosts() {

			// set resource identifier
			super("ConnectedHosts");

			// set display name
			getAttributes().setTitle("ConnectedHosts");
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			// respond to the request
			exchange.respond("RSU");
		}
	}

	/*
	 * Definition of the RealTimeSchedullingParameter Resource
	 */
	class RealTimeSchedullingParameter extends CoapResource {

		public RealTimeSchedullingParameter() {
			// set resource identifier
			super("RealTimeSchedullingParameter");

			// set display name
			getAttributes().setTitle("RealTimeSchedullingParameter");
		}

		@Override
		public void handleGET(CoapExchange exchange) {
			// respond to the request
			exchange.respond(
					"4131320799,19903842,3328215103,3510290603,3731015349,362649195,1717053403,1116134882,3436267957,1852908280");
		}

	}
}
