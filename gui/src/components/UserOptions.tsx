import { Button, ButtonGroup, IconName, Popover, Intent, Classes, PopoverInteractionKind, Text } from "@blueprintjs/core";
import * as React from "react";
import APIService from "../services/APIService";
import { AppToaster } from "../services/Toaster";

interface IProps {
    group_name: string
    username: string
    password: string
}

interface IState {
    group_name: string
    username: string
    password: string
}

export default class UserOptions extends React.PureComponent<IProps, IState> {

    constructor(props: IProps) {
        
        super(props)

        this.state = {
            group_name: this.props.group_name,
            username: this.props.username,
            password: this.props.password
        }

    }

    public render() {       

        return (
            <ButtonGroup {...this.state} style={{ minWidth: 120 }} fill={false}>
                {this.renderButton("Create", "document", Intent.NONE, 2, PopoverInteractionKind.CLICK)}
                { this.state.username && this.renderButton("View", "eye-open", Intent.NONE, 1, PopoverInteractionKind.HOVER)}
                {/* { this.state.username && this.renderButton("Expire", "time", Intent.NONE, 0, PopoverInteractionKind.CLICK)}       */}
            </ButtonGroup>
        );
    }

    private renderButton(text: string, iconName: IconName, intent: Intent, content: number, interaction: PopoverInteractionKind) {                
        
        const { ...popoverProps } = this.state;

        return (            
            <Popover                                
                popoverClassName={Classes.POPOVER_CONTENT}
                content={<div className="popover">{this.getContents(content)}</div>}
                hoverOpenDelay={100}
                hoverCloseDelay={200}
                interactionKind={interaction}
                {...popoverProps}      
            >
                <Button 
                    icon={iconName} 
                    text={text}
                    intent={intent}
                 />
                 
            </Popover>
        );
    }

    private getContents(index: number): JSX.Element {
        return [
            <div key="text">
                <h3>Confirm expiration</h3>
                <p>Are you sure you want to expire this user?</p>
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                    <Button className={Classes.POPOVER_DISMISS} style={{ marginRight: 10 }}>
                        Cancel
                    </Button>
                    <Button intent={Intent.DANGER} className={Classes.POPOVER_DISMISS}>
                        Expire
                    </Button>
                </div>
            </div>,
            <div key="text">
                <h3>Copy credentials</h3>
                
                    <code>Username: {this.state.username}</code><br/>
                    <code>Password:&nbsp;{this.state.password}</code>
                
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                    <Button 
                        intent={Intent.SUCCESS} 
                        className={Classes.POPOVER_DISMISS} 
                        style={{ marginRight: 10 }}
                        name={this.state.username}
                        onClick={(e: React.MouseEvent<HTMLElement, MouseEvent>) => this.copyUsername(e.currentTarget.name) }
                    >
                        Copy username
                    </Button>
                    <Button 
                        intent={Intent.PRIMARY} 
                        className={Classes.POPOVER_DISMISS} 
                        name={this.state.password}
                        onClick={(e: React.MouseEvent<HTMLElement, MouseEvent>) => this.copyPassword(e.currentTarget.name) }
                    >
                        Copy password
                    </Button>
                </div>
            </div>,
            <div key="text">
                <h3>Confirm creation</h3>
                <p>Do you want to create a new user?</p>
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                    <Button
                        className={Classes.POPOVER_DISMISS} 
                        style={{ marginRight: 10 }}
                    >
                        Cancel
                    </Button>
                    <Button 
                        intent={Intent.SUCCESS} 
                        className={Classes.POPOVER_DISMISS}
                        name={this.state.group_name}
                        onClick={(e: React.MouseEvent<HTMLElement, MouseEvent>) => this.CreateUser(e.currentTarget.name) }
                    >
                        Create
                    </Button>
                </div>
            </div>,
                 
        ][index];
    }



    private copyUsername(username: string) {

        this.copyToClipBoard(username)
        AppToaster.show(
            {
                intent: Intent.SUCCESS, 
                message: "Username copied to clipboard." 
            }
        )

    }

    private copyPassword(password: string) {

        this.copyToClipBoard(password)
        AppToaster.show(
            {
                intent: Intent.PRIMARY, 
                message: "Password copied to clipboard."
            }
        )

    }

    private copyToClipBoard(string: string) {

        const el = document.createElement('textarea');
        el.value = string;
        document.body.appendChild(el);
        el.select();
        document.execCommand('copy');
        document.body.removeChild(el);

    }


    private async CreateUser(group_name: string) {        
       
        const user = await APIService.createUser(group_name)

        AppToaster.show(
            {
                intent: Intent.SUCCESS, 
                message: "User created successfully." 
            }
        )

        if(user) {
            this.setState({
                username: user.username,
                password: user.password,
                group_name: user.group_name
            })
        }

    }

}