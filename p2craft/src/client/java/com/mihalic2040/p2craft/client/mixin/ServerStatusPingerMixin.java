package com.mihalic2040.p2craft.client.mixin;

import com.mihalic2040.p2craft.client.connection.P2CraftAddress;
import com.mihalic2040.p2craft.client.nativebridge.P2CraftBridge;

import net.minecraft.client.multiplayer.ServerStatusPinger;
import net.minecraft.client.multiplayer.resolver.ServerAddress;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.spongepowered.asm.mixin.Mixin;
import org.spongepowered.asm.mixin.Unique;
import org.spongepowered.asm.mixin.injection.At;
import org.spongepowered.asm.mixin.injection.Redirect;

@Mixin(ServerStatusPinger.class)
public abstract class ServerStatusPingerMixin {

    @Unique
    private static final Logger p2craft$log =
            LoggerFactory.getLogger("P2Craft");

    @Redirect(
            method = "pingServer(Lnet/minecraft/client/multiplayer/ServerData;Ljava/lang/Runnable;Ljava/lang/Runnable;)V",
            at = @At(
                    value = "INVOKE",
                    target = "Lnet/minecraft/client/multiplayer/resolver/ServerAddress;parseString(Ljava/lang/String;)Lnet/minecraft/client/multiplayer/resolver/ServerAddress;"
            )
    )
    private ServerAddress p2craft$redirectServerAddress(String rawAddress) {
        ServerAddress original = ServerAddress.parseString(rawAddress);

        String peerId = P2CraftAddress.getPeerId(original.getHost());

        // Not a P2Craft address.
        if (peerId == null) {
            return original;
        }

        try {
            // Go returns the existing endpoint if this tunnel already exists.
            String endpoint = P2CraftBridge.createTunnel(peerId);

            if (endpoint == null) {
                p2craft$log.warn(
                        "Failed to create tunnel for peer: {}",
                        peerId
                );
                return original;
            }

            p2craft$log.debug(
                    "Redirecting server-list ping: {} -> {}",
                    peerId,
                    endpoint
            );

            return ServerAddress.parseString(endpoint);

        } catch (RuntimeException e) {
            p2craft$log.error(
                    "Failed to create tunnel for peer: {}",
                    peerId,
                    e
            );
            return original;
        }
    }
}