import React from 'react';
import { FormGroup, Intent, Button, Elevation, Card, HTMLSelect, Callout, ButtonGroup, InputGroup, ITagProps, Tag, } from "@blueprintjs/core";

import ValidatedInputGroup from './ValidatedInputGroup';
import APIService, { IGroupCustomProps } from '../services/APIService';
import LDAPGroupSearch from './LDAPGroupSearch';

interface IProps {
    onGroupCreateHandler: (success: boolean) => void,
}

interface IState {
    ldap_group_name: string
    group_name: string
    lease_time: number
    errorMessage: string
    custPropKey: string
    custPropValue: string
    tags: IGroupCustomProps
    disableCreate: boolean
}

type Nullable<T> = T | null

const lease_time_default = 720

export default class GroupCreate extends React.Component<IProps, IState> {

    private custPropsKeyRef: Nullable<HTMLInputElement>

    constructor (props: IProps) {

        super(props)

        this.custPropsKeyRef = null

        this.state = {
            tags: {},
            ldap_group_name: "",
            group_name: "",
            lease_time: 720,
            errorMessage: "",
            custPropKey: "",
            custPropValue: "",
            disableCreate: true
        }

    }

    private addTag(key: string, value: string) {

        let tags = this.state.tags

        // Append to current object
        tags[key] = value;

        this.setState({
            tags: tags,
            custPropKey: "",
            custPropValue: ""
        })

        this.custPropsKeyRef && this.custPropsKeyRef.focus()

    }

    private removeTag(key?: string) {

        if(key) {

            let tags = this.state.tags
            delete tags[key];

            this.setState({
                tags: tags,
                custPropKey: "",
                custPropValue: ""
            })
        }

    }

    private renderTags() {

        let stateTags = this.state.tags
        let tagElements = [] as JSX.Element[];

        const onRemove = (e: React.MouseEvent<HTMLButtonElement>, tagProps: ITagProps) => { this.removeTag(tagProps.id || undefined) }

        // Looping through arrays created from Object.keys
        const tagKeys = Object.keys(stateTags)

        for (const key of tagKeys) {

            tagElements.push(
                <Tag
                    className="tags"
                    minimal={true}
                    onRemove={onRemove}
                    key={key}
                    id={key}
                     >
                    {key}={stateTags[key]}
                </Tag>
            )
        }

        return tagElements
    }

    private async onSubmit() {

        if (this.validateFields()) {

            const group = await APIService.groupCreate(this.state.ldap_group_name, this.state.group_name, this.state.lease_time, this.state.tags)

                if(!APIService.success) {
                    this.setState({
                        errorMessage: APIService.error.error.messages[0].value
                    })
                } else {
                    this.setState({
                        errorMessage: ""
                    })
                }

                // call login handler
                this.props.onGroupCreateHandler(APIService.success)

        }

    }

    private onLDAPGroupChange(ldap_group_name: string) {
        this.setState({
            ldap_group_name: ldap_group_name
        })
    }

    render () {

        return (

            <div className="groups-create card">
                <Card interactive={true} elevation={Elevation.FOUR}>
                    <div className="groups-create-content">
                    <h2>Create Group</h2>
                    {this.state.errorMessage ? <Callout title="Error" className="group-create-error-message" intent={Intent.DANGER} >{this.state.errorMessage}</Callout> : null }

                    <FormGroup
                     label="LDAP group"
                     labelFor="ldap-groups"
                    >
                    <LDAPGroupSearch
                        id="ldap-groups"
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => this.onLDAPGroupChange(e.target.value)}
                    />
                    </FormGroup>

                    <FormGroup
                     label="Group name"
                     labelFor="text-input"
                    >
                    <ValidatedInputGroup
                        placeholder="Group name..."
                        onKeyEnter={() => {
                            this.onSubmit()
                        }}
                        value={this.state.group_name}
                        validate={(currentValue: string) => this.validateGroupName(currentValue)}
                        errorMessage={(currentValue: string) =>{
                            return "Length must be greater than 2, and be URL friendly"
                        }}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                            this.setState({
                                group_name: e.target.value
                            })
                        }}
                    />
                    </FormGroup>
                    <FormGroup
                     label="Lease time"
                     labelFor="lease-time"
                    >
                    <HTMLSelect
                        id="lease-time"
                        value={this.state.lease_time}
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                            this.setState({
                                lease_time: parseInt(e.target.value)
                            })
                        }}
                    >
                    <option value="60">1 hour</option>
                    <option value="720">12 hours</option>
                    <option value="1440">24 hours</option>
                    <option value="10080">1 week</option>
                    <option value="20160">2 weeks</option>
                    </HTMLSelect>
                    </FormGroup>

                    <FormGroup
                     label="Custom properties"
                     labelFor="custom-props"
                    >

                    <ButtonGroup
                    >
                        <ValidatedInputGroup
                            inputRef={(input) => this.custPropsKeyRef = input}
                            value={this.state.custPropKey}
                            validate={(currentValue: string) => this.validateTagKey(currentValue)}
                            errorMessage={(currentValue: string) =>{
                                return "Not a valid key"
                            }}
                            placeholder="Key"
                            onChange={(e) => {
                                this.setState({
                                    custPropKey: e.target.value,
                                })
                            }}
                        >
                        </ValidatedInputGroup>
                        &nbsp;
                        <ValidatedInputGroup
                            value={this.state.custPropValue}
                            onKeyEnter={() => {
                                this.addTag(this.state.custPropKey, this.state.custPropValue)
                            }}
                            validate={(currentValue: string) => this.validateTagValue(currentValue)}
                            errorMessage={(currentValue: string) =>{
                                return "Not a valid key"
                            }}
                            placeholder="Value"
                            onChange={(e) => {
                                this.setState({
                                    custPropValue: e.target.value
                                })
                            }}
                        >
                        </ValidatedInputGroup>
                        <Button
                            icon="plus"
                            text=""
                            onClick={(e: React.MouseEvent<HTMLElement, MouseEvent> ) => {
                                if(this.validateTagKey(this.state.custPropKey) && this.validateTagValue(this.state.custPropValue)) {
                                    this.addTag(this.state.custPropKey, this.state.custPropValue)
                                }
                            }}
                        ></Button>
                    </ButtonGroup>
                    </FormGroup>

                    {this.renderTags()}

                    <Button
                        style={{ width: "100%", marginTop: "20px" }}
                        large={false}
                        intent={Intent.PRIMARY}
                        onClick={() => this.onSubmit()}
                    >Create</Button>
                    </div>
                </Card>
            </div>

        )

    }

    private validateLDAPGroupName (ldap_group_name: string): boolean {
        return (ldap_group_name.length == 0 || ldap_group_name == "None" ) ? false : true
    }

    private validateGroupName (group_name: string): boolean {
        return (group_name.length == 0 || !(group_name.length > 2 && group_name.match(/^[_\-0-9a-z]+$/g)) ) ? false : true
    }

    private validateTagKey(str: string): boolean {
        return (str == "") ? false : true
    }

    private validateTagValue(str: string): boolean {
        return (str == "") ? false : true
    }

    private validateFields(): boolean {

        if( this.validateLDAPGroupName(this.state.ldap_group_name) == false ) {
            return false
        }

        if( this.validateGroupName(this.state.group_name) == false ) {
            return false
        }

        return true

    }

}
