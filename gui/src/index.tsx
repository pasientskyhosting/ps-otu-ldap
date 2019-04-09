import React from "react"
import ReactDOM from "react-dom"
import App from "./app"

//@ts-ignore
if (module.hot) {
    //@ts-ignore
    module.hot.accept();
 }

const token = localStorage.getItem('jwt.token')

let isVerified = false

if ( !token  ) {
    isVerified = false
} else {
    isVerified = true
}

ReactDOM.render(<App isAuthenticated={isVerified} />,document.getElementById("root"))