package gnosivela;

/** Raised when the control plane returns a non-2xx response. */
public final class GnosivelaException extends RuntimeException {
    private final int status;
    private final String body;

    public GnosivelaException(int status, String body, String path) {
        super("gnosivela: " + path + " -> " + status + ": " + body);
        this.status = status;
        this.body = body;
    }

    public int status() {
        return status;
    }

    public String body() {
        return body;
    }
}
