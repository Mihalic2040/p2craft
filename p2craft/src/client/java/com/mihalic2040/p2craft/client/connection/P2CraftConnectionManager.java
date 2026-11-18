
package com.mihalic2040.p2craft.client.connection;

import com.mihalic2040.p2craft.client.nativebridge.P2CraftBridge;

import net.minecraft.client.Minecraft;
import net.minecraft.client.gui.screens.Screen;
import net.minecraft.client.gui.screens.ConnectScreen;
import net.minecraft.client.multiplayer.ServerData;
import net.minecraft.client.multiplayer.TransferState;
import net.minecraft.client.multiplayer.resolver.ServerAddress;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;

public final class P2CraftConnectionManager {

    private static final Logger LOGGER =
            LoggerFactory.getLogger("P2Craft");

    private P2CraftConnectionManager() {}

    public static void connect(
            String peerId,
            Screen parent,
            Minecraft minecraft,
            ServerData serverData,
            boolean quickPlay,
            TransferState transferState
    ) {

        LOGGER.info("Connecting to P2Craft peer: {}", peerId);

        CompletableFuture
                .supplyAsync(() -> P2CraftBridge.createTunnel(peerId))
                .thenAccept(localEndpoint -> {

                    LOGGER.info(
                            "P2P connection established: {} -> {}",
                            peerId,
                            localEndpoint
                    );

                    ServerAddress localAddress =
                            ServerAddress.parseString((String) localEndpoint);

                    minecraft.execute(() -> {

                        ConnectScreen.startConnecting(
                                parent,
                                minecraft,
                                localAddress,
                                serverData,
                                quickPlay,
                                transferState
                        );

                    });

                })
                .exceptionally(error -> {

                    LOGGER.error(
                            "Failed to connect to P2Craft peer: {}",
                            peerId,
                            error
                    );

                    minecraft.execute(() -> minecraft.setScreen(parent));

                    return null;
                });
    }
}
