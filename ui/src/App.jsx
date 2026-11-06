import { useState } from "react";

export default function App() {
  const [url, setUrl] = useState("");
  const [short, setShort] = useState("");

  async function submit(e) {
    e.preventDefault();
    const res = await fetch("/api/shorten", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    });
    const data = await res.json();
    setShort(`${window.location.origin}/api/${data.code}`);
  }

  return (
    <div style={{ padding: 40 }}>
      <h1>URL Shortener</h1>
      <form onSubmit={submit}>
        <input
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          size={50}
        />
        <button type="submit">Shorten</button>
      </form>
      {short && (
        <p>
          Short URL: <a href={short}>{short}</a>
        </p>
      )}
    </div>
  );
}
