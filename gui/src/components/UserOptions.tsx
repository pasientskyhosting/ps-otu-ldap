import { Button, ButtonGroup, IconName, Popover, Position, Intent, Classes } from "@blueprintjs/core";
import * as React from "react";

export default class UserOptions extends React.PureComponent {

    public render() {        
        return (
            <ButtonGroup {...this.state} style={{ minWidth: 120 }} fill={false}>
                {this.renderButton("Create", "document", Intent.NONE, 2)}
                {this.renderButton("View", "eye-open", Intent.NONE, 1)}
                {this.renderButton("Expire", "time", Intent.NONE, 0)}                
            </ButtonGroup>
        );
    }

    private renderButton(text: string, iconName: IconName, intent: Intent, content: number) {                
        
        const { ...popoverProps } = this.state;

        return (            
            <Popover                
                popoverClassName={Classes.POPOVER_CONTENT_SIZING}
                {...popoverProps}
            >
                <Button 
                    icon={iconName} 
                    text={text}
                    intent={intent}
                 />
                 {this.getContents(content)}
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
                <p>Username: kfjrftr-4545-45456565</p>
                <p>Password: fjriouj4855865</p>
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                    <Button intent={Intent.SUCCESS} className={Classes.POPOVER_DISMISS} style={{ marginRight: 10 }}>
                        Username
                    </Button>
                    <Button intent={Intent.PRIMARY} className={Classes.POPOVER_DISMISS}>
                        Password
                    </Button>
                </div>
            </div>,
            <div key="text">
                <h3>Confirm creation</h3>
                <p>Do you want to create a new One Time User?</p>
                <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                    <Button className={Classes.POPOVER_DISMISS} style={{ marginRight: 10 }}>
                        Cancel
                    </Button>
                    <Button intent={Intent.SUCCESS} className={Classes.POPOVER_DISMISS}>
                        Create
                    </Button>
                </div>
            </div>,
                 
        ][index];
    }

}