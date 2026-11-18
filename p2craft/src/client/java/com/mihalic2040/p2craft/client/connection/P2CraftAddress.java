package com.mihalic2040.p2craft.client.connection;

public final class P2CraftAddress {

    private static final String SUFFIX = ".p2p";

    private P2CraftAddress() {}

    public static boolean isP2CraftHost(String host) {
        return host != null
                && host.length() >= SUFFIX.length()
                && host.regionMatches(
                true,
                host.length() - SUFFIX.length(),
                SUFFIX,
                0,
                SUFFIX.length()
        );
    }

    public static String getPeerId(String host) {
        if (!isP2CraftHost(host)) {
            return null;
        }

        String peerId = host.substring(
                0,
                host.length() - SUFFIX.length()
        );

        return peerId.isBlank() ? null : peerId;
    }
}