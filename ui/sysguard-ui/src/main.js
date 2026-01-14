document.getElementById("app").innerHTML = `
  <h1>SysGuard AI</h1>
  <button id="refresh">Refresh System Metrics</button>
  <pre id="output"></pre>
`;

document.getElementById("refresh").onclick = async () => {
  try {
    const res = await fetch("http://127.0.0.1:7878/metrics");
    const data = await res.json();
    document.getElementById("output").textContent =
      JSON.stringify(data, null, 2);
  } catch {
    document.getElementById("output").textContent =
      "Agent not running";
  }
};
