export interface ILDAPGroup {
    ldap_group_name: string
}

export interface IGroup {
    ldap_group_name: string
    group_name: string    
    lease_time: number
    create_time: number
    create_by: string
}

export interface IUser {
    username: string
    password: string
    group_name: string
    expire_time: number
    create_time: number
    create_by: string
}

export interface ILogin {
    token: string
}

export interface IAPIError {
    error: IAPIErrorContent
	status_code: number
}

interface IAPIErrorContent {
    messages: IAPIErrorMessage[]
}

interface IAPIErrorMessage {
	key:   string
	value: string
}


class APIService {

    private baseUrl: string
    private token: string

    public success: boolean  
    public status: number    
    public error: IAPIError

    constructor(baseUrl?: string) {

        this.baseUrl = baseUrl || ""        
        this.success = false        
        this.status = 0

        this.error = {
            error: {
                messages: []
            },
            status_code: 0
        }
        
        const token = localStorage.getItem('jwt.token')
            
        this.token = token || ""

    }

    private async parseResponse<TResponseType>(response: Response): Promise<TResponseType | null > {

        console.log("Response status is " + response.status)
        this.status = response.status

        switch(response.status) {           
                
            case 200:
            case 201:
            case 203:
            case 204:
                this.success = true
                break;            
            default:
                this.success = false
        }
        
        if(this.success) {  
                        
            try {                
                return await response.json()
            } catch (error) { 
                console.log(error)                
                return null
            }    
            
        } else {            
            
            this.error = await response.json()            
            return null
        }              

    }

    public async login(username: string, password: string) {        

        try {

            let response = await fetch(this.baseUrl + '/auth', {
                method: 'post',                        
                body: JSON.stringify({  
                    username, password
                })
            })
    
            const payload = await this.parseResponse<ILogin>(response)    
            
            if(payload) {
                localStorage.setItem('jwt.token', payload.token)                                   
                let tokenPayload = JSON.parse(atob(payload.token.split(".")[1]))
                this.token = payload.token                
            }

        } catch (error) {
            this.resetConn()           
        }        

    }    

    public async groupCreate(ldap_group_name: string, group_name: string, lease_time: number): Promise<IGroup | null>  {

        try {

            let response = await fetch(this.baseUrl + '/ldap-groups/' + ldap_group_name + '/groups', {
                    method: 'post',    
                    headers: { "Authorization": `Bearer ${this.token}`  },                   
                    body: JSON.stringify({  
                    group_name, lease_time
                })
            })

            return await (this.parseResponse<IGroup>(response)) || null

        } catch (error) {

            this.resetConn()
        }        

        return null
    }

    public async groupUpdate(ldap_group_name: string, group_name: string, lease_time: number): Promise<IGroup | null>  {

        try {

            let response = await fetch(this.baseUrl + '/ldap-groups/' + ldap_group_name + '/groups', {
                    method: 'patch',    
                    headers: { "Authorization": `Bearer ${this.token}`  },                   
                    body: JSON.stringify({  
                    group_name, lease_time
                })
            })

            return await (this.parseResponse<IGroup>(response)) || null

        } catch (error) {

            this.resetConn()
        }        

        return null
    }

    public async getAllGroups(): Promise<IGroup[]> {

        try {

            let response = await fetch(this.baseUrl + '/groups', {
                method: 'get',                       
                headers: { "Authorization": `Bearer ${this.token}`  }             
            })
        
            return (await this.parseResponse<IGroup[]>(response)) || []

        } catch (error) {
            
            this.resetConn()   
        } 

        return []

    }
    
    public async getAllLDAPGroups(): Promise<ILDAPGroup[]> {

        try {

            let response = await fetch(this.baseUrl + '/ldap-groups', {
                method: 'get',                       
                headers: { "Authorization": `Bearer ${this.token}`  }             
            })
        
            return (await this.parseResponse<ILDAPGroup[]>(response)) || []

        } catch (error) {
            
            this.resetConn()   
        } 

        return []

    }
    
    public async getAllActiveUsers(): Promise<IUser[]> {

        try {

            let response = await fetch(this.baseUrl + '/users', {
                method: 'get',                       
                headers: { "Authorization": `Bearer ${this.token}`  }             
            })
        
            return (await this.parseResponse<IUser[]>(response)) || []

        } catch (error) {

            this.resetConn()
        } 

        return []

    }
    
    public async createUser(group_name: string): Promise<IUser | null> {

        try {

            let response = await fetch(this.baseUrl + '/groups/' + group_name + '/users', {
                    method: 'post',    
                    headers: { "Authorization": `Bearer ${this.token}`  }                    
            })

            return await (this.parseResponse<IUser>(response)) || null

        } catch (error) {

            this.resetConn()        
        }        

        return null

    }

    public async deleteGroup(group_name: string) {
        
        try {

            let response = await fetch(this.baseUrl + '/groups/' + group_name, {
                method: 'delete',    
                headers: { "Authorization": `Bearer ${this.token}`  }
            })

            await this.parseResponse(response)

        } catch (error) {

            this.resetConn()
        }

    }

    private resetConn () {

        this.success = false
        this.status = 0
        this.error = {
            error: {
                messages: []
            },
            status_code: 0
        }

    }

}

export default new APIService('/api/v1');