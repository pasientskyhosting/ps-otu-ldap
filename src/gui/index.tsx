import React from "react"
import ReactDOM from "react-dom"
import App from "./app"

//@ts-ignore
if (module.hot) {
    //@ts-ignore
    module.hot.accept();
 }

ReactDOM.render(<App message="yes this is jeppe"/>,document.getElementById("root"))
