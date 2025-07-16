import { useRef, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";

const WS_URL = "ws://localhost:9000/ws";
const API_URL = "http://localhost:9000/project";

export default function App() {
  const [gitURL, setGitURL] = useState("");
  const [slug, setSlug] = useState("");
  const [generatedSlug, setGeneratedSlug] = useState("");
  const [buildUrl, setBuildUrl] = useState("");
  const [logs, setLogs] = useState([]);
  const socketRef = useRef(null);

  const handleBuild = async () => {
    const res = await fetch(API_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ gitURL, slug }),
    });

    const data = await res.json();
    if (data?.data?.projectSlug) {
      const finalSlug = data.data.projectSlug;
      setGeneratedSlug(finalSlug);
      setBuildUrl(data.data.url);
      setupWebSocket(finalSlug);
    }
  };

  const setupWebSocket = (slugChannel) => {
    const socket = new WebSocket(WS_URL);
    socketRef.current = socket;

    socket.onopen = () => {
      socket.send(slugChannel);
    };

    socket.onmessage = (msg) => {
      setLogs((prev) => [...prev, msg.data]);
    };

    socket.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
    };
  };

  return (
    <div className="min-h-screen bg-gray-100 p-6 flex flex-col items-center">
      <Card className="w-full max-w-2xl shadow-md p-6 bg-white">
        <CardContent className="space-y-6">
          <h1 className="text-2xl font-bold text-center">
            🚀 GoPloy: Deploy Your React App with One Click!
          </h1>
          <div className="space-y-2">
            <Label htmlFor="gitURL">📦 Git Repository URL *</Label>
            <Input
              id="gitURL"
              value={gitURL}
              onChange={(e) => setGitURL(e.target.value)}
              placeholder="https://github.com/user/repo.git"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="slug">🔗 Custom Slug (optional)</Label>
            <Input
              id="slug"
              value={slug}
              onChange={(e) => setSlug(e.target.value)}
              placeholder="cool-project"
            />
          </div>
          <div className="flex justify-center pt-2">
            <Button
              onClick={handleBuild}
              disabled={!gitURL}
              className="bg-black text-white hover:bg-gray-800 transition-colors px-6 py-2 rounded-md cursor-pointer"
            >
              ⚒️ Trigger Build
            </Button>
          </div>
          {generatedSlug && (
            <div className="mt-4 text-center">
              <p className="text-sm text-gray-600">
                Slug: <b>{generatedSlug}</b>
              </p>
              <a
                href={buildUrl}
                className="text-blue-600 underline"
                target="_blank"
                rel="noopener noreferrer"
              >
                Visit Project →
              </a>
            </div>
          )}
        </CardContent>
      </Card>

      {logs.length > 0 && (
        <Card className="w-full max-w-2xl mt-6 bg-black text-white p-4 overflow-y-auto max-h-[400px]">
          <CardContent>
            <h2 className="text-lg font-semibold mb-2">📜 Live Logs</h2>
            <pre className="text-sm whitespace-pre-wrap">{logs.join("\n")}</pre>
          </CardContent>
        </Card>
      )}
    </div>
  );
}