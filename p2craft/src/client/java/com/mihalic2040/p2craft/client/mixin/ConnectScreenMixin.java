package com.mihalic2040.p2craft.client.mixin;

import com.mihalic2040.p2craft.client.connection.P2CraftAddress;
import com.mihalic2040.p2craft.client.connection.P2CraftConnectionManager;
import net.minecraft.client.Minecraft;
import net.minecraft.client.gui.screens.ConnectScreen;
import net.minecraft.client.gui.screens.Screen;
import net.minecraft.client.multiplayer.ServerData;
import net.minecraft.client.multiplayer.TransferState;
import net.minecraft.client.multiplayer.resolver.ServerAddress;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.spongepowered.asm.mixin.Mixin;
import org.spongepowered.asm.mixin.Unique;
import org.spongepowered.asm.mixin.injection.At;
import org.spongepowered.asm.mixin.injection.Inject;
import org.spongepowered.asm.mixin.injection.callback.CallbackInfo;

import java.util.Locale;

@Mixin(ConnectScreen.class)
public abstract class ConnectScreenMixin {

    @Unique
    private static final Logger p2craft$log =
            LoggerFactory.getLogger("P2Craft");

    @Inject(
            method = "startConnecting",
            at = @At("HEAD"),
            cancellable = true
    )
    private static void p2craft$interceptConnection(
            Screen screen,
            Minecraft minecraft,
            ServerAddress serverAddress,
            ServerData serverData,
            boolean quickPlay,
            TransferState transferState,
            CallbackInfo ci
    ) {
        String host = serverAddress.getHost();

        if (!P2CraftAddress.isP2CraftHost(host)) {
            return;
        }

        ci.cancel();

        String peerId = P2CraftAddress.getPeerId(host);

        if (peerId == null) {
            p2craft$log.warn("Invalid P2Craft address: {}", host);
            minecraft.setScreen(screen);
            return;
        }

        p2craft$log.info("Connecting to P2Craft peer: {}", peerId);

         TODO: P2CraftConnectionManager.connect(
             peerId,
             screen,
             minecraft,
             serverData,
             quickPlay,
             transferState
         );
    }
}

