import axios from "axios"
const defaultErr = new Error('Network or unknown error')

export async function Delete<T>(endpoint: string, data?: any) {
    try {
        return (await axios.delete<T>(endpoint, data)).data
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}

export async function Put<T, t>(endpoint: string, data: t) {
    try {
        return (await axios.put<T>(endpoint, data)).data
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}

export async function Post<T, t>(endpoint: string, data: t) {
    try {
        return (await axios.post<T>(endpoint, data)).data
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}

export async function Get<T>(endpoint: string) {
    try {
        return (await axios.get<T>(endpoint)).data
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}
