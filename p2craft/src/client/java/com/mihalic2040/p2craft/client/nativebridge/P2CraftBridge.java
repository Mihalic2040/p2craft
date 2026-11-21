package com.mihalic2040.p2craft.client.nativebridge;

import com.sun.jna.Library;
import com.sun.jna.Native;
import com.sun.jna.Pointer;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Locale;

public final class P2CraftBridge {

    private P2CraftBridge() {}

    private interface NativeApi extends Library {

        int Connect();

        int GetStatus();

        Pointer CreateTunnel(String remoteID);

        void FreeString(Pointer ptr);

        int Disconnect();
    }

    private static final class NativeHolder {

        private static final NativeApi INSTANCE = load();

        private static NativeApi load() {
            String resource = getNativeResource();

            try (var input = P2CraftBridge.class.getResourceAsStream(resource)) {
                if (input == null) {
                    throw new IllegalStateException(
                            "Native library not found inside mod JAR: " + resource
                    );
                }

                String extension = resource.substring(resource.lastIndexOf('.'));

                Path extracted = Files.createTempFile("p2craft-", extension);

                Files.copy(
                        input,
                        extracted,
                        java.nio.file.StandardCopyOption.REPLACE_EXISTING
                );

                extracted.toFile().deleteOnExit();

                return Native.load(
                        extracted.toAbsolutePath().toString(),
                        NativeApi.class
                );

            } catch (Exception e) {
                throw new IllegalStateException(
                        "Failed to load P2Craft native library: " + resource,
                        e
                );
            }
        }

        private static String getNativeResource() {
            String os = System.getProperty("os.name").toLowerCase(Locale.ROOT);
            String arch = System.getProperty("os.arch").toLowerCase(Locale.ROOT);

            boolean x64 = arch.equals("amd64") || arch.equals("x86_64");
            boolean arm64 = arch.equals("aarch64") || arch.equals("arm64");

            if (os.contains("linux")) {
                if (x64) {
                    return "/natives/linux-x86-64/libp2craft.so";
                }

                if (arm64) {
                    return "/natives/linux-aarch64/libp2craft.so";
                }
            }

            if (os.contains("windows")) {
                if (x64) {
                    return "/natives/windows-x86-64/p2craft.dll";
                }
            }

            throw new UnsupportedOperationException(
                    "Unsupported platform: " + os + " / " + arch
            );
        }
    }

    // Starts the P2Craft node.
    public static int connect() {
        return NativeHolder.INSTANCE.Connect();
    }

    // Returns node status.
    public static int getStatus() {
        return NativeHolder.INSTANCE.GetStatus();
    }

    // Creates a tunnel to a peer and returns its local endpoint.
    public static String createTunnel(String peerId) {
        Pointer ptr = NativeHolder.INSTANCE.CreateTunnel(peerId);

        if (ptr == null) {
            return null;
        }

        try {
            return ptr.getString(0, StandardCharsets.UTF_8.name());
        } finally {
            NativeHolder.INSTANCE.FreeString(ptr);
        }
    }

    // Stops the P2Craft node.
    public static int disconnect() {
        return NativeHolder.INSTANCE.Disconnect();
    }
}