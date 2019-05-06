export interface IGroup {
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

class APIService {

    private baseUrl: string
    private token: string
    public success: boolean  
    public status: number    

    constructor(baseUrl: string) {

        this.baseUrl = baseUrl        
        this.success = false        
        this.status = 0
        
        const token = localStorage.getItem('jwt.token')
            
        this.token = token || ""

    }

    private async parseResponse<TResponseType>(response: Response): Promise<TResponseType | null > {

        console.log("Response status is " + response.status)
        this.status = response.status

        switch(response.status) {           
                
            case 400:                
                this.success = false                                   
                break
            case 401:                               
                this.success = false                
                break
            case 409:                
                this.success = false                
                break
            default:
                this.success = true
        }

        if(response.body) {            
            return await response.json()
        } else {
            return null
        }              

    }

    public async login(username: string, password: string) {        

        try {

            let response = await fetch(this.baseUrl + '/v1/api/auth/authorize', {
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

            this.success = false
            this.status = 0            
        }        

    }    

    public async groupCreate(group_name: string, lease_time: number): Promise<IGroup | null>  {

        try {

            let response = await fetch('/v1/api/groups', {
                    method: 'post',    
                    headers: { "Authorization": `Bearer ${this.token}`  },                   
                    body: JSON.stringify({  
                    group_name, lease_time
                })
            })

            return await (this.parseResponse<IGroup>(response)) || null

        } catch (error) {

            this.success = false
            this.status = 0            
        }        

        return null
    }

    public async getAllGroups(): Promise<IGroup[]> {

        try {

            let response = await fetch(this.baseUrl + '/v1/api/groups', {
                method: 'get',                       
                headers: { "Authorization": `Bearer ${this.token}`  }             
            })
        
            return (await this.parseResponse<IGroup[]>(response)) || []

        } catch (error) {

            this.success = false
            this.status = 0            
        } 

        return []

    }    

    public async deleteGroup(group_name: string) {
        
        try {

            let response = await fetch(this.baseUrl + '/v1/api/groups/' + group_name, {
                method: 'delete',    
                headers: { "Authorization": `Bearer ${this.token}`  }
            })

            await this.parseResponse(response)

        } catch (error) {

            this.success = false
            this.status = 0            
        }

    }

}

export default new APIService('http://localhost:8080');