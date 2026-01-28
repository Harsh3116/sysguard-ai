document.getElementById("app").innerHTML = `
  <h1>SysGuard AI</h1>
  <button id="health">Check System Health</button>
  <button id="shutdown">Analyze Last Shutdown</button>
  <div id="result"></div>
`;

document.getElementById("health").onclick = async () => {
  const res = await fetch("http://127.0.0.1:7878/health-score");
  const data = await res.json();
  document.getElementById("result").innerHTML = `
    <h2>Health Score: ${data.score}</h2>
    <h3>Status: ${data.status}</h3>
    <ul>${data.reasons.map(r => `<li>${r}</li>`).join("")}</ul>
  `;
};

document.getElementById("shutdown").onclick = async () => {
  const res = await fetch("http://127.0.0.1:7878/shutdown-analysis");
  const data = await res.json();
  document.getElementById("result").innerHTML = `
    <h2>Last Shutdown Analysis</h2>
    <p><strong>Cause:</strong> ${data.cause}</p>
    <p><strong>Severity:</strong> ${data.severity}</p>
    <p>${data.explanation}</p>
    <p><strong>Recommendation:</strong> ${data.recommendation}</p>
  `;
};
