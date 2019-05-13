    import { Button, ButtonGroup, IconName, Popover, Position, Intent, Classes, Toaster } from "@blueprintjs/core";
    import * as React from "react";
    import APIService from '../services/APIService'

    interface IProps {
        onGroupDeleteHandler: (success: boolean, status_code: number) => void,        
        group_name: string                        
    }

    export default class GroupOptions extends React.PureComponent<IProps> {

        constructor(props: IProps) {

            super(props)
            
        }
        
        public render() {
            
            return (
                <ButtonGroup {...this.state} style={{ minWidth: 120 }} fill={false}>                                  
                    {this.renderDeleteButton()}                
                </ButtonGroup>
            );
        }

        
        private renderDeleteButton() {                
            
            const position = Position.RIGHT_TOP;        

            const { ...popoverProps } = this.state;

            return (            
                <Popover                 
                    position={position}
                    popoverClassName={Classes.POPOVER_CONTENT_SIZING}
                    {...popoverProps}
                >
                    <Button 
                        icon="trash"
                        text="Delete"                    
                        large={false}
                    />
                    <div key="text">
                        <h3>Confirm deletion</h3>
                        <p>Are you sure you want to delete this group? All active users will be expired!</p>
                        <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 15 }}>
                            <Button className={Classes.POPOVER_DISMISS} style={{ marginRight: 10 }}>
                                Cancel
                            </Button>
                            <Button 
                                intent={Intent.DANGER} 
                                className={Classes.POPOVER_DISMISS}
                                name={this.props.group_name}
                                onClick={(e: React.MouseEvent<HTMLElement, MouseEvent>) => this.deleteGroup(e.currentTarget.name) }                        
                            >
                                Delete
                            </Button>
                        </div>
                    </div>
                </Popover>
            );
        }

        private async deleteGroup(groupName: string) {        
        
            const group = await APIService.deleteGroup(groupName)        
           
            // Call login handler
            this.props.onGroupDeleteHandler(APIService.success, APIService.status)

        }

    }