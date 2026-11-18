package com.mihalic2040.p2craft.client;

import com.mihalic2040.p2craft.client.nativebridge.P2CraftBridge;
import net.fabricmc.api.ClientModInitializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class P2craftClient implements ClientModInitializer {

	private static final Logger LOGGER =
			LoggerFactory.getLogger("P2Craft");

	@Override
	public void onInitializeClient() {

		LOGGER.info("Starting P2Craft node...");

		int result = P2CraftBridge.connect();

		if (result != 0) {
			LOGGER.error(
					"Failed to start P2Craft node. Status: {}",
					P2CraftBridge.getStatus()
			);
			return;
		}

		LOGGER.info("P2Craft node started successfully.");
	}
}