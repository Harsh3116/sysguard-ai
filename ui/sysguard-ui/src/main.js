document.getElementById("app").innerHTML = `
  <h1>SysGuard AI</h1>

  <button id="health">Check System Health</button>
  <button id="shutdown">Analyze Last Shutdown</button>
  <button id="startup">Startup Health Summary</button>

  <div id="result"></div>
`;

document.getElementById("health").onclick = async () => {
  const res = await fetch("http://127.0.0.1:7878/health-score");
  const data = await res.json();

  document.getElementById("result").innerHTML = `
    <h2>Health Score: ${data.score}</h2>
    <p>Status: ${data.status}</p>
    <ul>${data.reasons.map(r => `<li>${r}</li>`).join("")}</ul>
  `;
};

document.getElementById("shutdown").onclick = async () => {
  const res = await fetch("http://127.0.0.1:7878/shutdown-analysis");
  const data = await res.json();

  document.getElementById("result").innerHTML = `
    <h2>Last Shutdown Analysis</h2>
    <p><b>Cause:</b> ${data.cause}</p>
    <p><b>Severity:</b> ${data.severity}</p>
    <p>${data.explanation}</p>
    <p><b>Recommendation:</b> ${data.recommendation}</p>
  `;
};

document.getElementById("startup").onclick = async () => {
  const res = await fetch("http://127.0.0.1:7878/startup-summary");
  const data = await res.json();

  document.getElementById("result").innerHTML = `
    <h2>Startup Health: ${data.overall_status}</h2>
    <p>${data.message}</p>

    ${
      data.high_impact_apps.length > 0
        ? `<ul>${data.high_impact_apps.map(a => `<li>${a}</li>`).join("")}</ul>`
        : ""
    }

    <p><b>Recommendation:</b> ${data.recommendation}</p>
  `;
};
