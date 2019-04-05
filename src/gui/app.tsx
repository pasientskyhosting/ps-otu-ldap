import React from "react"
import { Button, Intent, Spinner, Card, Elevation } from "@blueprintjs/core";


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
            <div style={{ marginLeft: '25px' }} className="std-view-container">
            <Spinner intent={Intent.PRIMARY} />
            <Card interactive={true} elevation={Elevation.TWO}>
                <h5><a href="#">Card heading</a></h5>
                <p>Card content</p>
                <Button>Submit</Button>
            </Card>
            {/* /* <BLReact.Panel>            
                <span>Yes, this is {this.props.message}?</span>
                <div>{this.state.count}</div>                
                <BLReact.Button onClick={this.incrementCount.bind(this)} blStyle="success" label="hest" />            
            </BLReact.Panel> */}
            </div>
        )
    }

}

export default App