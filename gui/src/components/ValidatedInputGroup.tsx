import React from 'react';
import { IInputGroupProps, InputGroup, Tooltip, Intent, Position, Label } from '@blueprintjs/core';

interface IProps extends IInputGroupProps {
    validate: (currentValue: string) => boolean
    errorMessage: (currentValue: string) => string
    onSubmit: () => void    
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    defaultValue?: string
}

interface IState {    
    currentValue: string
    valid: boolean    
}

export default class ValidatedInputGroup extends React.Component<IProps, IState> {

    constructor(props: IProps) {

        super(props)
        
        this.state = {
            currentValue: this.props.defaultValue || "",
            valid: true            
        }
    }
 
    private onChangeHandler(e: React.ChangeEvent<HTMLInputElement>) {
        
        this.setState({
            currentValue: e.target.value,            
            valid: this.props.validate(e.target.value)
        })

        this.props.onChange && this.props.onChange(e)
    }
    
    private renderWithTooltip() {
        return (
            <Tooltip isOpen={!(this.state.valid)} autoFocus={false} content={this.props.errorMessage(this.state.currentValue)} intent={Intent.DANGER} position={Position.RIGHT}>
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
            <InputGroup 
                onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => this.onKeyDownHandler(e)} 
                intent={(this.state.valid) ? "none" :  Intent.DANGER } 
                value={this.state.currentValue} 
                {...props}
            />            
        )
    }

    private onKeyDownHandler(e: React.KeyboardEvent<HTMLInputElement>) {

        if(e.keyCode == 13) {            
            if ( this.props.validate(this.state.currentValue) ) {
                this.setState( {valid: true} )
                this.props.onSubmit()
            } else {
                this.setState({valid: false})
            }
        }

    }

    render () {
       
        return this.renderWithTooltip()
       
    }
}
