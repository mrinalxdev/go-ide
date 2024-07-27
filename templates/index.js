document.querySelectorAll(".cosmic").forEach((btn) => {
  btn.addEventListener("click", () => {
    localStorage.setItem("selectedLanguage", btn.dataset.language);
  });
});

const editor = CodeMirror.fromTextArea(document.getElementById("editor"), {
  lineNumbers: true,
  theme: "monokai",
  mode: localStorage.getItem("selectedLanguage") || "python",
});

document
  .querySelector('button[hx-post="/run"]')
  .addEventListener("click", () => {
    const code = editor.getValue();
    const language = localStorage.getItem("selectedLanguage") || "python";
    htmx.ajax("POST", "/run", {
      target: "#output",
      swap: "innerHTML",
      values: { code: code, language: language },
    });
  });
