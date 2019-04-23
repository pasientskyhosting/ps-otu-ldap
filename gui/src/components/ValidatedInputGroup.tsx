import React from 'react';
import { IInputGroupProps, InputGroup, Tooltip, Intent, Position, Label } from '@blueprintjs/core';

interface IProps extends IInputGroupProps {
    validate: (currentValue: string) => boolean
    errorMessage: (currentValue: string) => string
    onKeyDown: (e: React.KeyboardEvent) => void
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
}

interface IState {    
    currentValue: string
}

export default class ValidatedInputGroup extends React.Component<IProps, IState> {
    constructor(props: IProps) {
        super(props)

        this.state = {
            currentValue: ""
        }
    }
 
    private onChangeHandler(e: React.ChangeEvent<HTMLInputElement>) {
        
        this.setState({
            currentValue: e.target.value
        })

        this.props.onChange && this.props.onChange(e)
    }
    
    private renderWithTooltip() {
        return (
            <Tooltip isOpen={!this.props.validate(this.state.currentValue)} autoFocus={false} content={this.props.errorMessage(this.state.currentValue)} intent={Intent.DANGER} position={Position.RIGHT}>
               {this.renderInputGroup()}
            </Tooltip>
        )
    }

    private renderInputGroup() {

        let props: IProps = Object.assign({}, this.props, {
            onChange: (e: React.ChangeEvent<HTMLInputElement>) => this.onChangeHandler(e)
        })
        

        delete props.errorMessage
        delete props.validate

        return (                
            <InputGroup intent={(this.props.validate(this.state.currentValue) ? "none" : Intent.DANGER )} value={this.state.currentValue} {...props}
            />            
        )
    }

    render () {
       
        return this.renderWithTooltip()
       
    }
}
