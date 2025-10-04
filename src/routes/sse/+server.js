// src/routes/sse/+server.js
import { produce } from "sveltekit-sse";

export async function POST() {
  return produce(async function start({ emit }) {
    try {
      const response = await fetch("http://localhost:8976/v1/sse", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Origin: "https://app1.local",
        },
      });

      if (!response.body) {
        emit("error", "No response body from Go SSE");
        return;
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        try {
          const { done, value } = await reader.read();
          if (done) {
            emit("close", "Upstream SSE closed");
            break;
          }

          buffer += decoder.decode(value, { stream: true });
          const parts = buffer.split("\n\n");
          buffer = parts.pop() || "";

          for (const part of parts) {
            const line = part.trim();

            // Heartbeat or empty message
            if (line === "data:" || line === "data") {
              emit("ping", "keep-alive");
              continue;
            }

            if (line.startsWith("data:")) {
              const payload = line.slice(5).trim();

              try {
                JSON.parse(payload);
                emit("message", payload);
              } catch {
                emit("debug", `Skipping non-JSON frame: ${payload}`);
              }
            }
          }
        } catch (err) {
          emit("error", `Reader error: ${err}`);
          break;
        }
      }
    } catch (err) {
      emit("error", `Connection failed: ${err}`);
    }
  });
}
