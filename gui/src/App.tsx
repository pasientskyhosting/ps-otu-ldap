import React from 'react';
import LoginForm from './components/LoginForm'

class App extends React.Component{    

    render() {        

        return (
            <div style={{ margin: "auto", width: "400px"}}>
            <LoginForm/>          
            </div>
        )
    }

}

export default App