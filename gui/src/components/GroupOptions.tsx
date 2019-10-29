import { Button, ButtonGroup, IconName, Popover, Position, Intent, Classes, Toaster } from "@blueprintjs/core";
import * as React from "react";
import APIService, { IGroup } from '../services/APIService'

interface IProps {
    onGroupDeleteHandler: (success: boolean, status_code: number) => void,
    onGroupUpdateHandler: () => void,
    group: IGroup
    disableUpdate: boolean
}

interface IState {
    group: IGroup
}

export default class GroupOptions extends React.Component<IProps, IState> {

    constructor(props: IProps) {

        super(props)

        this.state = {
            group: Object.assign({},this.props.group)
        }
    }

    public render() {

        return (
            <ButtonGroup style={{ minWidth: 120 }} fill={false}>
                {this.renderUpdateButton()}
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
                    intent={Intent.DANGER}
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
                            name={this.props.group.group_name}
                            onClick={(e: React.MouseEvent<HTMLElement, MouseEvent>) => this.deleteGroup(e.currentTarget.name) }
                        >
                            Delete
                        </Button>
                    </div>
                </div>
            </Popover>
        );
    }

    private renderUpdateButton() {

        const position = Position.RIGHT_TOP;

        const { ...popoverProps } = this.state;

        return (
            <Button
                icon="tick"
                text="Update"
                large={false}
                disabled={this.props.disableUpdate}
                intent={Intent.SUCCESS}
                onClick={ () => {
                    this.props.onGroupUpdateHandler()
                }}
            />
        );
    }

    private async deleteGroup(groupName: string) {

        const group = await APIService.deleteGroup(groupName)

        // Call login handler
        this.props.onGroupDeleteHandler(APIService.success, APIService.status)

    }

}
