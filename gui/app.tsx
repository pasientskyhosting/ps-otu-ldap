import React from "react"
import BLReact from  "bl-react"

interface IProps {
    message: string
}

interface IState {
    count: number
}

class App extends React.Component<IProps, IState> {

    constructor(props) {

        super(props)

        this.state = {
            count: 5
        }

    }

    incrementCount() {
        this.setState({
            count: this.state.count + 1
        })
    }

    render() {
        return (
            <div className="std-view-container">
            <BLReact.Panel>            
                <span>Yes, this is {this.props.message}?</span>
                <div>{this.state.count}</div>                
                <BLReact.Button onClick={this.incrementCount.bind(this)} blStyle="success" label="hest" />            
            </BLReact.Panel>
            </div>
        )
    }

}

export default App