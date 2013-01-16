window.onerror = (message, url, line) =>
  window.$logger ||= new WSLogger("ws://10.0.0.3:8081/")
  window.$logger.log(message: message, url: url, line: line)

window.addEventListener("load", () -> 
  new Controller(document.querySelector("#controller"))
)
