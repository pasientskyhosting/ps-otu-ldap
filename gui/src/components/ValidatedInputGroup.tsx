import React from 'react';
import { IInputGroupProps, InputGroup, Tooltip, Intent, Position, Label } from '@blueprintjs/core';

interface IProps extends IInputGroupProps {
    validate: (currentValue: string) => boolean
    errorMessage: (currentValue: string) => string
    onKeyEnter?: () => void    
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    defaultValue?: string
}

interface IState {    
    currentValue: string
    valid: boolean    
    displayError: boolean
}

export default class ValidatedInputGroup extends React.Component<IProps, IState> {

    constructor(props: IProps) {

        super(props)
        
        this.state = {
            currentValue: this.props.defaultValue || "",
            valid: true,         
            displayError: false
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
            <Tooltip isOpen={this.state.displayError} autoFocus={false} content={this.props.errorMessage(this.state.currentValue)} intent={Intent.DANGER} position={Position.RIGHT}>
               {this.renderInputGroup()}
            </Tooltip>
        )
    }

    // Dont show error onBlur if empty
    private onBlurHandler(e: any) {

        if(this.state.currentValue == "" ) {
            this.setState({                                
                valid: true,
                displayError: false
            })
        } else {
            this.setState({                                
                valid: this.props.validate(this.state.currentValue),
                displayError: !this.props.validate(this.state.currentValue)
            })
        }
        
    }
       
    private renderInputGroup() {

        let props: IProps = Object.assign({}, this.props, {
            onChange: (e: React.ChangeEvent<HTMLInputElement>) => this.onChangeHandler(e),
            onBlur: (e: any) => this.onBlurHandler(e),
            onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => this.onKeyEnterHandler(e) 
        })        

        delete props.errorMessage
        delete props.validate
        delete props.onKeyEnter

        return (                
            <InputGroup                                      
                intent={!(this.state.valid) ? Intent.DANGER : "none" } 
                value={this.state.currentValue} 
                {...props}
            />            
        )
    }

    private onKeyEnterHandler(e: React.KeyboardEvent<HTMLInputElement>) {

        if(e.keyCode == 13) {            
            if ( this.props.validate(this.state.currentValue) ) {
                this.setState( {valid: true} )
                this.props.onKeyEnter && this.props.onKeyEnter()
            } else {
                this.setState({valid: false, displayError: true})
            }
        }

    }

    render () {
       
        return this.renderWithTooltip()
       
    }
}
